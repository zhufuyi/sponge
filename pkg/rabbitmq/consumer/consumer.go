// Package consumer is the generic consumer-side processing logic for the four modes direct, topic, fanout, headers
package consumer

import (
	"context"
	"strconv"
	"time"

	"github.com/zhufuyi/sponge/pkg/rabbitmq"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

// Queue session
type Queue struct {
	name string
	c    *rabbitmq.Connection
	ch   *amqp.Channel

	ctx     context.Context
	autoAck bool

	consumeOption *consumeOptions
	zapLog        *zap.Logger
}

// Close queue
func (q *Queue) Close() {
	if q.ch != nil {
		_ = q.ch.Close()
	}
}

// Handler message
type Handler func(ctx context.Context, data []byte, tagID ...string) error

// NewQueue create a queue
func NewQueue(ctx context.Context, name string, c *rabbitmq.Connection, opts ...ConsumeOption) (*Queue, error) {
	o := defaultConsumeOptions()
	o.apply(opts...)

	q := &Queue{
		name: name,

		c:             c,
		consumeOption: o,
		zapLog:        c.ZapLog,
		autoAck:       o.autoAck,
		ctx:           ctx,
	}

	return q, nil
}

func (q *Queue) newChannel() error {
	q.c.Mutex.Lock()

	// crate a new channel
	ch, err := q.c.Conn.Channel()
	if err != nil {
		q.c.Mutex.Unlock()
		return err
	}
	q.ch = ch

	q.c.Mutex.Unlock()

	// setting the prefetch value
	// set channel.Qos on the consumer side to limit the number of messages consumed at a time,
	// balancing message throughput and fairness, and preventing consumers from being hit by bursty message traffic.
	o := q.consumeOption
	if o.enableQos {
		err = ch.Qos(o.qos.prefetchCount, o.qos.prefetchSize, o.qos.global)
		if err != nil {
			_ = ch.Close()
			return err
		}
	}

	q.zapLog.Info("[rabbitmq consumer] create a queue success", zap.String("name", q.name), zap.Bool("autoAck", o.autoAck))
	return nil
}

func (q *Queue) consumeWithContext() (<-chan amqp.Delivery, error) {
	return q.ch.ConsumeWithContext(
		q.ctx,
		q.name,
		q.consumeOption.consumer,
		q.consumeOption.autoAck,
		q.consumeOption.exclusive,
		q.consumeOption.noLocal,
		q.consumeOption.noWait,
		q.consumeOption.args,
	)
}

// Consume messages for loop in goroutine
func (q *Queue) Consume(handler Handler) {
	go func() {
		ticker := time.NewTicker(time.Second * 2)
		for {
			ticker.Reset(time.Second * 2)

			// check connection for loop
			select {
			case <-ticker.C:
				if !q.c.CheckConnected() {
					continue
				}
			case <-q.c.Exit:
				q.Close()
				return
			}
			ticker.Stop()

			err := q.newChannel()
			if err != nil {
				q.zapLog.Warn("[rabbitmq consumer] create a channel error", zap.String("err", err.Error()))
				continue
			}

			delivery, err := q.consumeWithContext()
			if err != nil {
				q.zapLog.Warn("[rabbitmq consumer] execution of consumption error", zap.String("err", err.Error()))
				continue
			}
			q.zapLog.Info("[rabbitmq consumer] queue is ready and waiting for messages, queue=" + q.name)

			isContinueConsume := false
			for {
				select {
				case <-q.c.Exit:
					q.Close()
					return
				case d, ok := <-delivery:
					if !ok {
						q.zapLog.Warn("[rabbitmq consumer] queue receive message exception exit, queue=" + q.name)
						isContinueConsume = true
						break
					}
					tagID := q.name + "/" + strconv.FormatUint(d.DeliveryTag, 10)
					err = handler(q.ctx, d.Body, tagID)
					if err != nil {
						q.zapLog.Warn("[rabbitmq consumer] handle message error", zap.String("err", err.Error()))
						continue
					}
					if !q.autoAck {
						if err = d.Ack(false); err != nil {
							q.zapLog.Warn("[rabbitmq consumer] manual ack error", zap.String("err", err.Error()), zap.String("tagID", tagID))
							continue
						}
						q.zapLog.Info("[rabbitmq consumer] manual ack done", zap.String("tagID", tagID))
					}
				}

				if isContinueConsume {
					break
				}
			}
		}
	}()
}
