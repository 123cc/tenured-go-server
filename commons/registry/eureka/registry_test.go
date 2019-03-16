package eureka

import (
	"github.com/ihaiker/tenured-go-server/commons"
	"github.com/ihaiker/tenured-go-server/commons/registry"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var config *registry.PluginConfig

func init() {
	config, _ = registry.ParseConfig("eureka://127.0.0.1:8761")
}

type NLister struct {
}

func (this *NLister) OnNotify(status registry.RegistionStatus, serverInstances []registry.ServerInstance) {
	if status == registry.UNREGISTER {
		logrus.Info("OnNotify deregister: ", serverInstances)
	} else {
		logrus.Info("OnNotify register  : ", serverInstances)
	}
}

func TestEurekaServiceRegistry_Register(t *testing.T) {
	si, sr, err := registry.GetRegistry(config)
	assert.Nil(t, err)

	err = sr.Subscribe("test", &NLister{})
	t.Log(err)

	si.Id = "b102c658-830a-4d63-ba08-6a1ab75823d8"
	si.Name = "testApp"
	si.Address = "127.0.0.1:8761"
	si.Metadata = map[string]string{"instanceId": "user:password"}
	//修改配置，存活检测周期
	si.PluginAttrs.Config(map[string]string{"interval": "1s"})

	err = sr.Register(*si)
	assert.Nil(t, err)

	ss, err := sr.Lookup("testApp", nil)
	assert.Nil(t, err)
	for _, s := range ss {
		t.Log(s)
	}
	time.Sleep(time.Second * 5)

	err = sr.Unregister(si.Name)

	time.Sleep(time.Second * 5)
	(sr.(commons.Service)).Shutdown()

	assert.Nil(t, err)
}
