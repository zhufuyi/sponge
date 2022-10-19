package etcd

import (
	"encoding/json"

	"github.com/zhufuyi/sponge/pkg/servicerd/registry"
)

func marshal(si *registry.ServiceInstance) (string, error) {
	data, err := json.Marshal(si)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func unmarshal(data []byte) (si *registry.ServiceInstance, err error) {
	err = json.Unmarshal(data, &si)
	return
}
