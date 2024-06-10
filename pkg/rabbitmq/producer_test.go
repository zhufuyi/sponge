package rabbitmq

import (
	"context"
	"strconv"
	"testing"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/assert"

	"github.com/zhufuyi/sponge/pkg/utils"
)

func TestProducerOptions(t *testing.T) {
	opts := []ProducerOption{
		WithProducerExchangeDeclareOptions(
			WithExchangeDeclareAutoDelete(true),
			WithExchangeDeclareInternal(true),
			WithExchangeDeclareNoWait(true),
			WithExchangeDeclareArgs(map[string]interface{}{"foo1": "bar1"}),
		),
		WithProducerQueueDeclareOptions(
			WithQueueDeclareAutoDelete(true),
			WithQueueDeclareExclusive(true),
			WithQueueDeclareNoWait(true),
			WithQueueDeclareArgs(map[string]interface{}{"foo": "bar"}),
		),
		WithProducerQueueBindOptions(
			WithQueueBindNoWait(true),
			WithQueueBindArgs(map[string]interface{}{"foo2": "bar2"}),
		),
		WithProducerPersistent(true),
		WithProducerMandatory(true),

		WithDeadLetterOptions(WithDeadLetter("dl-exchange", "dl-queue", "dl-routing-key")),
	}

	o := defaultProducerOptions()
	o.apply(opts...)

	assert.True(t, o.queueDeclare.autoDelete)
	assert.True(t, o.queueDeclare.exclusive)
	assert.True(t, o.queueDeclare.noWait)
	assert.Equal(t, "bar", o.queueDeclare.args["foo"])

	assert.True(t, o.exchangeDeclare.autoDelete)
	assert.True(t, o.exchangeDeclare.internal)
	assert.True(t, o.exchangeDeclare.noWait)
	assert.Equal(t, "bar1", o.exchangeDeclare.args["foo1"])

	assert.True(t, o.queueBind.noWait)
	assert.Equal(t, "bar2", o.queueBind.args["foo2"])

	assert.True(t, o.isPersistent)
	assert.True(t, o.mandatory)
}

func TestProducer_direct(t *testing.T) {
	utils.SafeRunWithTimeout(time.Second*3, func(cancel context.CancelFunc) {
		defer cancel()
		connection, err := NewConnection(url)
		if err != nil {
			t.Log(err)
			return
		}
		defer connection.Close()

		ctx := context.Background()
		exchangeName := "direct-exchange-demo"
		queueName := "direct-queue-demo"
		routingKey := "info"
		exchange := NewDirectExchange(exchangeName, routingKey)
		p, err := NewProducer(exchange, queueName, connection)
		if err != nil {
			t.Log(err)
			return
		}
		defer p.Close()

		for i := 1; i <= 10; i++ {
			err = p.PublishDirect(ctx, []byte(routingKey+" say hello "+strconv.Itoa(i)))
			if err != nil {
				t.Error(err)
				return
			}
		}
	})
}

func TestProducer_topic(t *testing.T) {
	utils.SafeRunWithTimeout(time.Second*3, func(cancel context.CancelFunc) {
		defer cancel()
		connection, err := NewConnection(url)
		if err != nil {
			t.Log(err)
			return
		}
		defer connection.Close()

		ctx := context.Background()
		exchangeName := "topic-exchange-demo"

		queueName := "topic-queue-1"
		routingKey := "*.orange.*"
		exchange := NewTopicExchange(exchangeName, routingKey)
		p, err := NewProducer(exchange, queueName, connection)
		if err != nil {
			t.Log(err)
			return
		}
		defer p.Close()
		key := "key1.orange.key3"
		err = p.PublishTopic(ctx, key, []byte(key+" say hello"))

		queueName = "topic-queue-2"
		routingKey = "*.*.rabbit"
		exchange = NewTopicExchange(exchangeName, routingKey)
		p, err = NewProducer(exchange, queueName, connection)
		if err != nil {
			t.Log(err)
			return
		}
		defer p.Close()
		key = "key1.key2.rabbit"
		err = p.PublishTopic(ctx, key, []byte(key+" say hello"))

		queueName = "topic-queue-2"
		routingKey = "lazy.#"
		exchange = NewTopicExchange(exchangeName, routingKey)
		p, err = NewProducer(exchange, queueName, connection)
		if err != nil {
			t.Log(err)
			return
		}
		defer p.Close()
		key = "lazy.key2.key3"
		err = p.PublishTopic(ctx, key, []byte(key+" say hello"))
	})
}

func TestProducer_fanout(t *testing.T) {
	utils.SafeRunWithTimeout(time.Second*3, func(cancel context.CancelFunc) {
		defer cancel()
		connection, err := NewConnection(url)
		if err != nil {
			t.Log(err)
			return
		}
		defer connection.Close()

		ctx := context.Background()
		exchangeName := "fanout-exchange-demo"
		queueNames := []string{"fanout-queue-1", "fanout-queue-2", "fanout-queue-3"}

		for _, queueName := range queueNames {
			exchange := NewFanoutExchange(exchangeName)
			p, err := NewProducer(exchange, queueName, connection)
			if err != nil {
				t.Log(err)
				return
			}
			defer p.Close()
			err = p.PublishFanout(ctx, []byte(queueName+" say hello"))
			if err != nil {
				t.Error(err)
				return
			}
		}
	})
}

func TestProducer_headers(t *testing.T) {
	utils.SafeRunWithTimeout(time.Second*3, func(cancel context.CancelFunc) {
		defer cancel()
		connection, err := NewConnection(url)
		if err != nil {
			t.Log(err)
			return
		}
		defer connection.Close()

		ctx := context.Background()
		exchangeName := "headers-exchange-demo"

		// the message is only received if there is an exact match for headers
		queueName := "headers-queue-1"
		kv1 := map[string]interface{}{"hello1": "world1", "foo1": "bar1"}
		exchange := NewHeadersExchange(exchangeName, HeadersTypeAll, kv1)
		p, err := NewProducer(exchange, queueName, connection)
		if err != nil {
			t.Log(err)
			return
		}
		defer p.Close()
		headersKey1 := kv1 // exact match, consumer queue can receive messages
		err = p.PublishHeaders(ctx, headersKey1, []byte("say hello 1"))
		if err != nil {
			t.Error(err)
			return
		}
		headersKey1 = map[string]interface{}{"foo": "bar"} // there is a complete mismatch and the consumer queue cannot receive the message
		err = p.PublishHeaders(ctx, headersKey1, []byte("say hello 2"))
		if err != nil {
			t.Error(err)
			return
		}
		headersKey1 = map[string]interface{}{"foo1": "bar1"} // partial match, consumer queue cannot receive message
		err = p.PublishHeaders(ctx, headersKey1, []byte("say hello 3"))
		if err != nil {
			t.Error(err)
			return
		}

		// only partial matches of headers are needed to receive the message
		queueName = "headers-queue-2"
		kv2 := map[string]interface{}{"hello2": "world2", "foo2": "bar2"}
		exchange = NewHeadersExchange(exchangeName, HeadersTypeAny, kv2)
		p, err = NewProducer(exchange, queueName, connection)
		if err != nil {
			t.Error(err)
			return
		}
		defer p.Close()
		headersKey2 := kv2 // exact match, consumer queue can receive messages
		err = p.PublishHeaders(ctx, headersKey2, []byte("say hello 4"))
		if err != nil {
			t.Error(err)
			return
		}
		headersKey2 = map[string]interface{}{"foo": "bar"} // there is a complete mismatch and the consumer queue cannot receive the message
		err = p.PublishHeaders(ctx, headersKey2, []byte("say hello 5"))
		if err != nil {
			t.Error(err)
			return
		}
		headersKey2 = map[string]interface{}{"foo2": "bar2"} // partial match, the consumer queue can receive the message
		err = p.PublishHeaders(ctx, headersKey2, []byte("say hello 6"))
		if err != nil {
			t.Error(err)
			return
		}
	})
}

func TestProducer_delayedMessage(t *testing.T) {
	utils.SafeRunWithTimeout(time.Second*6, func(cancel context.CancelFunc) {
		defer cancel()
		connection, err := NewConnection(url)
		if err != nil {
			t.Log(err)
			return
		}
		defer connection.Close()

		ctx := context.Background()
		exchangeName := "delayed-message-exchange-demo"
		queueName := "delayed-message-queue"
		routingKey := "delayed-key"
		e := NewDirectExchange("", routingKey)
		exchange := NewDelayedMessageExchange(exchangeName, e)
		p, err := NewProducer(exchange, queueName, connection)
		if err != nil {
			t.Log(err)
			return
		}
		defer p.Close()

		for i := 0; i < 3; i++ {
			err = p.PublishDelayedMessage(ctx, time.Second*10, []byte("say hello "+time.Now().Format(datetimeLayout)))
			if err != nil {
				t.Error(err)
				return
			}
			time.Sleep(time.Second)
		}
	})
}

func TestPublishErr(t *testing.T) {
	p := &Producer{
		QueueName: "test",
		Exchange: &Exchange{
			name:       "test",
			eType:      "unknown",
			routingKey: "test",
		},
	}

	ctx := context.Background()
	err := p.PublishDirect(ctx, []byte("data"))
	assert.Error(t, err)
	err = p.PublishFanout(ctx, []byte("data"))
	assert.Error(t, err)
	err = p.PublishTopic(ctx, "", []byte("data"))
	assert.Error(t, err)
	err = p.PublishHeaders(ctx, nil, []byte("data"))
	assert.Error(t, err)
	err = p.PublishDelayedMessage(ctx, time.Second, []byte("data"))
	assert.Error(t, err)
}

func TestPublishDirect(t *testing.T) {
	p := &Producer{
		QueueName:    "foo",
		conn:         &amqp.Connection{},
		ch:           &amqp.Channel{},
		isPersistent: true,
		mandatory:    true,
	}
	defer func() { recover() }()
	ctx := context.Background()

	p.Exchange = NewDirectExchange("foo", "bar")
	_ = p.PublishDirect(ctx, []byte("data"))
}

func TestPublishTopic(t *testing.T) {
	p := &Producer{
		QueueName:    "foo",
		conn:         &amqp.Connection{},
		ch:           &amqp.Channel{},
		isPersistent: true,
		mandatory:    true,
	}
	defer func() { recover() }()
	ctx := context.Background()

	p.Exchange = NewDirectExchange("foo", "bar")
	_ = p.PublishTopic(ctx, "foo", []byte("data"))
	p.Exchange = NewTopicExchange("foo", "bar")
	_ = p.PublishTopic(ctx, "foo", []byte("data"))
}

func TestPublishFanout(t *testing.T) {
	p := &Producer{
		QueueName:    "foo",
		conn:         &amqp.Connection{},
		ch:           &amqp.Channel{},
		isPersistent: true,
		mandatory:    true,
	}
	defer func() { recover() }()
	ctx := context.Background()

	p.Exchange = NewFanoutExchange("foo")
	_ = p.PublishFanout(ctx, []byte("data"))
}

func TestPublishHeaders(t *testing.T) {
	p := &Producer{
		QueueName:    "foo",
		conn:         &amqp.Connection{},
		ch:           &amqp.Channel{},
		isPersistent: true,
		mandatory:    true,
	}
	defer func() { recover() }()
	ctx := context.Background()

	p.Exchange = NewDirectExchange("foo", "bar")
	_ = p.PublishHeaders(ctx, nil, []byte("data"))
	p.Exchange = NewHeadersExchange("foo", "bar", nil)
	_ = p.PublishHeaders(ctx, nil, []byte("data"))
}

func TestPublishDelayedMessage(t *testing.T) {
	p := &Producer{
		QueueName:    "foo",
		conn:         &amqp.Connection{},
		ch:           &amqp.Channel{},
		isPersistent: true,
		mandatory:    true,
	}
	defer func() { recover() }()
	ctx := context.Background()

	p.Exchange = NewDelayedMessageExchange("foo", NewTopicExchange("", "bar"))
	_ = p.PublishDelayedMessage(ctx, time.Second, []byte("data"))
	p.Exchange = NewDelayedMessageExchange("foo", NewHeadersExchange("", HeadersTypeAll, nil))
	_ = p.PublishDelayedMessage(ctx, time.Second, []byte("data"))
	p.Exchange = NewDelayedMessageExchange("foo", NewDirectExchange("", "bar"))
	_ = p.PublishDelayedMessage(ctx, time.Second, []byte("data"))
}

func TestProducerErr(t *testing.T) {
	exchangeName := "direct-exchange-demo"
	queueName := "direct-queue-1"
	routeKey := "direct-key-1"
	exchange := NewDirectExchange(exchangeName, routeKey)

	utils.SafeRunWithTimeout(time.Second, func(cancel context.CancelFunc) {
		defer cancel()
		_, err := NewProducer(exchange, queueName, &Connection{conn: &amqp.Connection{}})
		if err != nil {
			t.Log(err)
			return
		}
	})

	p := &Producer{conn: &amqp.Connection{}, ch: &amqp.Channel{}}
	utils.SafeRun(context.Background(), func(ctx context.Context) {
		_ = p.PublishDirect(context.Background(), []byte("hello world"))
	})
	utils.SafeRun(context.Background(), func(ctx context.Context) {
		p.Close()
	})
}

func Test_printFields(t *testing.T) {
	exchange := NewDirectExchange("foo", "bar")
	fields := logFields("queue", exchange)
	t.Log(fields)

	exchange = NewHeadersExchange("foo", HeadersTypeAny, map[string]interface{}{"hello": "world"})
	fields = logFields("queue", exchange)
	t.Log(fields)

	e := NewDirectExchange("", "bar")
	exchange = NewDelayedMessageExchange("foo", e)
	fields = logFields("queue", exchange)
	t.Log(fields)

	e = NewHeadersExchange("", HeadersTypeAny, map[string]interface{}{"hello": "world"})
	exchange = NewDelayedMessageExchange("foo", e)
	fields = logFields("queue", exchange)
	t.Log(fields)
}
