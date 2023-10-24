package rabbitmq

import (
	"context"
	"strconv"
	"strings"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

// ConsumerOption consumer option.
type ConsumerOption func(*consumerOptions)

type consumerOptions struct {
	exchangeDeclare *exchangeDeclareOptions
	queueDeclare    *queueDeclareOptions
	queueBind       *queueBindOptions
	qos             *qosOptions
	consume         *consumeOptions

	isPersistent bool // persistent or not
	isAutoAck    bool // auto-answer or not, if false, manual ACK required
}

func (o *consumerOptions) apply(opts ...ConsumerOption) {
	for _, opt := range opts {
		opt(o)
	}
}

// default consumer settings
func defaultConsumerOptions() *consumerOptions {
	return &consumerOptions{
		exchangeDeclare: defaultExchangeDeclareOptions(),
		queueDeclare:    defaultQueueDeclareOptions(),
		queueBind:       defaultQueueBindOptions(),
		qos:             defaultQosOptions(),
		consume:         defaultConsumeOptions(),

		isPersistent: true,
		isAutoAck:    true,
	}
}

// WithConsumerExchangeDeclareOptions set exchange declare option.
func WithConsumerExchangeDeclareOptions(opts ...ExchangeDeclareOption) ConsumerOption {
	return func(o *consumerOptions) {
		o.exchangeDeclare.apply(opts...)
	}
}

// WithConsumerQueueDeclareOptions set queue declare option.
func WithConsumerQueueDeclareOptions(opts ...QueueDeclareOption) ConsumerOption {
	return func(o *consumerOptions) {
		o.queueDeclare.apply(opts...)
	}
}

// WithConsumerQueueBindOptions set queue bind option.
func WithConsumerQueueBindOptions(opts ...QueueBindOption) ConsumerOption {
	return func(o *consumerOptions) {
		o.queueBind.apply(opts...)
	}
}

// WithConsumerQosOptions set consume qos option.
func WithConsumerQosOptions(opts ...QosOption) ConsumerOption {
	return func(o *consumerOptions) {
		o.qos.apply(opts...)
	}
}

// WithConsumerConsumeOptions set consumer consume option.
func WithConsumerConsumeOptions(opts ...ConsumeOption) ConsumerOption {
	return func(o *consumerOptions) {
		o.consume.apply(opts...)
	}
}

// WithConsumerAutoAck set consumer auto ack option.
func WithConsumerAutoAck(enable bool) ConsumerOption {
	return func(o *consumerOptions) {
		o.isAutoAck = enable
	}
}

// WithConsumerPersistent set consumer persistent option.
func WithConsumerPersistent(enable bool) ConsumerOption {
	return func(o *consumerOptions) {
		o.isPersistent = enable
	}
}

// -------------------------------------------------------------------------------------------

// ConsumeOption consume option.
type ConsumeOption func(*consumeOptions)

type consumeOptions struct {
	consumer  string     // used to distinguish between multiple consumers
	exclusive bool       // only available to the program that created it
	noLocal   bool       // if set to true, a message sent by a producer in the same Connection cannot be passed to a consumer in this Connection.
	noWait    bool       // block processing
	args      amqp.Table // additional properties
}

func (o *consumeOptions) apply(opts ...ConsumeOption) {
	for _, opt := range opts {
		opt(o)
	}
}

// default consume settings
func defaultConsumeOptions() *consumeOptions {
	return &consumeOptions{
		consumer:  "",
		exclusive: false,
		noLocal:   false,
		noWait:    false,
		args:      nil,
	}
}

// WithConsumeConsumer set consume consumer option.
func WithConsumeConsumer(consumer string) ConsumeOption {
	return func(o *consumeOptions) {
		o.consumer = consumer
	}
}

// WithConsumeExclusive set consume exclusive option.
func WithConsumeExclusive(enable bool) ConsumeOption {
	return func(o *consumeOptions) {
		o.exclusive = enable
	}
}

// WithConsumeNoLocal set consume noLocal option.
func WithConsumeNoLocal(enable bool) ConsumeOption {
	return func(o *consumeOptions) {
		o.noLocal = enable
	}
}

// WithConsumeNoWait set consume no wait option.
func WithConsumeNoWait(enable bool) ConsumeOption {
	return func(o *consumeOptions) {
		o.noWait = enable
	}
}

// WithConsumeArgs set consume args option.
func WithConsumeArgs(args map[string]interface{}) ConsumeOption {
	return func(o *consumeOptions) {
		o.args = args
	}
}

// -------------------------------------------------------------------------------------------

// QosOption qos option.
type QosOption func(*qosOptions)

type qosOptions struct {
	enable        bool
	prefetchCount int
	prefetchSize  int
	global        bool
}

func (o *qosOptions) apply(opts ...QosOption) {
	for _, opt := range opts {
		opt(o)
	}
}

// default qos settings
func defaultQosOptions() *qosOptions {
	return &qosOptions{
		enable:        false,
		prefetchCount: 0,
		prefetchSize:  0,
		global:        false,
	}
}

// WithQosEnable set qos enable option.
func WithQosEnable() QosOption {
	return func(o *qosOptions) {
		o.enable = true
	}
}

// WithQosPrefetchCount set qos prefetch count option.
func WithQosPrefetchCount(count int) QosOption {
	return func(o *qosOptions) {
		o.prefetchCount = count
	}
}

// WithQosPrefetchSize set qos prefetch size option.
func WithQosPrefetchSize(size int) QosOption {
	return func(o *qosOptions) {
		o.prefetchSize = size
	}
}

// WithQosPrefetchGlobal set qos global option.
func WithQosPrefetchGlobal(enable bool) QosOption {
	return func(o *qosOptions) {
		o.global = enable
	}
}

// -------------------------------------------------------------------------------------------

// Consumer session
type Consumer struct {
	Exchange   *Exchange
	QueueName  string
	connection *Connection
	ch         *amqp.Channel

	exchangeDeclareOption *exchangeDeclareOptions
	queueDeclareOption    *queueDeclareOptions
	queueBindOption       *queueBindOptions
	qosOption             *qosOptions
	consumeOption         *consumeOptions

	isPersistent bool // persistent or not
	isAutoAck    bool // auto ack or not

	zapLog *zap.Logger
}

// Handler message
type Handler func(ctx context.Context, data []byte, tagID string) error

//type Handler func(ctx context.Context, d *amqp.Delivery, isAutoAck bool) error

// NewConsumer create a consumer
func NewConsumer(exchange *Exchange, queueName string, connection *Connection, opts ...ConsumerOption) (*Consumer, error) {
	o := defaultConsumerOptions()
	o.apply(opts...)

	c := &Consumer{
		Exchange:   exchange,
		QueueName:  queueName,
		connection: connection,

		exchangeDeclareOption: o.exchangeDeclare,
		queueDeclareOption:    o.queueDeclare,
		queueBindOption:       o.queueBind,
		qosOption:             o.qos,
		consumeOption:         o.consume,

		isPersistent: o.isPersistent,
		isAutoAck:    o.isAutoAck,

		zapLog: connection.zapLog,
	}

	return c, nil
}

// initialize a consumer session
func (c *Consumer) initialize() error {
	c.connection.mutex.Lock()
	// crate a new channel
	ch, err := c.connection.conn.Channel()
	if err != nil {
		c.connection.mutex.Unlock()
		return err
	}
	c.ch = ch
	c.connection.mutex.Unlock()

	if c.Exchange.eType == exchangeTypeDelayedMessage {
		if c.exchangeDeclareOption.args == nil {
			c.exchangeDeclareOption.args = amqp.Table{
				"x-delayed-type": c.Exchange.delayedMessageType,
			}
		} else {
			c.exchangeDeclareOption.args["x-delayed-type"] = c.Exchange.delayedMessageType
		}
	}
	// declare the exchange type
	err = ch.ExchangeDeclare(
		c.Exchange.name,
		c.Exchange.eType,
		c.isPersistent,
		c.exchangeDeclareOption.autoDelete,
		c.exchangeDeclareOption.internal,
		c.exchangeDeclareOption.noWait,
		c.exchangeDeclareOption.args,
	)
	if err != nil {
		_ = ch.Close()
		return err
	}

	// declare a queue and create it automatically if it doesn't exist, or skip creation if it does.
	queue, err := ch.QueueDeclare(
		c.QueueName,
		c.isPersistent,
		c.queueDeclareOption.autoDelete,
		c.queueDeclareOption.exclusive,
		c.queueDeclareOption.noWait,
		c.queueDeclareOption.args,
	)
	if err != nil {
		_ = ch.Close()
		return err
	}

	args := c.queueBindOption.args
	if c.Exchange.eType == exchangeTypeHeaders {
		args = c.Exchange.headersKeys
	}
	// binding queue and exchange
	err = ch.QueueBind(
		queue.Name,
		c.Exchange.routingKey,
		c.Exchange.name,
		c.queueBindOption.noWait,
		args,
	)
	if err != nil {
		_ = ch.Close()
		return err
	}

	// setting the prefetch value, set channel.Qos on the consumer side to limit the number of messages consumed at a time,
	// balancing message throughput and fairness, and prevent consumers from being hit by sudden bursts of information traffic.
	if c.qosOption.enable {
		err = ch.Qos(c.qosOption.prefetchCount, c.qosOption.prefetchSize, c.qosOption.global)
		if err != nil {
			_ = ch.Close()
			return err
		}
	}

	fields := logFields(c.QueueName, c.Exchange)
	fields = append(fields, zap.Bool("autoAck", c.isAutoAck))
	c.zapLog.Info("[rabbitmq consumer] initialized", fields...)
	return nil
}

func (c *Consumer) consumeWithContext(ctx context.Context) (<-chan amqp.Delivery, error) {
	return c.ch.ConsumeWithContext(
		ctx,
		c.QueueName,
		c.consumeOption.consumer,
		c.isAutoAck,
		c.consumeOption.exclusive,
		c.consumeOption.noLocal,
		c.consumeOption.noWait,
		c.consumeOption.args,
	)
}

// Consume messages for loop in goroutine
func (c *Consumer) Consume(ctx context.Context, handler Handler) {
	go func() {
		ticker := time.NewTicker(time.Second * 2)
		isFirst := true
		for {
			if isFirst {
				isFirst = false
				ticker.Reset(time.Millisecond * 10)
			} else {
				ticker.Reset(time.Second * 2)
			}

			// check connection for loop
			select {
			case <-ticker.C:
				if !c.connection.CheckConnected() {
					continue
				}
			case <-c.connection.exit:
				c.Close()
				return
			}
			ticker.Stop()

			err := c.initialize()
			if err != nil {
				c.zapLog.Warn("[rabbitmq consumer] initialize consumer error", zap.String("err", err.Error()), zap.String("queue", c.QueueName))
				continue
			}

			delivery, err := c.consumeWithContext(ctx)
			if err != nil {
				c.zapLog.Warn("[rabbitmq consumer] execution of consumption error", zap.String("err", err.Error()), zap.String("queue", c.QueueName))
				continue
			}
			c.zapLog.Info("[rabbitmq consumer] queue is ready and waiting for messages, queue=" + c.QueueName)

			isContinueConsume := false
			for {
				select {
				case <-c.connection.exit:
					c.Close()
					return
				case d, ok := <-delivery:
					if !ok {
						c.zapLog.Warn("[rabbitmq consumer] exit consume message, queue=" + c.QueueName)
						isContinueConsume = true
						break
					}
					tagID := strings.Join([]string{d.Exchange, c.QueueName, strconv.FormatUint(d.DeliveryTag, 10)}, "/")
					err = handler(ctx, d.Body, tagID)
					if err != nil {
						c.zapLog.Warn("[rabbitmq consumer] handle message error", zap.String("err", err.Error()), zap.String("tagID", tagID))
						continue
					}
					if !c.isAutoAck {
						if err = d.Ack(false); err != nil {
							c.zapLog.Warn("[rabbitmq consumer] manual ack error", zap.String("err", err.Error()), zap.String("tagID", tagID))
							continue
						}
						c.zapLog.Info("[rabbitmq consumer] manual ack done", zap.String("tagID", tagID))
					}
				}

				if isContinueConsume {
					break
				}
			}
		}
	}()
}

// Close consumer
func (c *Consumer) Close() {
	if c.ch != nil {
		_ = c.ch.Close()
	}
}
