## rabbitmq

rabbitmq library wrapped in [github.com/rabbitmq/amqp091-go](github.com/rabbitmq/amqp091-go), supports automatic reconnection and customized setting of queue parameters.

### Example of use

#### Consumer code

This is a consumer code example common to the four types direct, topic, fanout, and headers.

```go
package main

import (
	"context"
	"strings"

	"github.com/zhufuyi/sponge/pkg/logger"
	"github.com/zhufuyi/sponge/pkg/rabbitmq"
	"github.com/zhufuyi/sponge/pkg/rabbitmq/consumer"
)

var handler = func(ctx context.Context, data []byte, tag ...string) error {
	tagID := strings.Join(tag, ",")
	logger.Infof("tagID=%s, receive message: %s", tagID, string(data))
	return nil
}

func main() {
	url := rabbitmq.DefaultURL
	c, err := rabbitmq.NewConnection(url, rabbitmq.WithLogger(logger.Get())) // here you can set the connection parameters, such as tls, reconnect time interval
	if err != nil {
		logger.Error("NewConnection err",logger.Err(err))
		return
	}
	defer c.Close()

	queue, err := consumer.NewQueue(context.Background(), "yourQueueName", c, consumer.WithConsumeAutoAck(false)) // here you can set the consume parameter
	if err != nil {
		logger.Error("NewQueue err",logger.Err(err))
		return
	}

	queue.Consume(handler)

	exit := make(chan struct{})
	<-exit
}    
```

<br>

#### Direct Type Code

```go
package main

import (
	"context"

	"github.com/zhufuyi/sponge/pkg/logger"
	"github.com/zhufuyi/sponge/pkg/rabbitmq"
	"github.com/zhufuyi/sponge/pkg/rabbitmq/producer"
)

func main() {
	url := rabbitmq.DefaultURL
	c, err := rabbitmq.NewConnection(url, rabbitmq.WithLogger(logger.Get())) // here you can set the connection parameters, such as tls, reconnect time interval
	if err != nil {
		logger.Error("NewConnection err",logger.Err(err))
		return
	}
	defer c.Close()

	exchangeName := "direct-exchange-demo"
	queueName := "direct-queue-1"
	routeKey := "direct-key-1"
	exchange := producer.NewDirectExchange(exchangeName, routeKey)
	q, err := producer.NewQueue(queueName, c.Conn, exchange) // here you can set the producer parameter
	if err != nil {
		logger.Error("NewQueue err",logger.Err(err))
		return
	}
	defer q.Close()

	err = q.Publish(context.Background(), []byte(routeKey+" say hello"))
	if err != nil {
		logger.Error("Publish err",logger.Err(err))
		return
	}
}    
```

<br>

#### Topic Type Code

```go
package main

import (
	"context"

	"github.com/zhufuyi/sponge/pkg/logger"
	"github.com/zhufuyi/sponge/pkg/rabbitmq"
	"github.com/zhufuyi/sponge/pkg/rabbitmq/producer"
)

func main() {
	url := rabbitmq.DefaultURL
	c, err := rabbitmq.NewConnection(url, rabbitmq.WithLogger(logger.Get())) // here you can set the connection parameters, such as tls, reconnect time interval
	if err != nil {
		logger.Error("NewConnection err",logger.Err(err))
		return
	}
	defer c.Close()

	exchangeName := "topic-exchange-demo"
	queueName := "topic-queue-1"
	routingKey := "key1.key2.*"
	exchange := producer.NewTopicExchange(exchangeName, routingKey)
	q, err := producer.NewQueue(queueName, c.Conn, exchange) // here you can set the producer parameter
	if err != nil {
		logger.Error("NewQueue err",logger.Err(err))
		return
	}
	defer q.Close()

	key:="key1.key2.key3"
	err = q.PublishTopic(context.Background(), key, []byte(key+" say hello "))
	if err != nil {
		logger.Error("PublishTopic err",logger.Err(err))
		return
	}
}    
```

<br>

#### Fanout Type Code

```go
package main

import (
	"context"

	"github.com/zhufuyi/sponge/pkg/logger"
	"github.com/zhufuyi/sponge/pkg/rabbitmq"
	"github.com/zhufuyi/sponge/pkg/rabbitmq/producer"
)

func main() {
	url := rabbitmq.DefaultURL
	c, err := rabbitmq.NewConnection(url, rabbitmq.WithLogger(logger.Get())) // here you can set the connection parameters, such as tls, reconnect time interval
	if err != nil {
		logger.Error("NewConnection err",logger.Err(err))
		return
	}
	defer c.Close()
	
	exchangeName := "fanout-exchange-demo"
	queueName := "fanout-queue-1"
	exchange := producer.NewFanOutExchange(exchangeName)
	q, err := producer.NewQueue(queueName, c.Conn, exchange) // here you can set the producer parameter
	if err != nil {
		logger.Error("NewQueue err",logger.Err(err))
		return
	}
	defer q.Close()

	err = q.Publish(context.Background(), []byte("say hello"))
	if err != nil {
		logger.Error("Publish err",logger.Err(err))
		return
	}
}    
```

<br>

#### Headers Type Code

```go
package main

import (
	"context"

	"github.com/zhufuyi/sponge/pkg/logger"
	"github.com/zhufuyi/sponge/pkg/rabbitmq"
	"github.com/zhufuyi/sponge/pkg/rabbitmq/producer"
)

func main() {
	url := rabbitmq.DefaultURL
	c, err := rabbitmq.NewConnection(url, rabbitmq.WithLogger(logger.Get())) // here you can set the connection parameters, such as tls, reconnect time interval
	if err != nil {
		logger.Error("NewConnection err",logger.Err(err))
		return
	}
	defer c.Close()


	exchangeName := "headers-exchange-demo"
	// the message is only received if there is an exact match for headers
	queueName := "headers-queue-1"
	kv1 := map[string]interface{}{"hello1": "world1", "foo1": "bar1"}
	exchange := producer.NewHeaderExchange(exchangeName, producer.HeadersTypeAll, kv1)
	q, err := producer.NewQueue(queueName, c.Conn, exchange) // here you can set the producer parameter
	if err != nil {
		logger.Error("NewQueue err",logger.Err(err))
		return
	}
	defer q.Close()
	headersKey1 := kv1 // exact match, consumer queue can receive messages
	err = q.PublishHeaders(context.Background(), headersKey1, []byte("say hello"))
	if err != nil {
		logger.Error("PublishHeaders err",logger.Err(err))
		return
	}
}    
```

<br>

#### Publish Error Handling

If the error is caused by the network, you can check if the reconnection is successful and resend it again.

```go
    err := q.Publish(context.Background(), []byte(routeKey+" say hello"))
    if err != nil {
        if errors.Is(err, producer.ErrClosed) && c.CheckConnected() { // check connection
            q, err = producer.NewQueue(queueName, c.Conn, exchange)
            if err != nil {
                logger.Info("queue reconnect failed", logger.Err(err))
            }else{
                logger.Info("queue reconnect success")
            }
        }
    }
```

