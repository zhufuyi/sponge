package metrics

import (
	"context"
	"testing"
	"time"

	"github.com/zhufuyi/sponge/pkg/utils"

	"github.com/stretchr/testify/assert"
)

func TestClientHTTPService(t *testing.T) {
	serverAddr, _ := utils.GetLocalHTTPAddrPairs()

	s := ClientHTTPService(serverAddr)
	ctx, _ := context.WithTimeout(context.Background(), time.Second)
	time.Sleep(time.Millisecond * 100)
	err := s.Shutdown(ctx)
	assert.NoError(t, err)
}

func TestStreamClientMetrics(t *testing.T) {
	metrics := StreamClientMetrics()
	assert.NotNil(t, metrics)
}

func TestUnaryClientMetrics(t *testing.T) {
	metrics := UnaryClientMetrics()
	assert.NotNil(t, metrics)
}

func Test_cliRegisterMetrics(t *testing.T) {
	cliRegisterMetrics()
}
