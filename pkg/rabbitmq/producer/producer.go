// Package producer is the generic producer-side processing logic for the four modes direct, topic, fanout, headers.
package producer

import (
	"context"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

// ErrClosed closed
var ErrClosed = amqp.ErrClosed

const (
	exchangeTypeDirect  = "direct"
	exchangeTypeTopic   = "topic"
	exchangeTypeFanout  = "fanout"
	exchangeTypeHeaders = "headers"

	// HeadersTypeAll all
	HeadersTypeAll HeadersType = "all"
	// HeadersTypeAny any
	HeadersTypeAny HeadersType = "any"
)

// HeadersType headers type
type HeadersType = string

// Exchange rabbitmq minimum management unit
type Exchange struct {
	name       string                 // exchange name
	eType      string                 // exchange type: direct, topic, fanout, headers
	routingKey string                 // route key
	Headers    map[string]interface{} // this field is required if eType=headers.
}

// NewDirectExchange create a direct exchange
func NewDirectExchange(exchangeName string, routingKey string) *Exchange {
	return &Exchange{
		name:       exchangeName,
		eType:      exchangeTypeDirect,
		routingKey: routingKey,
	}
}

// NewTopicExchange create a topic exchange
func NewTopicExchange(exchangeName string, routingKey string) *Exchange {
	return &Exchange{
		name:       exchangeName,
		eType:      exchangeTypeTopic,
		routingKey: routingKey,
	}
}

// NewFanOutExchange create a fanout exchange
func NewFanOutExchange(exchangeName string) *Exchange {
	return &Exchange{
		name:       exchangeName,
		eType:      exchangeTypeFanout,
		routingKey: "",
	}
}

// NewHeaderExchange create a headers exchange, the headerType supports "all" and "any"
func NewHeaderExchange(exchangeName string, headersType HeadersType, kv map[string]interface{}) *Exchange {
	if kv == nil {
		kv = make(map[string]interface{})
	}

	switch headersType {
	case HeadersTypeAll, HeadersTypeAny:
		kv["x-match"] = headersType
	default:
		kv["x-match"] = HeadersTypeAll
	}

	return &Exchange{
		name:       exchangeName,
		eType:      exchangeTypeHeaders,
		routingKey: "",
		Headers:    kv,
	}
}

// -------------------------------------------------------------------------------------------

// Queue session
type Queue struct {
	queueName string           // queue name
	exchange  *Exchange        // exchange
	conn      *amqp.Connection // rabbitmq connection
	ch        *amqp.Channel    // rabbitmq channel

	// If true, the message will be returned to the sender if the queue cannot be
	// found according to its own exchange type and routeKey rules.
	mandatory bool
	// If true, when exchange sends a message to the queue and finds that there
	// are no consumers on the queue, it returns the message to the sender
	immediate bool
}

// NewQueue create a queue
func NewQueue(queueName string, conn *amqp.Connection, exchange *Exchange, opts ...QueueOption) (*Queue, error) {
	o := defaultProducerOptions()
	o.apply(opts...)

	// crate a new channel
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	// declare a queue and create it automatically if it doesn't exist, or skip creation if it does.
	q, err := ch.QueueDeclare(
		queueName,
		o.queueDeclare.durable,
		o.queueDeclare.autoDelete,
		o.queueDeclare.exclusive,
		o.queueDeclare.noWait,
		o.queueDeclare.args,
	)
	if err != nil {
		return nil, err
	}

	// declare the exchange type
	err = ch.ExchangeDeclare(
		exchange.name,
		exchange.eType,
		o.exchangeDeclare.durable,
		o.exchangeDeclare.autoDelete,
		o.exchangeDeclare.internal,
		o.exchangeDeclare.noWait,
		o.exchangeDeclare.args,
	)
	if err != nil {
		return nil, err
	}

	args := o.queueBind.args
	if exchange.eType == exchangeTypeHeaders {
		args = exchange.Headers
	}
	// Binding queue and exchange
	err = ch.QueueBind(
		q.Name,
		exchange.routingKey,
		exchange.name,
		o.queueBind.noWait,
		args,
	)
	if err != nil {
		return nil, err
	}

	return &Queue{
		queueName: queueName,
		conn:      conn,
		ch:        ch,
		exchange:  exchange,
		mandatory: o.mandatory,
		immediate: o.immediate,
	}, nil
}

// Publish send direct or fanout type message
func (q *Queue) Publish(ctx context.Context, body []byte) error {
	if q.exchange.eType != exchangeTypeDirect && q.exchange.eType != exchangeTypeFanout {
		return fmt.Errorf("invalid exchange type (%s), only supports direct or fanout types", q.exchange.eType)
	}
	return q.ch.PublishWithContext(
		ctx,
		q.exchange.name,
		q.exchange.routingKey,
		q.mandatory,
		q.immediate,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		},
	)
}

// PublishTopic send topic type message
func (q *Queue) PublishTopic(ctx context.Context, topicKey string, body []byte) error {
	if q.exchange.eType != exchangeTypeTopic {
		return fmt.Errorf("invalid exchange type (%s), only supports topic type", q.exchange.eType)
	}
	return q.ch.PublishWithContext(
		ctx,
		q.exchange.name,
		topicKey,
		q.mandatory,
		q.immediate,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		},
	)
}

// PublishHeaders send headers type message
func (q *Queue) PublishHeaders(ctx context.Context, headersKey map[string]interface{}, body []byte) error {
	if q.exchange.eType != exchangeTypeHeaders {
		return fmt.Errorf("invalid exchange type (%s), only supports headers type", q.exchange.eType)
	}
	return q.ch.PublishWithContext(
		ctx,
		q.exchange.name,
		q.exchange.routingKey,
		q.mandatory,
		q.immediate,
		amqp.Publishing{
			Headers:     headersKey,
			ContentType: "text/plain",
			Body:        body,
		},
	)
}

// Close the queue
func (q *Queue) Close() {
	if q.ch != nil {
		_ = q.ch.Close()
	}
}
