package rabbitmq

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExchange(t *testing.T) {
	e := NewDirectExchange("foo", "bar")
	assert.Equal(t, e.eType, exchangeTypeDirect)
	e = NewTopicExchange("foo", "bar.*")
	assert.Equal(t, e.eType, exchangeTypeTopic)
	e = NewFanoutExchange("foo")
	assert.Equal(t, e.eType, exchangeTypeFanout)
	e = NewHeadersExchange("foo", HeadersTypeAll, nil)
	assert.Equal(t, e.eType, exchangeTypeHeaders)
	e = NewHeadersExchange("foo", "unknown", nil)
	assert.Equal(t, e.eType, exchangeTypeHeaders)
	e = NewDelayedMessageExchange("foobar", NewDirectExchange("", "key"))
	assert.Equal(t, e.eType, exchangeTypeDelayedMessage)

	e = NewDelayedMessageExchange("foobar", NewDirectExchange("", "key"))
	assert.Equal(t, e.name, e.Name())
	assert.Equal(t, e.eType, e.Type())
	assert.Equal(t, e.routingKey, e.RoutingKey())
	assert.Equal(t, e.delayedMessageType, e.DelayedMessageType())
	assert.Equal(t, e.headersKeys, e.HeadersKeys())
}

func TestExchangeDeclareOptions(t *testing.T) {
	opts := []ExchangeDeclareOption{
		WithExchangeDeclareAutoDelete(true),
		WithExchangeDeclareInternal(true),
		WithExchangeDeclareNoWait(true),
		WithExchangeDeclareArgs(map[string]interface{}{"foo": "bar"}),
	}

	o := defaultExchangeDeclareOptions()
	o.apply(opts...)

	assert.True(t, o.autoDelete)
	assert.True(t, o.internal)
	assert.True(t, o.noWait)
	assert.Equal(t, "bar", o.args["foo"])
}

func TestQueueDeclareOptions(t *testing.T) {
	opts := []QueueDeclareOption{
		WithQueueDeclareAutoDelete(true),
		WithQueueDeclareExclusive(true),
		WithQueueDeclareNoWait(true),
		WithQueueDeclareArgs(map[string]interface{}{"foo": "bar"}),
	}

	o := defaultQueueDeclareOptions()
	o.apply(opts...)

	assert.True(t, o.autoDelete)
	assert.True(t, o.exclusive)
	assert.True(t, o.noWait)
	assert.Equal(t, "bar", o.args["foo"])
}

func TestQueueBindOptions(t *testing.T) {
	opts := []QueueBindOption{
		WithQueueBindNoWait(true),
		WithQueueBindArgs(map[string]interface{}{"foo": "bar"}),
	}

	o := defaultQueueBindOptions()
	o.apply(opts...)

	assert.True(t, o.noWait)
	assert.Equal(t, "bar", o.args["foo"])
}

func TestDelayedMessagePublishOptions(t *testing.T) {
	opts := []DelayedMessagePublishOption{
		WithDelayedMessagePublishTopicKey(""),
		WithDelayedMessagePublishTopicKey("key1.key2"),
		WithDelayedMessagePublishHeadersKeys(nil),
		WithDelayedMessagePublishHeadersKeys(map[string]interface{}{"foo": "bar"}),
	}

	o := defaultDelayedMessagePublishOptions()
	o.apply(opts...)

	assert.Equal(t, "key1.key2", o.topicKey)
	assert.Equal(t, "bar", o.headersKeys["foo"])
}
