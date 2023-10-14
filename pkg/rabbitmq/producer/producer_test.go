package producer

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/zhufuyi/sponge/pkg/rabbitmq"
	"github.com/zhufuyi/sponge/pkg/utils"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/assert"
)

const (
	url = "amqp://guest:guest@192.168.3.37:5672/"
)

func TestProducer_direct(t *testing.T) {
	utils.SafeRunWithTimeout(time.Second*3, func(cancel context.CancelFunc) {

		c, err := rabbitmq.NewConnection(url)
		if err != nil {
			t.Log(err)
			return
		}
		defer c.Close()

		ctx := context.Background()
		exchangeName := "direct-exchange-demo"
		queueName := "direct-queue-1"
		routeKey := "direct-key-1"
		exchange := NewDirectExchange(exchangeName, routeKey)
		q, err := NewQueue(queueName, c.Conn, exchange)
		if err != nil {
			t.Log(err)
			return
		}
		defer q.Close()

		for i := 0; i < 10; i++ {
			err = q.Publish(ctx, []byte(routeKey+" say hello "+strconv.Itoa(i)))
			if err != nil {
				t.Error(err)
				return
			}
		}

	})
}

func TestProducer_topic(t *testing.T) {
	utils.SafeRunWithTimeout(time.Second*3, func(cancel context.CancelFunc) {

		c, err := rabbitmq.NewConnection(url)
		if err != nil {
			t.Log(err)
			return
		}
		defer c.Close()

		ctx := context.Background()
		exchangeName := "topic-exchange-demo"

		queueName := "topic-queue-1"
		routingKey := "key1.key2.*"
		exchange := NewTopicExchange(exchangeName, routingKey)
		q, err := NewQueue(queueName, c.Conn, exchange)
		if err != nil {
			t.Log(err)
			return
		}
		defer q.Close()

		queueName = "topic-queue-2"
		routingKey = "*.key2"
		exchange = NewTopicExchange(exchangeName, routingKey)
		q, err = NewQueue(queueName, c.Conn, exchange)
		if err != nil {
			t.Log(err)
			return
		}
		defer q.Close()

		queueName = "topic-queue-3"
		routingKey = "key1.#"
		exchange = NewTopicExchange(exchangeName, routingKey)
		q, err = NewQueue(queueName, c.Conn, exchange)
		if err != nil {
			t.Log(err)
			return
		}
		defer q.Close()

		queueName = "topic-queue-4"
		routingKey = "#.key3"
		exchange = NewTopicExchange(exchangeName, routingKey)
		q, err = NewQueue(queueName, c.Conn, exchange)
		if err != nil {
			t.Log(err)
			return
		}
		defer q.Close()

		keys := []string{
			"key1",           // only match queue 3
			"key1.key2",      // only match queue 2 and 3
			"key2.key3",      // only match queue 4
			"key1.key2.key3", // match queue 1,2,3,4
		}
		for _, key := range keys {
			err = q.PublishTopic(ctx, key, []byte(key+" say hello "))
			if err != nil {
				t.Error(err)
				return
			}
		}

	})
}

func TestProducer_fanout(t *testing.T) {
	utils.SafeRunWithTimeout(time.Second*3, func(cancel context.CancelFunc) {

		c, err := rabbitmq.NewConnection(url)
		if err != nil {
			t.Log(err)
			return
		}
		defer c.Close()

		ctx := context.Background()
		exchangeName := "fanout-exchange-demo"

		queueName := "fanout-queue-1"
		exchange := NewFanOutExchange(exchangeName)
		q, err := NewQueue(queueName, c.Conn, exchange)
		if err != nil {
			t.Log(err)
			return
		}
		defer q.Close()

		queueName = "fanout-queue-2"
		exchange = NewFanOutExchange(exchangeName)
		q, err = NewQueue(queueName, c.Conn, exchange)
		if err != nil {
			t.Log(err)
			return
		}
		defer q.Close()

		queueName = "fanout-queue-3"
		exchange = NewFanOutExchange(exchangeName)
		q, err = NewQueue(queueName, c.Conn, exchange)
		if err != nil {
			t.Log(err)
			return
		}
		defer q.Close()

		// queues 1,2 and 3 can receive the same messages.
		for i := 0; i < 10; i++ {
			err = q.Publish(ctx, []byte(" say hello "+strconv.Itoa(i)))
			if err != nil {
				t.Error(err)
				return
			}
		}

	})
}

func TestProducer_headers(t *testing.T) {
	utils.SafeRunWithTimeout(time.Second*3, func(cancel context.CancelFunc) {

		c, err := rabbitmq.NewConnection(url)
		if err != nil {
			t.Log(err)
			return
		}
		defer c.Close()

		ctx := context.Background()
		exchangeName := "headers-exchange-demo"

		// the message is only received if there is an exact match for headers
		queueName := "headers-queue-1"
		kv1 := map[string]interface{}{"hello1": "world1", "foo1": "bar1"}
		exchange := NewHeaderExchange(exchangeName, HeadersTypeAll, kv1)
		q, err := NewQueue(queueName, c.Conn, exchange)
		if err != nil {
			t.Log(err)
			return
		}
		defer q.Close()
		headersKey1 := kv1 // exact match, consumer queue can receive messages
		err = q.PublishHeaders(ctx, headersKey1, []byte("say hello 1"))
		if err != nil {
			t.Error(err)
			return
		}
		headersKey1 = map[string]interface{}{"foo": "bar"} // there is a complete mismatch and the consumer queue cannot receive the message
		err = q.PublishHeaders(ctx, headersKey1, []byte("say hello 2"))
		if err != nil {
			t.Error(err)
			return
		}
		headersKey1 = map[string]interface{}{"foo1": "bar1"} // partial match, consumer queue cannot receive message
		err = q.PublishHeaders(ctx, headersKey1, []byte("say hello 3"))
		if err != nil {
			t.Error(err)
			return
		}

		// only partial matches of headers are needed to receive the message
		queueName = "headers-queue-2"
		kv2 := map[string]interface{}{"hello2": "world2", "foo2": "bar2"}
		exchange = NewHeaderExchange(exchangeName, HeadersTypeAny, kv2)
		q, err = NewQueue(queueName, c.Conn, exchange)
		if err != nil {
			t.Error(err)
			return
		}
		defer q.Close()
		headersKey2 := kv2 // exact match, consumer queue can receive messages
		err = q.PublishHeaders(ctx, headersKey2, []byte("say hello 4"))
		if err != nil {
			t.Error(err)
			return
		}
		headersKey2 = map[string]interface{}{"foo": "bar"} // there is a complete mismatch and the consumer queue cannot receive the message
		err = q.PublishHeaders(ctx, headersKey2, []byte("say hello 5"))
		if err != nil {
			t.Error(err)
			return
		}
		headersKey2 = map[string]interface{}{"foo2": "bar2"} // partial match, the consumer queue can receive the message
		err = q.PublishHeaders(ctx, headersKey2, []byte("say hello 6"))
		if err != nil {
			t.Error(err)
			return
		}

	})
}

func TestQueueErr(t *testing.T) {
	q := &Queue{
		queueName: "test",
		exchange: &Exchange{
			name:       "test",
			eType:      "unknown",
			routingKey: "test",
		},
		//queue: amqp.Queue{},
	}

	ctx := context.Background()
	err := q.Publish(ctx, []byte("test"))
	assert.Error(t, err)
	err = q.PublishTopic(ctx, "", []byte("test"))
	assert.Error(t, err)
	err = q.PublishHeaders(ctx, nil, []byte("test"))
	assert.Error(t, err)
}

func TestNewExchange(t *testing.T) {
	NewDirectExchange("foo", "bar")
	NewTopicExchange("foo", "bar")
	NewFanOutExchange("foo")
	NewHeaderExchange("foo", HeadersTypeAll, nil)
	NewHeaderExchange("foo", "bar", nil)
}

func TestNewQueue(t *testing.T) {
	defer func() { recover() }()
	q, err := NewQueue("foo", &amqp.Connection{}, NewDirectExchange("foo", "bar"))
	if err != nil {
		t.Log(err)
		return
	}
	q.Close()
}

func TestPublish(t *testing.T) {
	q := Queue{
		queueName: "foo",
		conn:      &amqp.Connection{},
		ch:        &amqp.Channel{},
		mandatory: false,
		immediate: false,
	}
	defer func() { recover() }()

	q.exchange = NewTopicExchange("foo", "bar")
	_ = q.Publish(context.Background(), []byte("test"))
	q.exchange = NewDirectExchange("foo", "bar")

	_ = q.Publish(context.Background(), []byte("test"))
}

func TestPublishTopic(t *testing.T) {
	q := Queue{
		queueName: "foo",
		conn:      &amqp.Connection{},
		ch:        &amqp.Channel{},
		mandatory: false,
		immediate: false,
	}
	defer func() { recover() }()

	q.exchange = NewDirectExchange("foo", "bar")
	_ = q.PublishTopic(context.Background(), "foo", []byte("bar"))
	q.exchange = NewTopicExchange("foo", "bar")
	_ = q.PublishTopic(context.Background(), "foo", []byte("bar"))
}

func TestPublishHeaders(t *testing.T) {
	q := Queue{
		queueName: "foo",
		conn:      &amqp.Connection{},
		ch:        &amqp.Channel{},
		mandatory: false,
		immediate: false,
	}
	defer func() { recover() }()

	q.exchange = NewDirectExchange("foo", "bar")
	_ = q.PublishHeaders(context.Background(), nil, []byte("bar"))
	q.exchange = NewHeaderExchange("foo", "bar", nil)
	_ = q.PublishHeaders(context.Background(), nil, []byte("bar"))
}
