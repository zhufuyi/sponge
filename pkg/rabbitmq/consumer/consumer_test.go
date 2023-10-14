package consumer

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/zhufuyi/sponge/pkg/rabbitmq"
	"github.com/zhufuyi/sponge/pkg/rabbitmq/producer"
	"github.com/zhufuyi/sponge/pkg/utils"
)

var url = "amqp://guest:guest@192.168.3.37:5672/"

var handler = func(ctx context.Context, data []byte, tag ...string) error {
	tagID := strings.Join(tag, ",")
	fmt.Printf("tagID=%s, receive message: %s\n", tagID, string(data))
	return nil
}

func consume(ctx context.Context, queueNames ...string) error {
	var consumeErr error
	utils.SafeRunWithTimeout(time.Second*3, func(cancel context.CancelFunc) {

		c, err := rabbitmq.NewConnection(url)
		if err != nil {
			consumeErr = err
			return
		}

		for _, queueName := range queueNames {
			queue, err := NewQueue(ctx, queueName, c, WithConsumeAutoAck(false))
			if err != nil {
				consumeErr = err
				return
			}
			queue.Consume(handler)
		}

	})
	return consumeErr
}

func TestConsumer_direct(t *testing.T) {
	queueName := "direct-queue-1"

	err := producerDirect(queueName)
	if err != nil {
		t.Log(err)
		return
	}

	ctx, _ := context.WithTimeout(context.Background(), time.Second*3)
	err = consume(ctx, queueName)
	if err != nil {
		t.Log(err)
		return
	}
	<-ctx.Done()
	time.Sleep(time.Millisecond * 100)
}

func TestConsumer_topic(t *testing.T) {
	queueName := "topic-queue-1"

	err := producerTopic(queueName)
	if err != nil {
		t.Log(err)
		return
	}

	ctx, _ := context.WithTimeout(context.Background(), time.Second*3)
	err = consume(ctx, queueName)
	if err != nil {
		t.Log(err)
		return
	}
	<-ctx.Done()
	time.Sleep(time.Millisecond * 100)
}

func TestConsumer_fanout(t *testing.T) {
	queueName := "fanout-queue-1"
	err := producerFanout(queueName)
	if err != nil {
		t.Log(err)
		return
	}

	ctx, _ := context.WithTimeout(context.Background(), time.Second*3)
	err = consume(ctx, queueName)
	if err != nil {
		t.Log(err)
		return
	}
	<-ctx.Done()
	time.Sleep(time.Millisecond * 100)
}

func TestConsumer_headers(t *testing.T) {
	queueName := "headers-queue-1"
	err := producerHeaders(queueName)
	if err != nil {
		t.Log(err)
		return
	}

	ctx, _ := context.WithTimeout(context.Background(), time.Second*3)
	err = consume(ctx, queueName)
	if err != nil {
		t.Log(err)
		return
	}
	<-ctx.Done()
	time.Sleep(time.Millisecond * 100)
}

func producerDirect(queueName string) error {
	var producerErr error
	utils.SafeRunWithTimeout(time.Second*3, func(cancel context.CancelFunc) {

		c, err := rabbitmq.NewConnection(url)
		if err != nil {
			producerErr = err
			return
		}
		defer c.Close()

		ctx := context.Background()
		exchangeName := "direct-exchange-demo"
		routeKey := "direct-key-1"
		exchange := producer.NewDirectExchange(exchangeName, routeKey)
		q, err := producer.NewQueue(queueName, c.Conn, exchange)
		if err != nil {
			producerErr = err
			return
		}
		defer q.Close()

		_ = q.Publish(ctx, []byte("say hello 1"))
		_ = q.Publish(ctx, []byte("say hello 2"))
		producerErr = q.Publish(ctx, []byte("say hello 3"))

	})

	return producerErr
}

func producerTopic(queueName string) error {
	var producerErr error
	utils.SafeRunWithTimeout(time.Second*3, func(cancel context.CancelFunc) {

		c, err := rabbitmq.NewConnection(url)
		if err != nil {
			producerErr = err
			return
		}
		defer c.Close()

		ctx := context.Background()
		exchangeName := "topic-exchange-demo"

		routingKey := "key1.key2.*"
		exchange := producer.NewTopicExchange(exchangeName, routingKey)
		q, err := producer.NewQueue(queueName, c.Conn, exchange)
		if err != nil {
			producerErr = err
			return
		}
		defer q.Close()

		key := "key1.key2.key3"
		producerErr = q.PublishTopic(ctx, key, []byte(key+" say hello"))

	})

	return producerErr
}

func producerFanout(queueName string) error {
	var producerErr error
	utils.SafeRunWithTimeout(time.Second*3, func(cancel context.CancelFunc) {

		c, err := rabbitmq.NewConnection(url)
		if err != nil {
			producerErr = err
			return
		}
		defer c.Close()

		ctx := context.Background()
		exchangeName := "fanout-exchange-demo"

		exchange := producer.NewFanOutExchange(exchangeName)
		q, err := producer.NewQueue(queueName, c.Conn, exchange)
		if err != nil {
			producerErr = err
			return
		}
		defer q.Close()

		producerErr = q.Publish(ctx, []byte(" say hello "))

	})
	return producerErr
}

func producerHeaders(queueName string) error {
	var producerErr error
	utils.SafeRunWithTimeout(time.Second*3, func(cancel context.CancelFunc) {

		c, err := rabbitmq.NewConnection(url)
		if err != nil {
			producerErr = err
			return
		}
		defer c.Close()

		ctx := context.Background()
		exchangeName := "headers-exchange-demo"

		kv1 := map[string]interface{}{"hello1": "world1", "foo1": "bar1"}
		exchange := producer.NewHeaderExchange(exchangeName, producer.HeadersTypeAll, kv1) // all
		q, err := producer.NewQueue(queueName, c.Conn, exchange)
		if err != nil {
			producerErr = err
			return
		}
		defer q.Close()

		headersKey1 := kv1
		err = q.PublishHeaders(ctx, headersKey1, []byte("say hello 1"))
		if err != nil {
			producerErr = err
			return
		}
		headersKey1 = map[string]interface{}{"foo": "bar"}
		producerErr = q.PublishHeaders(ctx, headersKey1, []byte("say hello 2"))

	})
	return producerErr
}
