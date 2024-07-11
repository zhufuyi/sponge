package kafka

import (
	"fmt"

	"github.com/IBM/sarama"
)

// ClientManager client manager
type ClientManager struct {
	client        sarama.Client
	offsetManager sarama.OffsetManager
}

// Backlog info
type Backlog struct {
	Partition         int32 `json:"partition"`  // partition id
	Backlog           int64 `json:"backlog"`    // data backlog
	NextConsumeOffset int64 `json:"nextOffset"` // offset for next consumption
}

// InitClientManager init client manager
func InitClientManager(addrs []string, groupID string) (*ClientManager, error) {
	config := sarama.NewConfig()
	client, err := sarama.NewClient(addrs, config)
	if err != nil {
		return nil, err
	}

	offsetManager, err := sarama.NewOffsetManagerFromClient(groupID, client)
	if err != nil {
		return nil, err
	}

	return &ClientManager{
		client:        client,
		offsetManager: offsetManager,
	}, nil
}

// GetBacklog get topic backlog
func (m *ClientManager) GetBacklog(topic string) (int64, []*Backlog, error) {
	if m == nil || m.client == nil {
		return 0, nil, fmt.Errorf("client manager is nil")
	}

	var (
		total             int64
		partitionBacklogs []*Backlog
	)

	partitions, err := m.client.Partitions(topic)
	if err != nil {
		return 0, nil, err
	}

	for _, partition := range partitions {
		// get offset from kafka
		offset, err := m.client.GetOffset(topic, partition, -1)
		if err != nil {
			return 0, nil, err
		}

		// create topic/partition manager
		pom, err := m.offsetManager.ManagePartition(topic, partition)
		if err != nil {
			return 0, nil, err
		}

		var backlog int64
		// call sarama The NextOffset method of PartitionOffsetManager. Return the offset for the next consumption
		// if the consumer group has not consumed the data for this section, the return value will be -1
		n, str := pom.NextOffset()
		if str != "" {
			return 0, nil, fmt.Errorf("partition %d, %s", partition, str)
		}
		if n == -1 {
			backlog = offset
		} else {
			backlog = offset - n
		}
		total += backlog

		partitionBacklogs = append(partitionBacklogs, &Backlog{
			Partition:         partition,
			Backlog:           backlog,
			NextConsumeOffset: n,
		})
	}

	return total, partitionBacklogs, nil
}

// Close topic backlog
func (m *ClientManager) Close() error {
	if m != nil && m.client != nil {
		return m.client.Close()
	}
	return nil
}
