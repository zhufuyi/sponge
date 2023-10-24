package rabbitmq

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/zhufuyi/sponge/pkg/utils"

	amqp "github.com/rabbitmq/amqp091-go"
)

var testChannelName = "pub-sub"

func runPublisher(ctx context.Context, channelName string) error {
	var publisherErr error
	utils.SafeRunWithTimeout(time.Second*3, func(cancel context.CancelFunc) {
		defer cancel()
		connection, err := NewConnection(url)
		if err != nil {
			publisherErr = err
			return
		}
		defer connection.Close()

		p, err := NewPublisher(channelName, connection)
		if err != nil {
			publisherErr = err
			return
		}
		defer p.Close()

		data := []byte("hello world " + time.Now().Format(datetimeLayout))
		err = p.Publish(ctx, data)
		if err != nil {
			publisherErr = err
			return
		}
		fmt.Printf("[send]: %s\n", data)
	})
	return publisherErr
}

func TestPublisher(t *testing.T) {
	err := runPublisher(context.Background(), testChannelName)
	if err != nil {
		t.Log(err)
		return
	}
}

func TestPublisherErr(t *testing.T) {
	utils.SafeRunWithTimeout(time.Second, func(cancel context.CancelFunc) {
		defer cancel()
		_, err := NewPublisher(testChannelName, &Connection{conn: &amqp.Connection{}})
		if err != nil {
			t.Log(err)
			return
		}
	})

	p := &Publisher{&Producer{conn: &amqp.Connection{}, ch: &amqp.Channel{}}}
	utils.SafeRun(context.Background(), func(ctx context.Context) {
		_ = p.Publish(context.Background(), []byte("hello world"))
	})
	utils.SafeRun(context.Background(), func(ctx context.Context) {
		p.Close()
	})
}
