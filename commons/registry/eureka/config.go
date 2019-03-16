package eureka

import (
	"github.com/ihaiker/tenured-go-server/commons/registry"
	"time"
)

type EurekaConfig struct {
	config *registry.PluginConfig
}

func (this *EurekaConfig) Scheme() string {
	return this.config.Get("scheme", "http")
}

func (this *EurekaConfig) Address() string {
	return this.config.Address[0]
}

func (this *EurekaConfig) Datacenter() string {
	return this.config.Get("datacenter", "dcl")
}

func (this *EurekaConfig) Token() string {
	return this.config.Get("token", "")
}

func (this *EurekaConfig) HealthWaiTime() time.Duration {
	return time.Second * time.Duration(this.config.GetInt("healthWaiTime", 5))
}

func (this *EurekaConfig) HealthFailTime() time.Duration {
	return time.Second * time.Duration(this.config.GetInt("failHealthTime", 1))
}
