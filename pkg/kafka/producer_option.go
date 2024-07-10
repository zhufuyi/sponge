package kafka

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"time"

	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

// -------------------------------------- sync producer ------------------------------------

// SyncProducerOption set options.
type SyncProducerOption func(*syncProducerOptions)

type syncProducerOptions struct {
	version         sarama.KafkaVersion           // default V2_1_0_0
	requiredAcks    sarama.RequiredAcks           // default WaitForAll
	partitioner     sarama.PartitionerConstructor // default NewHashPartitioner
	returnSuccesses bool                          // default true
	clientID        string                        // default "sarama"
	tlsConfig       *tls.Config                   // default nil

	// custom config, if not nil, it will override the default config, the above parameters are invalid
	config *sarama.Config // default nil
}

func (o *syncProducerOptions) apply(opts ...SyncProducerOption) {
	for _, opt := range opts {
		opt(o)
	}
}

func defaultSyncProducerOptions() *syncProducerOptions {
	return &syncProducerOptions{
		version:         sarama.V2_1_0_0,
		requiredAcks:    sarama.WaitForAll,
		partitioner:     sarama.NewHashPartitioner,
		returnSuccesses: true,
		clientID:        "sarama",
	}
}

// SyncProducerWithVersion set kafka version.
func SyncProducerWithVersion(version sarama.KafkaVersion) SyncProducerOption {
	return func(o *syncProducerOptions) {
		o.version = version
	}
}

// SyncProducerWithRequiredAcks set requiredAcks.
func SyncProducerWithRequiredAcks(requiredAcks sarama.RequiredAcks) SyncProducerOption {
	return func(o *syncProducerOptions) {
		o.requiredAcks = requiredAcks
	}
}

// SyncProducerWithPartitioner set partitioner.
func SyncProducerWithPartitioner(partitioner sarama.PartitionerConstructor) SyncProducerOption {
	return func(o *syncProducerOptions) {
		o.partitioner = partitioner
	}
}

// SyncProducerWithReturnSuccesses set returnSuccesses.
func SyncProducerWithReturnSuccesses(returnSuccesses bool) SyncProducerOption {
	return func(o *syncProducerOptions) {
		o.returnSuccesses = returnSuccesses
	}
}

// SyncProducerWithClientID set clientID.
func SyncProducerWithClientID(clientID string) SyncProducerOption {
	return func(o *syncProducerOptions) {
		o.clientID = clientID
	}
}

// SyncProducerWithTLS set tlsConfig, if isSkipVerify is true, crypto/tls accepts any certificate presented by
// the server and any host name in that certificate.
func SyncProducerWithTLS(certFile, keyFile, caFile string, isSkipVerify bool) SyncProducerOption {
	return func(o *syncProducerOptions) {
		var err error
		o.tlsConfig, err = getTLSConfig(certFile, keyFile, caFile, isSkipVerify)
		if err != nil {
			fmt.Println("SyncProducerWithTLS error:", err)
		}
	}
}

// SyncProducerWithConfig set custom config.
func SyncProducerWithConfig(config *sarama.Config) SyncProducerOption {
	return func(o *syncProducerOptions) {
		o.config = config
	}
}

// -------------------------------------- async producer -----------------------------------

// AsyncSendFailedHandlerFn is a function that handles failed messages.
type AsyncSendFailedHandlerFn func(msg *sarama.ProducerMessage) error

// AsyncProducerOption set options.
type AsyncProducerOption func(*asyncProducerOptions)

type asyncProducerOptions struct {
	version         sarama.KafkaVersion           // default V2_1_0_0
	requiredAcks    sarama.RequiredAcks           // default WaitForLocal
	partitioner     sarama.PartitionerConstructor // default NewHashPartitioner
	returnSuccesses bool                          // default true
	clientID        string                        // default "sarama"
	flushMessages   int                           // default 20
	flushFrequency  time.Duration                 // default 2 second
	flushBytes      int                           // default 0
	tlsConfig       *tls.Config

	// custom config, if not nil, it will override the default config, the above parameters are invalid
	config *sarama.Config // default nil

	zapLogger      *zap.Logger              // default NewProduction
	handleFailedFn AsyncSendFailedHandlerFn // default nil
}

func (o *asyncProducerOptions) apply(opts ...AsyncProducerOption) {
	for _, opt := range opts {
		opt(o)
	}
}

func defaultAsyncProducerOptions() *asyncProducerOptions {
	zapLogger, _ := zap.NewProduction()
	return &asyncProducerOptions{
		version:         sarama.V2_1_0_0,
		requiredAcks:    sarama.WaitForLocal,
		partitioner:     sarama.NewHashPartitioner,
		returnSuccesses: true,
		clientID:        "sarama",
		flushMessages:   20,
		flushFrequency:  2 * time.Second,
		zapLogger:       zapLogger,
	}
}

// AsyncProducerWithVersion set kafka version.
func AsyncProducerWithVersion(version sarama.KafkaVersion) AsyncProducerOption {
	return func(o *asyncProducerOptions) {
		o.version = version
	}
}

// AsyncProducerWithRequiredAcks set requiredAcks.
func AsyncProducerWithRequiredAcks(requiredAcks sarama.RequiredAcks) AsyncProducerOption {
	return func(o *asyncProducerOptions) {
		o.requiredAcks = requiredAcks
	}
}

// AsyncProducerWithPartitioner set partitioner.
func AsyncProducerWithPartitioner(partitioner sarama.PartitionerConstructor) AsyncProducerOption {
	return func(o *asyncProducerOptions) {
		o.partitioner = partitioner
	}
}

// AsyncProducerWithReturnSuccesses set returnSuccesses.
func AsyncProducerWithReturnSuccesses(returnSuccesses bool) AsyncProducerOption {
	return func(o *asyncProducerOptions) {
		o.returnSuccesses = returnSuccesses
	}
}

// AsyncProducerWithClientID set clientID.
func AsyncProducerWithClientID(clientID string) AsyncProducerOption {
	return func(o *asyncProducerOptions) {
		o.clientID = clientID
	}
}

// AsyncProducerWithFlushMessages set flushMessages.
func AsyncProducerWithFlushMessages(flushMessages int) AsyncProducerOption {
	return func(o *asyncProducerOptions) {
		o.flushMessages = flushMessages
	}
}

// AsyncProducerWithFlushFrequency set flushFrequency.
func AsyncProducerWithFlushFrequency(flushFrequency time.Duration) AsyncProducerOption {
	return func(o *asyncProducerOptions) {
		o.flushFrequency = flushFrequency
	}
}

// AsyncProducerWithFlushBytes set flushBytes.
func AsyncProducerWithFlushBytes(flushBytes int) AsyncProducerOption {
	return func(o *asyncProducerOptions) {
		o.flushBytes = flushBytes
	}
}

// AsyncProducerWithTLS set tlsConfig, if isSkipVerify is true, crypto/tls accepts any certificate presented by
// the server and any host name in that certificate.
func AsyncProducerWithTLS(certFile, keyFile, caFile string, isSkipVerify bool) AsyncProducerOption {
	return func(o *asyncProducerOptions) {
		var err error
		o.tlsConfig, err = getTLSConfig(certFile, keyFile, caFile, isSkipVerify)
		if err != nil {
			fmt.Println("AsyncProducerWithTLS error:", err)
		}
	}
}

// AsyncProducerWithZapLogger set zapLogger.
func AsyncProducerWithZapLogger(zapLogger *zap.Logger) AsyncProducerOption {
	return func(o *asyncProducerOptions) {
		if zapLogger != nil {
			o.zapLogger = zapLogger
		}
	}
}

// AsyncProducerWithHandleFailed set handleFailedFn.
func AsyncProducerWithHandleFailed(handleFailedFn AsyncSendFailedHandlerFn) AsyncProducerOption {
	return func(o *asyncProducerOptions) {
		o.handleFailedFn = handleFailedFn
	}
}

// AsyncProducerWithConfig set custom config.
func AsyncProducerWithConfig(config *sarama.Config) AsyncProducerOption {
	return func(o *asyncProducerOptions) {
		o.config = config
	}
}

func getTLSConfig(certFile, keyFile, caFile string, isSkipVerify bool) (*tls.Config, error) {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, err
	}

	caCert, err := os.ReadFile(caFile)
	if err != nil {
		return nil, err
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	return &tls.Config{
		Certificates:       []tls.Certificate{cert},
		RootCAs:            caCertPool,
		InsecureSkipVerify: isSkipVerify,
	}, nil
}
