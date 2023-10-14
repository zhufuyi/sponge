package producer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQueueDeclareOptions(t *testing.T) {
	opts := []QueueDeclareOption{
		WithQueueDeclareDurable(true),
		WithQueueDeclareAutoDelete(true),
		WithQueueDeclareExclusive(true),
		WithQueueDeclareNoWait(true),
		WithQueueDeclareArgs(map[string]interface{}{"foo": "bar"}),
	}

	o := defaultQueueDeclareOptions()
	o.apply(opts...)

	assert.True(t, o.durable)
	assert.True(t, o.autoDelete)
	assert.True(t, o.exclusive)
	assert.True(t, o.noWait)
	assert.Equal(t, "bar", o.args["foo"])
}

func TestExchangeDeclareOptions(t *testing.T) {
	opts := []ExchangeDeclareOption{
		WithExchangeDeclareDurable(true),
		WithExchangeDeclareAutoDelete(true),
		WithExchangeDeclareInternal(true),
		WithExchangeDeclareNoWait(true),
		WithExchangeDeclareArgs(map[string]interface{}{"foo1": "bar1"}),
	}

	o := defaultExchangeDeclareOptions()
	o.apply(opts...)

	assert.True(t, o.durable)
	assert.True(t, o.autoDelete)
	assert.True(t, o.internal)
	assert.True(t, o.noWait)
	assert.Equal(t, "bar1", o.args["foo1"])
}

func TestQueueBindOptions(t *testing.T) {
	opts := []QueueBindOption{
		WithQueueBindNoWait(true),
		WithQueueBindArgs(map[string]interface{}{"foo2": "bar2"}),
	}

	o := defaultQueueBindOptions()
	o.apply(opts...)

	assert.True(t, o.noWait)
	assert.Equal(t, "bar2", o.args["foo2"])
}

func TestProducerOptions(t *testing.T) {
	opts := []QueueOption{
		WithQueueDeclareOptions(
			WithQueueDeclareDurable(true),
			WithQueueDeclareAutoDelete(true),
			WithQueueDeclareExclusive(true),
			WithQueueDeclareNoWait(true),
			WithQueueDeclareArgs(map[string]interface{}{"foo": "bar"}),
		),
		WithExchangeDeclareOptions(
			WithExchangeDeclareDurable(true),
			WithExchangeDeclareAutoDelete(true),
			WithExchangeDeclareInternal(true),
			WithExchangeDeclareNoWait(true),
			WithExchangeDeclareArgs(map[string]interface{}{"foo1": "bar1"}),
		),
		WithQueueBindOptions(
			WithQueueBindNoWait(true),
			WithQueueBindArgs(map[string]interface{}{"foo2": "bar2"}),
		),
		WithQueuePublishMandatory(true),
		WithQueuePublishImmediate(true)}

	o := defaultProducerOptions()
	o.apply(opts...)

	assert.True(t, o.queueDeclare.durable)
	assert.True(t, o.queueDeclare.autoDelete)
	assert.True(t, o.queueDeclare.exclusive)
	assert.True(t, o.queueDeclare.noWait)
	assert.Equal(t, "bar", o.queueDeclare.args["foo"])

	assert.True(t, o.exchangeDeclare.durable)
	assert.True(t, o.exchangeDeclare.autoDelete)
	assert.True(t, o.exchangeDeclare.internal)
	assert.True(t, o.exchangeDeclare.noWait)
	assert.Equal(t, "bar1", o.exchangeDeclare.args["foo1"])

	assert.True(t, o.queueBind.noWait)
	assert.Equal(t, "bar2", o.queueBind.args["foo2"])

	assert.True(t, o.mandatory)
	assert.True(t, o.immediate)
}
