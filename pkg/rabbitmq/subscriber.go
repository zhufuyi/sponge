package rabbitmq

import (
	"context"
)

// Subscriber session
type Subscriber struct {
	*Consumer
}

// NewSubscriber create a subscriber, channelName is exchange name, identifier is queue name
func NewSubscriber(channelName string, identifier string, connection *Connection, opts ...ConsumerOption) (*Subscriber, error) {
	exchange := NewFanoutExchange(channelName)
	queueName := identifier
	c, err := NewConsumer(exchange, queueName, connection, opts...)
	if err != nil {
		return nil, err
	}
	return &Subscriber{c}, nil
}

// Subscribe and handle message
func (s *Subscriber) Subscribe(ctx context.Context, handler Handler) {
	s.Consume(ctx, handler)
}

// Close subscriber
func (s *Subscriber) Close() {
	if s.ch != nil {
		_ = s.ch.Close()
	}
}
