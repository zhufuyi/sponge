package consumer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConsumeOptions(t *testing.T) {
	opts := []ConsumeOption{
		WithConsumeConsumer("test"),
		WithConsumeAutoAck(true),
		WithConsumeExclusive(true),
		WithConsumeNoLocal(true),
		WithConsumeNoWait(true),
		WithConsumeArgs(map[string]interface{}{"foo": "bar"}),
		WithConsumeQos(
			WithQosPrefetchCount(1),
			WithQosPrefetchSize(4096),
			WithQosPrefetchGlobal(true),
		),
	}

	o := defaultConsumeOptions()
	o.apply(opts...)

	assert.Equal(t, "test", o.consumer)
	assert.True(t, o.autoAck)
	assert.True(t, o.exclusive)
	assert.True(t, o.noLocal)
	assert.True(t, o.noWait)
	assert.Equal(t, "bar", o.args["foo"])
	assert.True(t, o.enableQos)
	assert.Equal(t, 1, o.qos.prefetchCount)
	assert.Equal(t, 4096, o.qos.prefetchSize)
	assert.True(t, o.qos.global)
}

func TestQosOptions(t *testing.T) {
	opts := []QosOption{
		WithQosPrefetchCount(1),
		WithQosPrefetchSize(4096),
		WithQosPrefetchGlobal(true),
	}

	o := defaultQosOptions()
	o.apply(opts...)

	assert.Equal(t, 1, o.prefetchCount)
	assert.Equal(t, 4096, o.prefetchSize)
	assert.True(t, o.global)
}
