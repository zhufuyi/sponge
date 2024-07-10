package kafka

import (
	"crypto/tls"
	"fmt"
	"time"

	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

// HandleMessageFn is a function that handles a message from a partition consumer
type HandleMessageFn func(msg *sarama.ConsumerMessage) error

// ConsumerOption set options.
type ConsumerOption func(*consumerOptions)

type consumerOptions struct {
	version   sarama.KafkaVersion // default V2_1_0_0
	clientID  string              // default "sarama"
	tlsConfig *tls.Config         // default nil

	// consumer group options
	offsetsInitial            int64         // default OffsetOldest
	offsetsAutoCommitEnable   bool          // default true
	offsetsAutoCommitInterval time.Duration // default 1s, when offsetsAutoCommitEnable is true

	// custom config, if not nil, it will override the default config, the above parameters are invalid
	config *sarama.Config // default nil

	zapLogger *zap.Logger // default NewProduction
}

func (o *consumerOptions) apply(opts ...ConsumerOption) {
	for _, opt := range opts {
		opt(o)
	}
}

func defaultConsumerOptions() *consumerOptions {
	zapLogger, _ := zap.NewProduction()
	return &consumerOptions{
		version:                   sarama.V2_1_0_0,
		offsetsInitial:            sarama.OffsetOldest,
		offsetsAutoCommitEnable:   true,
		offsetsAutoCommitInterval: time.Second,
		clientID:                  "sarama",
		zapLogger:                 zapLogger,
	}
}

// ConsumerWithVersion set kafka version.
func ConsumerWithVersion(version sarama.KafkaVersion) ConsumerOption {
	return func(o *consumerOptions) {
		o.version = version
	}
}

// ConsumerWithOffsetsInitial set offsetsInitial.
func ConsumerWithOffsetsInitial(offsetsInitial int64) ConsumerOption {
	return func(o *consumerOptions) {
		o.offsetsInitial = offsetsInitial
	}
}

// ConsumerWithOffsetsAutoCommitEnable set offsetsAutoCommitEnable.
func ConsumerWithOffsetsAutoCommitEnable(offsetsAutoCommitEnable bool) ConsumerOption {
	return func(o *consumerOptions) {
		o.offsetsAutoCommitEnable = offsetsAutoCommitEnable
	}
}

// ConsumerWithOffsetsAutoCommitInterval set offsetsAutoCommitInterval.
func ConsumerWithOffsetsAutoCommitInterval(offsetsAutoCommitInterval time.Duration) ConsumerOption {
	return func(o *consumerOptions) {
		o.offsetsAutoCommitInterval = offsetsAutoCommitInterval
	}
}

// ConsumerWithClientID set clientID.
func ConsumerWithClientID(clientID string) ConsumerOption {
	return func(o *consumerOptions) {
		o.clientID = clientID
	}
}

// ConsumerWithTLS set tlsConfig, if isSkipVerify is true, crypto/tls accepts any certificate presented by
// the server and any host name in that certificate.
func ConsumerWithTLS(certFile, keyFile, caFile string, isSkipVerify bool) ConsumerOption {
	return func(o *consumerOptions) {
		var err error
		o.tlsConfig, err = getTLSConfig(certFile, keyFile, caFile, isSkipVerify)
		if err != nil {
			fmt.Println("ConsumerWithTLS error:", err)
		}
	}
}

// ConsumerWithZapLogger set zapLogger.
func ConsumerWithZapLogger(zapLogger *zap.Logger) ConsumerOption {
	return func(o *consumerOptions) {
		if zapLogger != nil {
			o.zapLogger = zapLogger
		}
	}
}

// ConsumerWithConfig set custom config.
func ConsumerWithConfig(config *sarama.Config) ConsumerOption {
	return func(o *consumerOptions) {
		o.config = config
	}
}
