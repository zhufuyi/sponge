package rabbitmq

import (
	"context"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

// ProducerOption producer option.
type ProducerOption func(*producerOptions)

type producerOptions struct {
	exchangeDeclare *exchangeDeclareOptions
	queueDeclare    *queueDeclareOptions
	queueBind       *queueBindOptions
	deadLetter      *deadLetterOptions

	isPersistent bool // is it persistent

	// If true, the message will be returned to the sender if the queue cannot be
	// found according to its own exchange type and routeKey rules.
	mandatory bool
}

func (o *producerOptions) apply(opts ...ProducerOption) {
	for _, opt := range opts {
		opt(o)
	}
}

// default producer settings
func defaultProducerOptions() *producerOptions {
	return &producerOptions{
		exchangeDeclare: defaultExchangeDeclareOptions(),
		queueDeclare:    defaultQueueDeclareOptions(),
		queueBind:       defaultQueueBindOptions(),
		deadLetter:      defaultDeadLetterOptions(),

		isPersistent: true,
		mandatory:    true,
	}
}

// WithProducerExchangeDeclareOptions set exchange declare option.
func WithProducerExchangeDeclareOptions(opts ...ExchangeDeclareOption) ProducerOption {
	return func(o *producerOptions) {
		o.exchangeDeclare.apply(opts...)
	}
}

// WithProducerQueueDeclareOptions set queue declare option.
func WithProducerQueueDeclareOptions(opts ...QueueDeclareOption) ProducerOption {
	return func(o *producerOptions) {
		o.queueDeclare.apply(opts...)
	}
}

// WithProducerQueueBindOptions set queue bind option.
func WithProducerQueueBindOptions(opts ...QueueBindOption) ProducerOption {
	return func(o *producerOptions) {
		o.queueBind.apply(opts...)
	}
}

// WithDeadLetterOptions set dead letter options.
func WithDeadLetterOptions(opts ...DeadLetterOption) ProducerOption {
	return func(o *producerOptions) {
		o.deadLetter.apply(opts...)
	}
}

// WithProducerPersistent set producer persistent option.
func WithProducerPersistent(enable bool) ProducerOption {
	return func(o *producerOptions) {
		o.isPersistent = enable
	}
}

// WithProducerMandatory set producer mandatory option.
func WithProducerMandatory(enable bool) ProducerOption {
	return func(o *producerOptions) {
		o.mandatory = enable
	}
}

// -------------------------------------------------------------------------------------------

// Producer session
type Producer struct {
	Exchange  *Exchange        // exchange
	QueueName string           // queue name
	conn      *amqp.Connection // rabbitmq connection
	ch        *amqp.Channel    // rabbitmq channel

	// persistent or not
	isPersistent bool
	deliveryMode uint8 // amqp.Persistent or amqp.Transient

	// If true, the message will be returned to the sender if the queue cannot be
	// found according to its own exchange type and routeKey rules.
	mandatory bool

	zapLog *zap.Logger

	exchangeArgs  amqp.Table
	queueArgs     amqp.Table
	queueBindArgs amqp.Table
}

// NewProducer create a producer
func NewProducer(exchange *Exchange, queueName string, connection *Connection, opts ...ProducerOption) (*Producer, error) {
	o := defaultProducerOptions()
	o.apply(opts...)

	// crate a new channel
	ch, err := connection.conn.Channel()
	if err != nil {
		return nil, err
	}

	if exchange.eType == exchangeTypeDelayedMessage {
		if o.exchangeDeclare.args == nil {
			o.exchangeDeclare.args = amqp.Table{
				"x-delayed-type": exchange.delayedMessageType,
			}
		} else {
			o.exchangeDeclare.args["x-delayed-type"] = exchange.delayedMessageType
		}
	}
	// declare the exchange type
	err = ch.ExchangeDeclare(
		exchange.name,
		exchange.eType,
		o.isPersistent,
		o.exchangeDeclare.autoDelete,
		o.exchangeDeclare.internal,
		o.exchangeDeclare.noWait,
		o.exchangeDeclare.args,
	)
	if err != nil {
		_ = ch.Close()
		return nil, err
	}

	// declare a queue and create it automatically if it doesn't exist, or skip creation if it does.
	if o.deadLetter.isEnabled() {
		if o.queueDeclare.args == nil {
			o.queueDeclare.args = amqp.Table{
				"x-dead-letter-exchange":    o.deadLetter.exchangeName,
				"x-dead-letter-routing-key": o.deadLetter.routingKey,
			}
		} else {
			o.queueDeclare.args["x-dead-letter-exchange"] = o.deadLetter.exchangeName
			o.queueDeclare.args["x-dead-letter-routing-key"] = o.deadLetter.routingKey
		}
	}
	q, err := ch.QueueDeclare(
		queueName,
		o.isPersistent,
		o.queueDeclare.autoDelete,
		o.queueDeclare.exclusive,
		o.queueDeclare.noWait,
		o.queueDeclare.args,
	)
	if err != nil {
		_ = ch.Close()
		return nil, err
	}

	args := o.queueBind.args
	if exchange.eType == exchangeTypeHeaders {
		args = exchange.headersKeys
	}
	// binding queue and exchange
	err = ch.QueueBind(
		q.Name,
		exchange.routingKey,
		exchange.name,
		o.queueBind.noWait,
		args,
	)
	if err != nil {
		_ = ch.Close()
		return nil, err
	}

	fields := logFields(queueName, exchange)
	fields = append(fields, zap.Bool("isPersistent", o.isPersistent))

	// create dead letter exchange and queue if enabled
	if o.deadLetter.isEnabled() {
		err = createDeadLetter(ch, o.deadLetter)
		if err != nil {
			_ = ch.Close()
			return nil, err
		}
		fields = append(fields, zap.Any("deadLetter", map[string]string{
			"exchange":   o.deadLetter.exchangeName,
			"queue":      o.deadLetter.queueName,
			"routingKey": o.deadLetter.routingKey,
			"type":       exchangeTypeDirect,
		}))
	}

	deliveryMode := amqp.Persistent
	if !o.isPersistent {
		deliveryMode = amqp.Transient
	}

	connection.zapLog.Info("[rabbit producer] initialized", fields...)

	return &Producer{
		QueueName:    queueName,
		conn:         connection.conn,
		ch:           ch,
		Exchange:     exchange,
		isPersistent: o.isPersistent,
		deliveryMode: deliveryMode,
		mandatory:    o.mandatory,
		zapLog:       connection.zapLog,

		exchangeArgs:  o.exchangeDeclare.args,
		queueArgs:     o.queueDeclare.args,
		queueBindArgs: o.queueBind.args,
	}, nil
}

// PublishDirect send direct type message
func (p *Producer) PublishDirect(ctx context.Context, body []byte) error {
	if p.Exchange.eType != exchangeTypeDirect {
		return fmt.Errorf("invalid exchange type (%s), only supports direct type", p.Exchange.eType)
	}
	return p.ch.PublishWithContext(
		ctx,
		p.Exchange.name,
		p.Exchange.routingKey,
		p.mandatory,
		false,
		amqp.Publishing{
			DeliveryMode: p.deliveryMode,
			ContentType:  "text/plain",
			Body:         body,
		},
	)
}

// PublishFanout send fanout type message
func (p *Producer) PublishFanout(ctx context.Context, body []byte) error {
	if p.Exchange.eType != exchangeTypeFanout {
		return fmt.Errorf("invalid exchange type (%s), only supports fanout type", p.Exchange.eType)
	}
	return p.ch.PublishWithContext(
		ctx,
		p.Exchange.name,
		p.Exchange.routingKey,
		p.mandatory,
		false,
		amqp.Publishing{
			DeliveryMode: p.deliveryMode,
			ContentType:  "text/plain",
			Body:         body,
		},
	)
}

// PublishTopic send topic type message
func (p *Producer) PublishTopic(ctx context.Context, topicKey string, body []byte) error {
	if p.Exchange.eType != exchangeTypeTopic {
		return fmt.Errorf("invalid exchange type (%s), only supports topic type", p.Exchange.eType)
	}
	return p.ch.PublishWithContext(
		ctx,
		p.Exchange.name,
		topicKey,
		p.mandatory,
		false,
		amqp.Publishing{
			DeliveryMode: p.deliveryMode,
			ContentType:  "text/plain",
			Body:         body,
		},
	)
}

// PublishHeaders send headers type message
func (p *Producer) PublishHeaders(ctx context.Context, headersKeys map[string]interface{}, body []byte) error {
	if p.Exchange.eType != exchangeTypeHeaders {
		return fmt.Errorf("invalid exchange type (%s), only supports headers type", p.Exchange.eType)
	}
	return p.ch.PublishWithContext(
		ctx,
		p.Exchange.name,
		p.Exchange.routingKey,
		p.mandatory,
		false,
		amqp.Publishing{
			DeliveryMode: p.deliveryMode,
			Headers:      headersKeys,
			ContentType:  "text/plain",
			Body:         body,
		},
	)
}

// PublishDelayedMessage send delayed type message
func (p *Producer) PublishDelayedMessage(ctx context.Context, delayTime time.Duration, body []byte, opts ...DelayedMessagePublishOption) error {
	if p.Exchange.eType != exchangeTypeDelayedMessage {
		return fmt.Errorf("invalid exchange type (%s), only supports x-delayed-message type", p.Exchange.eType)
	}

	routingKey := p.Exchange.routingKey
	headersKeys := make(map[string]interface{})
	o := defaultDelayedMessagePublishOptions()
	o.apply(opts...)
	switch p.Exchange.delayedMessageType {
	case exchangeTypeTopic:
		if o.topicKey == "" {
			return fmt.Errorf("topic key is required, please set topicKey in DelayedMessagePublishOption")
		}
		routingKey = o.topicKey
	case exchangeTypeHeaders:
		if o.headersKeys == nil {
			return fmt.Errorf("headers keys is required, please set headersKeys in DelayedMessagePublishOption")
		}
		headersKeys = o.headersKeys
	}
	headersKeys["x-delay"] = int(delayTime / time.Millisecond) // delay time: milliseconds

	return p.ch.PublishWithContext(
		ctx,
		p.Exchange.name,
		routingKey,
		p.mandatory,
		false,
		amqp.Publishing{
			DeliveryMode: p.deliveryMode,
			Headers:      headersKeys,
			ContentType:  "text/plain",
			Body:         body,
		},
	)
}

// Close the consumer
func (p *Producer) Close() {
	if p.ch != nil {
		_ = p.ch.Close()
	}
}

// ExchangeArgs returns the exchange declare args.
func (p *Producer) ExchangeArgs() amqp.Table {
	return p.exchangeArgs
}

// QueueArgs returns the queue declare args.
func (p *Producer) QueueArgs() amqp.Table {
	return p.queueArgs
}

// QueueBindArgs returns the queue bind args.
func (p *Producer) QueueBindArgs() amqp.Table {
	return p.queueBindArgs
}

func logFields(queueName string, exchange *Exchange) []zap.Field {
	fields := []zap.Field{
		zap.String("queue", queueName),
		zap.String("exchange", exchange.name),
		zap.String("exchangeType", exchange.eType),
	}
	switch exchange.eType {
	case exchangeTypeDirect, exchangeTypeTopic:
		fields = append(fields, zap.String("routingKey", exchange.routingKey))
	case exchangeTypeHeaders:
		fields = append(fields, zap.Any("headersKeys", exchange.headersKeys))
	case exchangeTypeDelayedMessage:
		fields = append(fields, zap.String("delayedMessageType", exchange.delayedMessageType))
		switch exchange.delayedMessageType {
		case exchangeTypeDirect, exchangeTypeTopic:
			fields = append(fields, zap.String("routingKey", exchange.routingKey))
		case exchangeTypeHeaders:
			fields = append(fields, zap.Any("headersKeys", exchange.headersKeys))
		}
	}
	return fields
}

// -------------------------------------------------------------------------------------------

func createDeadLetter(ch *amqp.Channel, o *deadLetterOptions) error {
	// declare the exchange type
	err := ch.ExchangeDeclare(
		o.exchangeName,
		exchangeTypeDirect,
		true,
		o.exchangeDeclare.autoDelete,
		o.exchangeDeclare.internal,
		o.exchangeDeclare.noWait,
		o.exchangeDeclare.args,
	)
	if err != nil {
		return err
	}

	// declare a queue and create it automatically if it doesn't exist, or skip creation if it does.
	q, err := ch.QueueDeclare(
		o.queueName,
		true,
		o.queueDeclare.autoDelete,
		o.queueDeclare.exclusive,
		o.queueDeclare.noWait,
		o.queueDeclare.args,
	)
	if err != nil {
		return err
	}

	// binding queue and exchange
	err = ch.QueueBind(
		q.Name,
		o.routingKey,
		o.exchangeName,
		o.queueBind.noWait,
		o.queueBind.args,
	)

	return err
}
