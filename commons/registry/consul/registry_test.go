package consul

import (
	"github.com/ihaiker/tenured-go-server/commons"
	"github.com/ihaiker/tenured-go-server/commons/registry"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var config *registry.PluginConfig

func init() {
	config, _ = registry.ParseConfig("consul://127.0.0.1:8500")
}

type NLister struct {
}

func (this *NLister) OnNotify(status registry.RegistionStatus, serverInstances []registry.ServerInstance) {
	//logrus.Debug(status)
}

func TestConsulServiceRegistry_Register(t *testing.T) {
	si, sr, err := registry.GetRegistry(config)
	assert.Nil(t, err)

	err = sr.Subscribe("test", &NLister{})
	t.Log(err)

	si.Id = "b102c658-830a-4d63-ba08-6a1ab75823d8"
	si.Name = "test"
	si.Address = "127.0.0.1:6071"
	si.Metadata = map[string]string{"test_metadata": "demo"}
	//修改配置，存活检测周期
	si.PluginAttrs.Config(map[string]string{"interval": "1s"})

	err = sr.Register(*si)
	assert.Nil(t, err)

	ss, err := sr.Lookup("test", nil)
	assert.Nil(t, err)
	for _, s := range ss {
		t.Log(s)
	}
	time.Sleep(time.Second * 5)

	t.Log("服务下线....")
	err = sr.Unregister(si.Id)

	time.Sleep(time.Second * 5)
	(sr.(commons.Service)).Shutdown()

	assert.Nil(t, err)
}
