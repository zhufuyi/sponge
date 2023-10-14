package producer

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

// QueueDeclareOption declare queue option.
type QueueDeclareOption func(*queueDeclareOptions)

type queueDeclareOptions struct {
	durable    bool       // is it persistent
	autoDelete bool       // delete automatically
	exclusive  bool       // exclusive (only available to the program that created it)
	noWait     bool       // block processing
	args       amqp.Table // additional properties
}

func (o *queueDeclareOptions) apply(opts ...QueueDeclareOption) {
	for _, opt := range opts {
		opt(o)
	}
}

// default queue declare settings
func defaultQueueDeclareOptions() *queueDeclareOptions {
	return &queueDeclareOptions{
		durable:    true,
		autoDelete: false,
		exclusive:  false,
		noWait:     false,
		args:       nil,
	}
}

// WithQueueDeclareDurable set queue declare durable option.
func WithQueueDeclareDurable(enable bool) QueueDeclareOption {
	return func(o *queueDeclareOptions) {
		o.durable = enable
	}
}

// WithQueueDeclareAutoDelete set queue declare auto delete option.
func WithQueueDeclareAutoDelete(enable bool) QueueDeclareOption {
	return func(o *queueDeclareOptions) {
		o.autoDelete = enable
	}
}

// WithQueueDeclareExclusive set queue declare exclusive option.
func WithQueueDeclareExclusive(enable bool) QueueDeclareOption {
	return func(o *queueDeclareOptions) {
		o.exclusive = enable
	}
}

// WithQueueDeclareNoWait set queue declare no wait option.
func WithQueueDeclareNoWait(enable bool) QueueDeclareOption {
	return func(o *queueDeclareOptions) {
		o.noWait = enable
	}
}

// WithQueueDeclareArgs set queue declare args option.
func WithQueueDeclareArgs(args map[string]interface{}) QueueDeclareOption {
	return func(o *queueDeclareOptions) {
		o.args = args
	}
}

// -------------------------------------------------------------------------------------------

// ExchangeDeclareOption declare exchange option.
type ExchangeDeclareOption func(*exchangeDeclareOptions)

type exchangeDeclareOptions struct {
	durable    bool       // is it persistent
	autoDelete bool       // delete automatically
	internal   bool       // public or not, false means public
	noWait     bool       // block processing
	args       amqp.Table // additional properties
}

func (o *exchangeDeclareOptions) apply(opts ...ExchangeDeclareOption) {
	for _, opt := range opts {
		opt(o)
	}
}

// default exchange declare settings
func defaultExchangeDeclareOptions() *exchangeDeclareOptions {
	return &exchangeDeclareOptions{
		durable:    true,
		autoDelete: false,
		internal:   false,
		noWait:     false,
		args:       nil,
	}
}

// WithExchangeDeclareDurable set exchange declare durable option.
func WithExchangeDeclareDurable(enable bool) ExchangeDeclareOption {
	return func(o *exchangeDeclareOptions) {
		o.durable = enable
	}
}

// WithExchangeDeclareAutoDelete set exchange declare auto delete option.
func WithExchangeDeclareAutoDelete(enable bool) ExchangeDeclareOption {
	return func(o *exchangeDeclareOptions) {
		o.autoDelete = enable
	}
}

// WithExchangeDeclareInternal set exchange declare internal option.
func WithExchangeDeclareInternal(enable bool) ExchangeDeclareOption {
	return func(o *exchangeDeclareOptions) {
		o.internal = enable
	}
}

// WithExchangeDeclareNoWait set exchange declare no wait option.
func WithExchangeDeclareNoWait(enable bool) ExchangeDeclareOption {
	return func(o *exchangeDeclareOptions) {
		o.noWait = enable
	}
}

// WithExchangeDeclareArgs set exchange declare args option.
func WithExchangeDeclareArgs(args map[string]interface{}) ExchangeDeclareOption {
	return func(o *exchangeDeclareOptions) {
		o.args = args
	}
}

// -------------------------------------------------------------------------------------------

// QueueBindOption declare queue bind option.
type QueueBindOption func(*queueBindOptions)

type queueBindOptions struct {
	noWait bool       // block processing
	args   amqp.Table // this parameter is invalid if the type is headers.
}

func (o *queueBindOptions) apply(opts ...QueueBindOption) {
	for _, opt := range opts {
		opt(o)
	}
}

// default queue bind settings
func defaultQueueBindOptions() *queueBindOptions {
	return &queueBindOptions{
		noWait: false,
		args:   nil,
	}
}

// WithQueueBindNoWait set queue bind no wait option.
func WithQueueBindNoWait(enable bool) QueueBindOption {
	return func(o *queueBindOptions) {
		o.noWait = enable
	}
}

// WithQueueBindArgs set queue bind args option.
func WithQueueBindArgs(args map[string]interface{}) QueueBindOption {
	return func(o *queueBindOptions) {
		o.args = args
	}
}

// -------------------------------------------------------------------------------------------

// QueueOption queue option.
type QueueOption func(*queueOptions)

type queueOptions struct {
	queueDeclare    *queueDeclareOptions
	exchangeDeclare *exchangeDeclareOptions
	queueBind       *queueBindOptions

	// If true, the message will be returned to the sender if the queue cannot be
	// found according to its own exchange type and routeKey rules.
	mandatory bool
	// If true, when exchange sends a message to the queue and finds that there
	// are no consumers on the queue, it returns the message to the sender
	immediate bool
}

func (o *queueOptions) apply(opts ...QueueOption) {
	for _, opt := range opts {
		opt(o)
	}
}

// default queue declare settings
func defaultProducerOptions() *queueOptions {
	return &queueOptions{
		queueDeclare:    defaultQueueDeclareOptions(),
		exchangeDeclare: defaultExchangeDeclareOptions(),
		queueBind:       defaultQueueBindOptions(),

		mandatory: false,
		immediate: false,
	}
}

// WithQueueDeclareOptions set queue declare option.
func WithQueueDeclareOptions(opts ...QueueDeclareOption) QueueOption {
	return func(o *queueOptions) {
		o.queueDeclare.apply(opts...)
	}
}

// WithExchangeDeclareOptions set exchange declare option.
func WithExchangeDeclareOptions(opts ...ExchangeDeclareOption) QueueOption {
	return func(o *queueOptions) {
		o.exchangeDeclare.apply(opts...)
	}
}

// WithQueueBindOptions set queue bind option.
func WithQueueBindOptions(opts ...QueueBindOption) QueueOption {
	return func(o *queueOptions) {
		o.queueBind.apply(opts...)
	}
}

// WithQueuePublishMandatory set queue publish mandatory option.
func WithQueuePublishMandatory(enable bool) QueueOption {
	return func(o *queueOptions) {
		o.mandatory = enable
	}
}

// WithQueuePublishImmediate set queue publish immediate option.
func WithQueuePublishImmediate(enable bool) QueueOption {
	return func(o *queueOptions) {
		o.immediate = enable
	}
}
