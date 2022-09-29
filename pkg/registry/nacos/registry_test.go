package nacos

import (
	"testing"

	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	namingClient, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  nil,
			ServerConfigs: nil,
		},
	)
	assert.Error(t, err)

	r := New(namingClient,
		WithPrefix("/micro"),
		WithWeight(1),
		WithCluster("cluster"),
		WithGroup("dev"),
		WithDefaultKind("grpc"),
	)
	assert.NotNil(t, r)
}
