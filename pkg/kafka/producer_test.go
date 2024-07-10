package kafka

import (
	"testing"
	"time"

	"github.com/IBM/sarama"
	"github.com/IBM/sarama/mocks"
	"go.uber.org/zap"

	"github.com/zhufuyi/sponge/pkg/grpc/gtls/certfile"
)

var (
	addrs = []string{"localhost:9092"}
	//addrs     = []string{"192.168.3.37:33001", "192.168.3.37:33002", "192.168.3.37:33003"}
	testTopic = "test_topic_1"

	testData = []interface{}{
		// (1) sarama.ProducerMessage type
		&sarama.ProducerMessage{
			Topic: testTopic,
			Value: sarama.StringEncoder("hello world " + time.Now().String()),
		},

		// (2) string type
		"hello world " + time.Now().String(),

		// (3) []byte type
		[]byte("hello world " + time.Now().String()),

		// (4) struct type, supports json.Marshal
		&struct {
			Name string `json:"name"`
			Age  int    `json:"age"`
		}{
			Name: "Alice",
			Age:  20,
		},

		// (5) Message type
		&Message{
			Topic: testTopic,
			Data:  []byte("hello world " + time.Now().String()),
			Key:   []byte("foobar"),
		},
	}
)

func TestInitSyncProducer(t *testing.T) {
	// Test InitSyncProducer default options
	p, err := InitSyncProducer(addrs)
	if err != nil {
		t.Log(err)
	}

	// Test InitSyncProducer with options
	p, err = InitSyncProducer(addrs,
		SyncProducerWithVersion(sarama.V3_6_0_0),
		SyncProducerWithClientID("my-client-id"),
		SyncProducerWithRequiredAcks(sarama.WaitForLocal),
		SyncProducerWithPartitioner(sarama.NewRandomPartitioner),
		SyncProducerWithReturnSuccesses(true),
		SyncProducerWithTLS(certfile.Path("two-way/server/server.pem"), certfile.Path("two-way/server/server.key"), certfile.Path("two-way/ca.pem"), true),
	)
	if err != nil {
		t.Log(err)
	}

	// Test InitSyncProducer custom options
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	p, err = InitSyncProducer(addrs, SyncProducerWithConfig(config))
	if err != nil {
		t.Log(err)
		return
	}

	time.Sleep(time.Second)
	_ = p.Close()
}

func TestSyncProducer_SendMessage(t *testing.T) {
	p, err := InitSyncProducer(addrs)
	if err != nil {
		t.Log(err)
		return
	}
	defer p.Close()

	partition, offset, err := p.SendMessage(&sarama.ProducerMessage{
		Topic: testTopic,
		Value: sarama.StringEncoder("hello world " + time.Now().String()),
	})
	if err != nil {
		t.Error(err)
		return
	}

	t.Log("partition:", partition, "offset:", offset)
}

func TestSyncProducer_SendData(t *testing.T) {
	p, err := InitSyncProducer(addrs)
	if err != nil {
		t.Log(err)
		return
	}
	defer p.Close()

	for _, data := range testData {
		partition, offset, err := p.SendData(testTopic, data)
		if err != nil {
			t.Log(err)
			continue
		}
		t.Log("partition:", partition, "offset:", offset)
	}
}

func TestInitAsyncProducer(t *testing.T) {
	// Test InitAsyncProducer default options
	p, err := InitAsyncProducer(addrs)
	if err != nil {
		t.Log(err)
	}

	// Test InitAsyncProducer with options
	p, err = InitAsyncProducer(addrs,
		AsyncProducerWithVersion(sarama.V3_6_0_0),
		AsyncProducerWithClientID("my-client-id"),
		AsyncProducerWithRequiredAcks(sarama.WaitForLocal),
		AsyncProducerWithPartitioner(sarama.NewRandomPartitioner),
		AsyncProducerWithReturnSuccesses(true),
		AsyncProducerWithFlushMessages(100),
		AsyncProducerWithFlushFrequency(time.Second),
		AsyncProducerWithFlushBytes(16*1024),
		AsyncProducerWithTLS(certfile.Path("two-way/server/server.pem"), certfile.Path("two-way/server/server.key"), certfile.Path("two-way/ca.pem"), true),
		AsyncProducerWithZapLogger(zap.NewExample()),
		AsyncProducerWithHandleFailed(func(msg *sarama.ProducerMessage) error {
			t.Logf("handle failed message: %v", msg)
			return nil
		}),
	)
	if err != nil {
		t.Log(err)
	}

	// Test InitAsyncProducer custom options
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	p, err = InitAsyncProducer(addrs, AsyncProducerWithConfig(config))
	if err != nil {
		t.Log(err)
		return
	}

	time.Sleep(time.Second)
	_ = p.Close()
}

func TestAsyncProducer_SendMessage(t *testing.T) {
	p, err := InitAsyncProducer(addrs, AsyncProducerWithFlushFrequency(time.Millisecond*100))
	if err != nil {
		t.Log(err)
		return
	}
	defer p.Close()

	msg1 := &sarama.ProducerMessage{
		Topic: testTopic,
		Value: sarama.StringEncoder("hello world " + time.Now().String()),
	}
	msg2 := &sarama.ProducerMessage{
		Topic: testTopic,
		Value: sarama.StringEncoder("foo bar " + time.Now().String()),
	}

	err = p.SendMessage(msg1, msg2)
	if err != nil {
		t.Error(err)
		return
	}
	time.Sleep(time.Millisecond * 200) // wait for messages to be sent, and flush them
}

func TestAsyncProducer_SendData(t *testing.T) {
	p, err := InitAsyncProducer(addrs, AsyncProducerWithFlushFrequency(time.Millisecond*100))
	if err != nil {
		t.Log(err)
		return
	}
	defer p.Close()

	err = p.SendData(testTopic, testData...)
	if err != nil {
		t.Error(err)
		return
	}
	time.Sleep(time.Millisecond * 200) // wait for messages to be sent, and flush them
}

func TestSyncProducer(t *testing.T) {
	sp := mocks.NewSyncProducer(t, nil)
	sp.ExpectSendMessageAndSucceed()
	p := &SyncProducer{Producer: sp}
	defer p.Close()

	msg := testData[0].(*sarama.ProducerMessage)
	partition, offset, err := p.SendMessage(msg)
	if err != nil {
		t.Log(err)
	} else {
		t.Log("partition:", partition, "offset:", offset)
	}

	for _, data := range testData {
		sp.ExpectSendMessageAndSucceed()
		p = &SyncProducer{Producer: sp}
		partition, offset, err := p.SendData(testTopic, data)
		if err != nil {
			t.Log(err)
			continue
		} else {
			t.Log("partition:", partition, "offset:", offset)
		}
	}
}

func TestAsyncProducer(t *testing.T) {
	ap := mocks.NewAsyncProducer(t, nil)
	ap.ExpectInputAndSucceed()
	p := &AsyncProducer{Producer: ap, exit: make(chan struct{}), zapLogger: zap.NewExample()}
	defer p.Close()
	go p.handleResponse(nil)

	msg := testData[0].(*sarama.ProducerMessage)
	err := p.SendMessage(msg)
	if err != nil {
		t.Log(err)
	} else {
		t.Log("send message success")
	}

	for _, data := range testData {
		ap.ExpectInputAndSucceed()
		p.Producer = ap
		err := p.SendData(testTopic, data)
		if err != nil {
			t.Log(err)
			continue
		} else {
			t.Log("send message success")
		}
	}
}
