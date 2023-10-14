package consumer

import amqp "github.com/rabbitmq/amqp091-go"

// ConsumeOption consume option.
type ConsumeOption func(*consumeOptions)

type consumeOptions struct {
	consumer  string     // used to distinguish between multiple consumers
	autoAck   bool       // auto-answer or not, if false, manual ACK required
	exclusive bool       // only available to the program that created it
	noLocal   bool       // if set to true, a message sent by a producer in the same Connection cannot be passed to a consumer in this Connection.
	noWait    bool       // block processing
	args      amqp.Table // additional properties

	enableQos bool
	qos       *qosOptions
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
		autoAck:   true,
		exclusive: false,
		noLocal:   false,
		noWait:    false,
		args:      nil,
		enableQos: false,
		qos:       defaultQosOptions(),
	}
}

// WithConsumeConsumer set consume consumer option.
func WithConsumeConsumer(consumer string) ConsumeOption {
	return func(o *consumeOptions) {
		o.consumer = consumer
	}
}

// WithConsumeAutoAck set consume auto ack option.
func WithConsumeAutoAck(enable bool) ConsumeOption {
	return func(o *consumeOptions) {
		o.autoAck = enable
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

// WithConsumeQos set consume qos option.
func WithConsumeQos(opts ...QosOption) ConsumeOption {
	return func(o *consumeOptions) {
		o.enableQos = true
		o.qos.apply(opts...)
	}
}

// -------------------------------------------------------------------------------------------

// QosOption qos option.
type QosOption func(*qosOptions)

type qosOptions struct {
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
		prefetchCount: 0,
		prefetchSize:  0,
		global:        false,
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
