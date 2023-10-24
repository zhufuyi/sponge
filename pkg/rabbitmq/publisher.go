package rabbitmq

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

// Publisher session
type Publisher struct {
	*Producer
}

// NewPublisher create a publisher, channelName is exchange name
func NewPublisher(channelName string, connection *Connection, opts ...ProducerOption) (*Publisher, error) {
	o := defaultProducerOptions()
	o.apply(opts...)

	exchange := NewFanoutExchange(channelName)

	// crate a new channel
	ch, err := connection.conn.Channel()
	if err != nil {
		return nil, err
	}

	// declare the exchange type
	err = ch.ExchangeDeclare(
		channelName,
		exchangeTypeFanout,
		o.isPersistent,
		o.exchangeDeclare.autoDelete,
		o.exchangeDeclare.internal,
		o.exchangeDeclare.noWait,
		o.exchangeDeclare.args,
	)
	if err != nil {
		_ = ch.Close()
		return nil, err
	}

	deliveryMode := amqp.Persistent
	if !o.isPersistent {
		deliveryMode = amqp.Transient
	}

	connection.zapLog.Info("[rabbit producer] initialized", zap.String("channel", channelName), zap.Bool("isPersistent", o.isPersistent))

	p := &Producer{
		Exchange:     exchange,
		conn:         connection.conn,
		ch:           ch,
		isPersistent: o.isPersistent,
		deliveryMode: deliveryMode,
		mandatory:    o.mandatory,
		zapLog:       connection.zapLog,
	}

	return &Publisher{p}, nil
}

func (p *Publisher) Publish(ctx context.Context, body []byte) error {
	return p.ch.PublishWithContext(
		ctx,
		p.Exchange.name,
		p.Exchange.routingKey,
		p.mandatory,
		false,
		amqp.Publishing{
			DeliveryMode: p.deliveryMode,
			ContentType:  "text/plain",
			Body:         body,
		},
	)
}

// Close publisher
func (p *Publisher) Close() {
	if p.ch != nil {
		_ = p.ch.Close()
	}
}
