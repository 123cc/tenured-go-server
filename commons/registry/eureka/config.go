package eureka

import (
	"github.com/ihaiker/tenured-go-server/commons/registry"
)

type EurekaConfig struct {
	config *registry.PluginConfig
}

func (this *EurekaConfig) Address() string {
	return this.config.Address[0]
}
