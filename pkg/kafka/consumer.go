package kafka

import (
	"context"

	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

// ---------------------------------- consume group---------------------------------------

// ConsumerGroup consume group
type ConsumerGroup struct {
	Group            sarama.ConsumerGroup
	groupID          string
	zapLogger        *zap.Logger
	autoCommitEnable bool
}

// InitConsumerGroup init consumer group
func InitConsumerGroup(addrs []string, groupID string, opts ...ConsumerOption) (*ConsumerGroup, error) {
	o := defaultConsumerOptions()
	o.apply(opts...)

	var config *sarama.Config
	if o.config != nil {
		config = o.config
	} else {
		config = sarama.NewConfig()
		config.Version = o.version
		config.Consumer.Offsets.Initial = o.offsetsInitial
		config.Consumer.Offsets.AutoCommit.Enable = o.offsetsAutoCommitEnable
		config.Consumer.Offsets.AutoCommit.Interval = o.offsetsAutoCommitInterval
		config.ClientID = o.clientID
		if o.tlsConfig != nil {
			config.Net.TLS.Config = o.tlsConfig
			config.Net.TLS.Enable = true
		}
	}

	consumer, err := sarama.NewConsumerGroup(addrs, groupID, config)
	if err != nil {
		return nil, err
	}
	return &ConsumerGroup{
		Group:            consumer,
		groupID:          groupID,
		zapLogger:        o.zapLogger,
		autoCommitEnable: config.Consumer.Offsets.AutoCommit.Enable,
	}, nil
}

// Consume consume messages
func (c *ConsumerGroup) Consume(ctx context.Context, topics []string, handleMessageFn HandleMessageFn) error {
	handler := &defaultConsumerHandler{
		ctx:              ctx,
		handleMessageFn:  handleMessageFn,
		zapLogger:        c.zapLogger,
		autoCommitEnable: c.autoCommitEnable,
	}

	err := c.Group.Consume(ctx, topics, handler)
	if err != nil {
		c.zapLogger.Error("failed to consume messages", zap.String("group_id", c.groupID), zap.Strings("topics", topics), zap.Error(err))
		return err
	}
	return nil
}

// ConsumeCustom consume messages for custom handler, you need to implement the sarama.ConsumerGroupHandler interface
func (c *ConsumerGroup) ConsumeCustom(ctx context.Context, topics []string, handler sarama.ConsumerGroupHandler) error {
	err := c.Group.Consume(ctx, topics, handler)
	if err != nil {
		c.zapLogger.Error("failed to consume messages", zap.String("group_id", c.groupID), zap.Strings("topics", topics), zap.Error(err))
		return err
	}
	return nil
}

func (c *ConsumerGroup) Close() error {
	if c == nil || c.Group == nil {
		return c.Group.Close()
	}
	return nil
}

type defaultConsumerHandler struct {
	ctx              context.Context
	handleMessageFn  HandleMessageFn
	zapLogger        *zap.Logger
	autoCommitEnable bool
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (h *defaultConsumerHandler) Setup(sess sarama.ConsumerGroupSession) error {
	h.zapLogger.Info("consumer group session [setup]", zap.Any("claims", sess.Claims()))
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (h *defaultConsumerHandler) Cleanup(sess sarama.ConsumerGroupSession) error {
	h.zapLogger.Info("consumer group session [cleanup]", zap.Any("claims", sess.Claims()))
	return nil
}

// ConsumeClaim consumes messages
func (h *defaultConsumerHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	defer func() {
		if e := recover(); e != nil {
			h.zapLogger.Error("panic occurred while consuming messages", zap.Any("error", e))
			_ = h.ConsumeClaim(sess, claim)
		}
	}()

	for {
		select {
		case <-h.ctx.Done():
			return nil
		case msg, ok := <-claim.Messages():
			if !ok {
				return nil
			}
			err := h.handleMessageFn(msg)
			if err != nil {
				h.zapLogger.Error("failed to handle message", zap.Error(err))
				continue
			}
			sess.MarkMessage(msg, "")
			if !h.autoCommitEnable {
				sess.Commit()
			}
		}
	}
}

// ---------------------------------- consume partition------------------------------------

// Consumer consume partition
type Consumer struct {
	C         sarama.Consumer
	zapLogger *zap.Logger
}

// InitConsumer init consumer
func InitConsumer(addrs []string, opts ...ConsumerOption) (*Consumer, error) {
	o := defaultConsumerOptions()
	o.apply(opts...)

	var config *sarama.Config
	if o.config != nil {
		config = o.config
	} else {
		config = sarama.NewConfig()
		config.Version = o.version
		config.Consumer.Return.Errors = true
		config.ClientID = o.clientID
		if o.tlsConfig != nil {
			config.Net.TLS.Config = o.tlsConfig
			config.Net.TLS.Enable = true
		}
	}

	consumer, err := sarama.NewConsumer(addrs, config)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		C:         consumer,
		zapLogger: o.zapLogger,
	}, nil
}

// ConsumePartition consumer one partition, blocking
func (c *Consumer) ConsumePartition(ctx context.Context, topic string, partition int32, offset int64, handleFn HandleMessageFn) {
	defer func() {
		if e := recover(); e != nil {
			c.zapLogger.Error("panic occurred while consuming messages", zap.Any("error", e))
			c.ConsumePartition(ctx, topic, partition, offset, handleFn)
		}
	}()

	pc, err := c.C.ConsumePartition(topic, partition, offset)
	if err != nil {
		c.zapLogger.Error("failed to create partition consumer", zap.Error(err), zap.String("topic", topic), zap.Int32("partition", partition))
		return
	}

	c.zapLogger.Info("start consuming partition", zap.String("topic", topic), zap.Int32("partition", partition), zap.Int64("offset", offset))

	for {
		select {
		case msg := <-pc.Messages():
			err = handleFn(msg)
			if err != nil {
				c.zapLogger.Warn("failed to handle message", zap.Error(err), zap.String("topic", topic), zap.Int32("partition", partition), zap.Int64("offset", msg.Offset))
			}
		case err := <-pc.Errors():
			c.zapLogger.Error("partition consumer error", zap.Any("err", err))
		case <-ctx.Done():
			return
		}
	}
}

// ConsumeAllPartition consumer all partitions, no blocking
func (c *Consumer) ConsumeAllPartition(ctx context.Context, topic string, offset int64, handleFn HandleMessageFn) {
	partitionList, err := c.C.Partitions(topic)
	if err != nil {
		c.zapLogger.Error("failed to get partition", zap.Error(err))
		return
	}

	for _, partition := range partitionList {
		go func(partition int32, offset int64) {
			c.ConsumePartition(ctx, topic, partition, offset, handleFn)
		}(partition, offset)
	}
}

// Close the consumer
func (c *Consumer) Close() error {
	if c == nil || c.C == nil {
		return c.C.Close()
	}
	return nil
}
