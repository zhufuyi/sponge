package rabbitmq

import amqp "github.com/rabbitmq/amqp091-go"

// ErrClosed closed
var ErrClosed = amqp.ErrClosed

const (
	exchangeTypeDirect         = "direct"
	exchangeTypeTopic          = "topic"
	exchangeTypeFanout         = "fanout"
	exchangeTypeHeaders        = "headers"
	exchangeTypeDelayedMessage = "x-delayed-message"

	// HeadersTypeAll all
	HeadersTypeAll HeadersType = "all"
	// HeadersTypeAny any
	HeadersTypeAny HeadersType = "any"
)

// HeadersType headers type
type HeadersType = string

// Exchange rabbitmq minimum management unit
type Exchange struct {
	name               string                 // exchange name
	eType              string                 // exchange type: direct, topic, fanout, headers, x-delayed-message
	routingKey         string                 // route key
	headersKeys        map[string]interface{} // this field is required if eType=headers.
	delayedMessageType string                 // this field is required if eType=headers, support direct, topic, fanout, headers
}

// Name exchange name
func (e *Exchange) Name() string {
	return e.name
}

// Type exchange type
func (e *Exchange) Type() string {
	return e.eType
}

// RoutingKey exchange routing key
func (e *Exchange) RoutingKey() string {
	return e.routingKey
}

// HeadersKeys exchange headers keys
func (e *Exchange) HeadersKeys() map[string]interface{} {
	return e.headersKeys
}

// DelayedMessageType exchange delayed message type
func (e *Exchange) DelayedMessageType() string {
	return e.delayedMessageType
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

// NewFanoutExchange create a fanout exchange
func NewFanoutExchange(exchangeName string) *Exchange {
	return &Exchange{
		name:       exchangeName,
		eType:      exchangeTypeFanout,
		routingKey: "",
	}
}

// NewHeadersExchange create a headers exchange, the headerType supports "all" and "any"
func NewHeadersExchange(exchangeName string, headersType HeadersType, keys map[string]interface{}) *Exchange {
	if keys == nil {
		keys = make(map[string]interface{})
	}

	switch headersType {
	case HeadersTypeAll, HeadersTypeAny:
		keys["x-match"] = headersType
	default:
		keys["x-match"] = HeadersTypeAll
	}

	return &Exchange{
		name:        exchangeName,
		eType:       exchangeTypeHeaders,
		routingKey:  "",
		headersKeys: keys,
	}
}

// NewDelayedMessageExchange create a delayed message exchange
func NewDelayedMessageExchange(exchangeName string, e *Exchange) *Exchange {
	return &Exchange{
		name:               exchangeName,
		eType:              "x-delayed-message",
		routingKey:         e.routingKey,
		delayedMessageType: e.eType,
		headersKeys:        e.headersKeys,
	}
}

// -------------------------------------------------------------------------------------------

// QueueDeclareOption declare queue option.
type QueueDeclareOption func(*queueDeclareOptions)

type queueDeclareOptions struct {
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
		autoDelete: false,
		exclusive:  false,
		noWait:     false,
		args:       nil,
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
		//durable:    true,
		autoDelete: false,
		internal:   false,
		noWait:     false,
		args:       nil,
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

// DelayedMessagePublishOption declare queue bind option.
type DelayedMessagePublishOption func(*delayedMessagePublishOptions)

type delayedMessagePublishOptions struct {
	topicKey    string                 // the topic message type must be required
	headersKeys map[string]interface{} // the headers message type must be required
}

func (o *delayedMessagePublishOptions) apply(opts ...DelayedMessagePublishOption) {
	for _, opt := range opts {
		opt(o)
	}
}

// default delayed message publish settings
func defaultDelayedMessagePublishOptions() *delayedMessagePublishOptions {
	return &delayedMessagePublishOptions{}
}

// WithDelayedMessagePublishTopicKey set delayed message publish topicKey option.
func WithDelayedMessagePublishTopicKey(topicKey string) DelayedMessagePublishOption {
	return func(o *delayedMessagePublishOptions) {
		if topicKey == "" {
			return
		}
		o.topicKey = topicKey
	}
}

// WithDelayedMessagePublishHeadersKeys set delayed message publish headersKeys option.
func WithDelayedMessagePublishHeadersKeys(headersKeys map[string]interface{}) DelayedMessagePublishOption {
	return func(o *delayedMessagePublishOptions) {
		if headersKeys == nil {
			return
		}
		o.headersKeys = headersKeys
	}
}

// -------------------------------------------------------------------------------------------

// DeadLetterOption declare dead letter option.
type DeadLetterOption func(*deadLetterOptions)

type deadLetterOptions struct {
	exchangeName string
	queueName    string
	routingKey   string

	exchangeDeclare *exchangeDeclareOptions
	queueDeclare    *queueDeclareOptions
	queueBind       *queueBindOptions
}

func (o *deadLetterOptions) apply(opts ...DeadLetterOption) {
	for _, opt := range opts {
		opt(o)
	}
}

func (o *deadLetterOptions) isEnabled() bool {
	if o.exchangeName != "" && o.queueName != "" {
		return true
	}
	return false
}

func defaultDeadLetterOptions() *deadLetterOptions {
	return &deadLetterOptions{
		exchangeDeclare: defaultExchangeDeclareOptions(),
		queueDeclare:    defaultQueueDeclareOptions(),
		queueBind:       defaultQueueBindOptions(),
	}
}

// WithDeadLetterExchangeDeclareOptions set dead letter exchange declare option.
func WithDeadLetterExchangeDeclareOptions(opts ...ExchangeDeclareOption) DeadLetterOption {
	return func(o *deadLetterOptions) {
		o.exchangeDeclare.apply(opts...)
	}
}

// WithDeadLetterQueueDeclareOptions set dead letter queue declare option.
func WithDeadLetterQueueDeclareOptions(opts ...QueueDeclareOption) DeadLetterOption {
	return func(o *deadLetterOptions) {
		o.queueDeclare.apply(opts...)
	}
}

// WithDeadLetterQueueBindOptions set dead letter queue bind option.
func WithDeadLetterQueueBindOptions(opts ...QueueBindOption) DeadLetterOption {
	return func(o *deadLetterOptions) {
		o.queueBind.apply(opts...)
	}
}

// WithDeadLetter set dead letter exchange, queue, routing key.
func WithDeadLetter(exchangeName string, queueName string, routingKey string) DeadLetterOption {
	return func(o *deadLetterOptions) {
		o.exchangeName = exchangeName
		o.queueName = queueName
		o.routingKey = routingKey
	}
}
