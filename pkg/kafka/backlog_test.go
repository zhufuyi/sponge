package kafka

import (
	"testing"

	"github.com/IBM/sarama"
)

func TestInitClientManager(t *testing.T) {
	m, err := InitClientManager(addrs, groupID)
	if err != nil {
		t.Log(err)
		return
	}
	defer m.Close()
}

func testConfig() *sarama.Config {
	config := sarama.NewConfig()
	config.Consumer.Retry.Backoff = 0
	config.Producer.Retry.Backoff = 0
	config.Version = sarama.MinVersion
	config.Metadata.Retry.Max = 0
	return config
}

func TestClientManager_GetBacklog(t *testing.T) {
	seedBroker := sarama.NewMockBroker(t, 1)
	leader := sarama.NewMockBroker(t, 2)

	metadata := new(sarama.MetadataResponse)
	metadata.AddTopicPartition("foo", 0, leader.BrokerID(), nil, nil, nil, sarama.ErrNoError)
	metadata.AddTopicPartition("foo", 1, leader.BrokerID(), nil, nil, nil, sarama.ErrNoError)
	metadata.AddBroker(leader.Addr(), leader.BrokerID())
	seedBroker.Returns(metadata)

	client, err := sarama.NewClient([]string{seedBroker.Addr()}, testConfig())
	if err != nil {
		t.Fatal(err)
	}

	offsetResponse := new(sarama.OffsetResponse)
	offsetResponse.AddTopicPartition("foo", 0, 123)
	leader.Returns(offsetResponse)

	leader.Returns(&sarama.ConsumerMetadataResponse{
		Coordinator: sarama.NewBroker(leader.Addr()),
	})

	offsetManager, err := sarama.NewOffsetManagerFromClient("group", client)
	if err != nil {
		t.Error(err)
		return
	}

	fetchResponse := new(sarama.OffsetFetchResponse)
	fetchResponse.AddBlock("foo", 0, &sarama.OffsetFetchResponseBlock{
		Err:      sarama.ErrNoError,
		Offset:   123,
		Metadata: "original_meta",
	})
	leader.Returns(fetchResponse)

	m := ClientManager{
		client:        client,
		offsetManager: offsetManager,
	}
	defer m.Close()

	total, backlogs, err := m.GetBacklog("foo")
	if err != nil {
		t.Log(err)
		return
	}
	t.Log(total, backlogs)
}
