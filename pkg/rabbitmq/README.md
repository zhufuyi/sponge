## rabbitmq

rabbitmq library wrapped in [github.com/rabbitmq/amqp091-go](github.com/rabbitmq/amqp091-go), supports automatic reconnection and customized setting parameters, includes `direct`, `topic`, `fanout`, `headers`, `delayed message`, `publisher subscriber` a total of six message types.

### Example of use

#### Code Example

The code example includes `direct`, `topic`, `fanout`, `headers`, `delayed message`, `publisher subscriber` a total of six message types.

```go
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/zhufuyi/sponge/pkg/logger"
	"github.com/zhufuyi/sponge/pkg/rabbitmq"
)

func main() {
	url := "amqp://guest:guest@127.0.0.1:5672/"

	directExample(url)

	topicExample(url)

	fanoutExample(url)

	headersExample(url)

	delayedMessageExample(url)

	publisherSubscriberExample(url)
}

func directExample(url string) {
	exchangeName := "direct-exchange-demo"
	queueName := "direct-queue-1"
	routeKey := "direct-key-1"
	exchange := rabbitmq.NewDirectExchange(exchangeName, routeKey)
	fmt.Printf("\n\n-------------------- direct --------------------\n")

	// producer-side direct message
	func() {
		connection, err := rabbitmq.NewConnection(url, rabbitmq.WithLogger(logger.Get()))
		checkErr(err)
		defer connection.Close()

		p, err := rabbitmq.NewProducer(exchange, queueName, connection)
		checkErr(err)
		defer p.Close()

		err = p.PublishDirect(context.Background(), []byte("[direct] say hello"))
		checkErr(err)
	}()

	// consumer-side direct message
	func() {
		runConsume(url, exchange, queueName)
	}()

	<-time.After(time.Second)
}

func topicExample(url string) {
	exchangeName := "topic-exchange-demo"
	queueName := "topic-queue-1"
	routingKey := "key1.key2.*"
	exchange := rabbitmq.NewTopicExchange(exchangeName, routingKey)
	fmt.Printf("\n\n-------------------- topic --------------------\n")

	// producer-side topic message
	func() {
		connection, err := rabbitmq.NewConnection(url, rabbitmq.WithLogger(logger.Get()))
		checkErr(err)
		defer connection.Close()

		p, err := rabbitmq.NewProducer(exchange, queueName, connection)
		checkErr(err)
		defer p.Close()

		key := "key1.key2.key3"
		err = p.PublishTopic(context.Background(), key, []byte("[topic] "+key+" say hello"))
		checkErr(err)
	}()

	// consumer-side topic message
	func() {
		runConsume(url, exchange, queueName)
	}()

	<-time.After(time.Second)
}

func fanoutExample(url string) {
	exchangeName := "fanout-exchange-demo"
	queueName := "fanout-queue-1"
	exchange := rabbitmq.NewFanoutExchange(exchangeName)
	fmt.Printf("\n\n-------------------- fanout --------------------\n")

	// producer-side fanout message
	func() {
		connection, err := rabbitmq.NewConnection(url, rabbitmq.WithLogger(logger.Get()))
		checkErr(err)
		defer connection.Close()

		p, err := rabbitmq.NewProducer(exchange, queueName, connection)
		checkErr(err)
		defer p.Close()

		err = p.PublishFanout(context.Background(), []byte("[fanout] say hello"))
		checkErr(err)
	}()

	// consumer-side fanout message
	func() {
		runConsume(url, exchange, queueName)
		queueName = "fanout-queue-2"
		runConsume(url, exchange, queueName)
	}()

	<-time.After(time.Second)
}

func headersExample(url string) {
	exchangeName := "headers-exchange-demo"
	queueName := "headers-queue-1"
	headersKeys := map[string]interface{}{"hello": "world", "foo": "bar"}
	exchange := rabbitmq.NewHeadersExchange(exchangeName, rabbitmq.HeadersTypeAll, headersKeys) // all, you can set HeadersTypeAny type
	fmt.Printf("\n\n-------------------- headers --------------------\n")

	// producer-side headers message
	func() {
		connection, err := rabbitmq.NewConnection(url, rabbitmq.WithLogger(logger.Get()))
		checkErr(err)
		defer connection.Close()

		p, err := rabbitmq.NewProducer(exchange, queueName, connection)
		checkErr(err)
		defer p.Close()

		ctx := context.Background()
		headersKeys1 := headersKeys
		err = p.PublishHeaders(ctx, headersKeys1, []byte("[headers] say hello 1"))
		checkErr(err)
		headersKeys2 := map[string]interface{}{"foo": "bar"}
		err = p.PublishHeaders(ctx, headersKeys2, []byte("[headers] say hello 2"))
		checkErr(err)
	}()

	// consumer-side headers message
	func() {
		runConsume(url, exchange, queueName)
	}()

	<-time.After(time.Second)
}

func delayedMessageExample(url string) {
	exchangeName := "delayed-message-exchange-demo"
	queueName := "delayed-message-queue"
	routingKey := "delayed-key"
	exchange := rabbitmq.NewDelayedMessageExchange(exchangeName, rabbitmq.NewDirectExchange("", routingKey))
	fmt.Printf("\n\n-------------------- delayed message --------------------\n")

	// producer-side delayed message
	func() {
		connection, err := rabbitmq.NewConnection(url, rabbitmq.WithLogger(logger.Get()))
		checkErr(err)
		defer connection.Close()

		p, err := rabbitmq.NewProducer(exchange, queueName, connection)
		checkErr(err)
		defer p.Close()

		ctx := context.Background()
		datetimeLayout := "2006-01-02 15:04:05.000"
		err = p.PublishDelayedMessage(ctx, time.Second*3, []byte("[delayed message] say hello "+time.Now().Format(datetimeLayout)))
		checkErr(err)
	}()

	// consumer-side delayed message
	func() {
		runConsume(url, exchange, queueName)
	}()

	<-time.After(time.Second * 4)
}

func publisherSubscriberExample(url string) {
	channelName := "pub-sub"
	fmt.Printf("\n\n-------------------- publisher subscriber --------------------\n")

	// publisher-side message
	func() {
		connection, err := rabbitmq.NewConnection(url, rabbitmq.WithLogger(logger.Get()))
		checkErr(err)
		defer connection.Close()

		p, err := rabbitmq.NewPublisher(channelName, connection)
		checkErr(err)
		defer p.Close()

		err = p.Publish(context.Background(), []byte("[pub-sub] say hello"))
		checkErr(err)
	}()

	// subscriber-side message
	func() {
		identifier := "pub-sub-queue-1"
		runSubscriber(url, channelName, identifier)
		identifier = "pub-sub-queue-2"
		runSubscriber(url, channelName, identifier)
	}()

	<-time.After(time.Second)
}

func runConsume(url string, exchange *rabbitmq.Exchange, queueName string) {
	connection, err := rabbitmq.NewConnection(url, rabbitmq.WithLogger(logger.Get()))
	checkErr(err)

	c, err := rabbitmq.NewConsumer(exchange, queueName, connection, rabbitmq.WithConsumerAutoAck(false))
	checkErr(err)

	c.Consume(context.Background(), handler)
}

func runSubscriber(url string, channelName string, identifier string) {
	connection, err := rabbitmq.NewConnection(url, rabbitmq.WithLogger(logger.Get()))
	checkErr(err)

	s, err := rabbitmq.NewSubscriber(channelName, identifier, connection, rabbitmq.WithConsumerAutoAck(false))
	checkErr(err)

	s.Subscribe(context.Background(), handler)
}

var handler = func(ctx context.Context, data []byte, tagID string) error {
	logger.Info("received message", logger.String("tagID", tagID), logger.String("data", string(data)))
	return nil
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
```

<br>

#### Example of Automatic Resumption of Publish

If the error of publish is caused by the network, you can check if the reconnection is successful and publish it again.

```go
package main

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/zhufuyi/sponge/pkg/logger"
	"github.com/zhufuyi/sponge/pkg/rabbitmq"
)

var url = "amqp://guest:guest@127.0.0.1:5672/"

func main() {
	ctx, _ := context.WithTimeout(context.Background(), time.Hour)
	exchangeName := "direct-exchange-demo"
	queueName := "direct-queue"
	routeKey := "info"
	exchange := rabbitmq.NewDirectExchange(exchangeName, routeKey)

	err := runConsume(ctx, exchange, queueName)
	if err != nil {
		logger.Error("runConsume failed", logger.Err(err))
		return
	}

	err = runProduce(ctx, exchange, queueName)
	if err != nil {
		logger.Error("runProduce failed", logger.Err(err))
		return
	}
}

func runProduce(ctx context.Context, exchange *rabbitmq.Exchange, queueName string) error {
	connection, err := rabbitmq.NewConnection(url, rabbitmq.WithLogger(logger.Get()))
	if err != nil {
		return err
	}
	defer connection.Close()

	p, err := rabbitmq.NewProducer(exchange, queueName, connection)
	if err != nil {
		return err
	}
	defer p.Close()

	count := 0
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			count++
			data := []byte("direct say hello" + strconv.Itoa(count))
			err = p.PublishDirect(ctx, data)
			if err != nil {
				if errors.Is(err, rabbitmq.ErrClosed) {
					for {
						if !connection.CheckConnected() { // check connection
							time.Sleep(time.Second * 2)
							continue
						}
						p, err = rabbitmq.NewProducer(exchange, queueName, connection)
						if err != nil {
							logger.Warn("reconnect failed", logger.Err(err))
							time.Sleep(time.Second * 2)
							continue
						}
						break
					}
				} else {
					logger.Warn("publish failed", logger.Err(err))
				}
			}
			logger.Info("publish message", logger.String("data", string(data)))
			time.Sleep(time.Second * 5)
		}
	}
}

func runConsume(ctx context.Context, exchange *rabbitmq.Exchange, queueName string) error {
	connection, err := rabbitmq.NewConnection(url, rabbitmq.WithLogger(logger.Get()))
	if err != nil {
		return err
	}

	c, err := rabbitmq.NewConsumer(exchange, queueName, connection, rabbitmq.WithConsumerAutoAck(false))
	if err != nil {
		return err
	}

	c.Consume(ctx, handler)

	return nil
}

var handler = func(ctx context.Context, data []byte, tagID string) error {
	logger.Info("received message", logger.String("tagID", tagID), logger.String("data", string(data)))
	return nil
}
```
