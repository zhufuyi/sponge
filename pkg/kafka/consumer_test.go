package kafka

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/IBM/sarama"
	"go.uber.org/zap"

	"github.com/zhufuyi/sponge/pkg/grpc/gtls/certfile"
)

var (
	groupID  = "my-group"
	waitTime = time.Second * 10

	handleMsgFn = func(msg *sarama.ConsumerMessage) error {
		fmt.Printf("received msg: topic=%s, partition=%d, offset=%d, key=%s, val=%s\n",
			msg.Topic, msg.Partition, msg.Offset, msg.Key, msg.Value)
		return nil
	}
)

type myConsumerGroupHandler struct {
	autoCommitEnable bool
}

func (h *myConsumerGroupHandler) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (h *myConsumerGroupHandler) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (h *myConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		fmt.Printf("received msg: topic=%s, partition=%d, offset=%d, key=%s, val=%s\n",
			msg.Topic, msg.Partition, msg.Offset, msg.Key, msg.Value)
		session.MarkMessage(msg, "")
		if !h.autoCommitEnable {
			session.Commit()
		}
	}
	return nil
}

func TestInitConsumerGroup(t *testing.T) {
	// Test InitConsumerGroup default options
	cg, err := InitConsumerGroup(addrs, groupID)
	if err != nil {
		t.Log(err)
	}

	// Test InitConsumerGroup with options
	cg, err = InitConsumerGroup(addrs, groupID,
		ConsumerWithVersion(sarama.V3_6_0_0),
		ConsumerWithClientID("my-client-id"),
		ConsumerWithGroupStrategies(sarama.NewBalanceStrategySticky()),
		ConsumerWithOffsetsInitial(sarama.OffsetOldest),
		ConsumerWithOffsetsAutoCommitEnable(true),
		ConsumerWithOffsetsAutoCommitInterval(time.Second),
		ConsumerWithTLS(certfile.Path("two-way/server/server.pem"), certfile.Path("two-way/server/server.key"), certfile.Path("two-way/ca.pem"), true),
		ConsumerWithZapLogger(zap.NewNop()),
	)
	if err != nil {
		t.Log(err)
	}

	// Test InitConsumerGroup custom options
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	cg, err = InitConsumerGroup(addrs, groupID, ConsumerWithConfig(config))
	if err != nil {
		t.Log(err)
		return
	}

	time.Sleep(time.Second)
	_ = cg.Close()
}

func TestConsumerGroup_Consume(t *testing.T) {
	cg, err := InitConsumerGroup(addrs, groupID,
		ConsumerWithVersion(sarama.V3_6_0_0),
		ConsumerWithOffsetsInitial(sarama.OffsetOldest),
		ConsumerWithOffsetsAutoCommitEnable(true),
		ConsumerWithOffsetsAutoCommitInterval(time.Second),
	)
	if err != nil {
		t.Log(err)
		return
	}
	defer cg.Close()

	go cg.Consume(context.Background(), []string{testTopic}, handleMsgFn)

	<-time.After(waitTime)
}

func TestConsumerGroup_ConsumeCustom(t *testing.T) {
	cg, err := InitConsumerGroup(addrs, groupID,
		ConsumerWithVersion(sarama.V3_6_0_0),
		ConsumerWithOffsetsAutoCommitEnable(false),
	)
	if err != nil {
		t.Log(err)
		return
	}
	defer cg.Close()

	cgh := &myConsumerGroupHandler{autoCommitEnable: cg.autoCommitEnable}
	go cg.ConsumeCustom(context.Background(), []string{testTopic}, cgh)

	<-time.After(waitTime)
}

func TestInitConsumer(t *testing.T) {
	// Test InitConsumer default options
	c, err := InitConsumer(addrs)
	if err != nil {
		t.Log(err)
	}

	// Test InitConsumer with options
	c, err = InitConsumer(addrs,
		ConsumerWithVersion(sarama.V3_6_0_0),
		ConsumerWithClientID("my-client-id"),
		ConsumerWithTLS(certfile.Path("two-way/server/server.pem"), certfile.Path("two-way/server/server.key"), certfile.Path("two-way/ca.pem"), true),
		ConsumerWithZapLogger(zap.NewNop()),
	)
	if err != nil {
		t.Log(err)
	}

	// Test InitConsumer custom options
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	c, err = InitConsumer(addrs, ConsumerWithConfig(config))
	if err != nil {
		t.Log(err)
		return
	}

	time.Sleep(time.Second)
	_ = c.Close()
}

func TestConsumer_ConsumePartition(t *testing.T) {
	c, err := InitConsumer(addrs,
		ConsumerWithVersion(sarama.V3_6_0_0),
		ConsumerWithClientID("my-client-id"),
	)
	if err != nil {
		t.Log(err)
		return
	}
	defer c.Close()

	go c.ConsumePartition(context.Background(), testTopic, 0, sarama.OffsetNewest, handleMsgFn)

	<-time.After(waitTime)
}

func TestConsumer_ConsumeAllPartition(t *testing.T) {
	c, err := InitConsumer(addrs,
		ConsumerWithVersion(sarama.V3_6_0_0),
		ConsumerWithClientID("my-client-id"),
	)
	if err != nil {
		t.Log(err)
		return
	}
	defer c.Close()

	c.ConsumeAllPartition(context.Background(), testTopic, sarama.OffsetNewest, handleMsgFn)

	<-time.After(waitTime)
}

func TestConsumerGroup(t *testing.T) {
	var (
		myTopic = "my-topic"
		myGroup = "my_group"
	)

	broker0 := sarama.NewMockBroker(t, 0)
	defer broker0.Close()

	mockData := map[string]sarama.MockResponse{
		"MetadataRequest": sarama.NewMockMetadataResponse(t).
			SetBroker(broker0.Addr(), broker0.BrokerID()).
			SetLeader(myTopic, 0, broker0.BrokerID()),
		"OffsetRequest": sarama.NewMockOffsetResponse(t).
			SetOffset(myTopic, 0, sarama.OffsetOldest, 0).
			SetOffset(myTopic, 0, sarama.OffsetNewest, 1),
		"FindCoordinatorRequest": sarama.NewMockFindCoordinatorResponse(t).
			SetCoordinator(sarama.CoordinatorGroup, myGroup, broker0),
		"HeartbeatRequest": sarama.NewMockHeartbeatResponse(t),
		"JoinGroupRequest": sarama.NewMockSequence(
			sarama.NewMockJoinGroupResponse(t).SetError(sarama.ErrOffsetsLoadInProgress),
			sarama.NewMockJoinGroupResponse(t).SetGroupProtocol(sarama.RangeBalanceStrategyName),
		),
		"SyncGroupRequest": sarama.NewMockSequence(
			sarama.NewMockSyncGroupResponse(t).SetError(sarama.ErrOffsetsLoadInProgress),
			sarama.NewMockSyncGroupResponse(t).SetMemberAssignment(
				&sarama.ConsumerGroupMemberAssignment{
					Version: 0,
					Topics: map[string][]int32{
						myTopic: {0},
					},
				}),
		),
		"OffsetFetchRequest": sarama.NewMockOffsetFetchResponse(t).SetOffset(
			myGroup, myTopic, 0, 0, "", sarama.ErrNoError,
		).SetError(sarama.ErrNoError),
		"FetchRequest": sarama.NewMockSequence(
			sarama.NewMockFetchResponse(t, 1).
				SetMessage(myTopic, 0, 0, sarama.StringEncoder("foo")).
				SetMessage(myTopic, 0, 1, sarama.StringEncoder("bar")),
			sarama.NewMockFetchResponse(t, 1),
		),
	}

	broker0.SetHandlerByMap(mockData)

	config := sarama.NewConfig()
	config.ClientID = t.Name()
	config.Version = sarama.V2_0_0_0
	config.Consumer.Return.Errors = true
	config.Consumer.Group.Rebalance.Retry.Max = 2
	config.Consumer.Group.Rebalance.Retry.Backoff = 0
	config.Consumer.Offsets.AutoCommit.Enable = false
	group, err := sarama.NewConsumerGroup([]string{broker0.Addr()}, myGroup, config)
	if err != nil {
		t.Fatal(err)
	}

	topics := []string{myTopic}
	g := &ConsumerGroup{
		Group:            group,
		groupID:          myGroup,
		zapLogger:        zap.NewExample(),
		autoCommitEnable: false,
	}
	defer g.Close()

	ctx, cancel := context.WithCancel(context.Background())

	go g.Consume(ctx, topics, handleMsgFn)

	<-time.After(time.Second)

	broker0.SetHandlerByMap(mockData)
	group, err = sarama.NewConsumerGroup([]string{broker0.Addr()}, myGroup, config)
	if err != nil {
		t.Fatal(err)
	}
	g.Group = group
	go g.ConsumeCustom(ctx, topics, &defaultConsumerHandler{
		ctx:              ctx,
		handleMessageFn:  handleMsgFn,
		zapLogger:        g.zapLogger,
		autoCommitEnable: g.autoCommitEnable,
	})

	<-time.After(time.Second)
	cancel()
}

func TestConsumerPartition(t *testing.T) {
	myTopic := "my-topic"
	testMsg := sarama.StringEncoder("Foo")
	broker0 := sarama.NewMockBroker(t, 0)

	manualOffset := int64(1234)
	offsetNewest := int64(2345)
	offsetNewestAfterFetchRequest := int64(3456)

	mockFetchResponse := sarama.NewMockFetchResponse(t, 1)

	mockFetchResponse.SetMessage(myTopic, 0, manualOffset-1, testMsg)

	for i := int64(0); i < 10; i++ {
		mockFetchResponse.SetMessage(myTopic, 0, i+manualOffset, testMsg)
	}

	mockFetchResponse.SetHighWaterMark(myTopic, 0, offsetNewestAfterFetchRequest)

	mockData := map[string]sarama.MockResponse{
		"MetadataRequest": sarama.NewMockMetadataResponse(t).
			SetBroker(broker0.Addr(), broker0.BrokerID()).
			SetLeader(myTopic, 0, broker0.BrokerID()),
		"OffsetRequest": sarama.NewMockOffsetResponse(t).
			SetOffset(myTopic, 0, sarama.OffsetOldest, 0).
			SetOffset(myTopic, 0, sarama.OffsetNewest, offsetNewest),
		"FetchRequest": mockFetchResponse,
	}
	broker0.SetHandlerByMap(mockData)

	master, err := sarama.NewConsumer([]string{broker0.Addr()}, sarama.NewConfig())
	if err != nil {
		t.Fatal(err)
	}

	c := &Consumer{
		C:         master,
		zapLogger: zap.NewExample(),
	}
	defer c.Close()

	ctx, cancel := context.WithCancel(context.Background())

	go c.ConsumePartition(ctx, myTopic, 0, manualOffset, handleMsgFn)
	<-time.After(time.Second)

	broker0.SetHandlerByMap(mockData)
	master, err = sarama.NewConsumer([]string{broker0.Addr()}, sarama.NewConfig())
	if err != nil {
		t.Fatal(err)
	}
	c.C = master
	go c.ConsumeAllPartition(ctx, myTopic, offsetNewest, handleMsgFn)
	<-time.After(time.Second)

	cancel()
}
