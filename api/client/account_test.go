package client

import (
	"github.com/ihaiker/tenured-go-server/api"
	"github.com/ihaiker/tenured-go-server/commons"
	"github.com/ihaiker/tenured-go-server/commons/registry"
	_ "github.com/ihaiker/tenured-go-server/commons/registry/consul"
	"github.com/ihaiker/tenured-go-server/commons/snowflake"
	"github.com/kataras/iris/core/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

var reg registry.ServiceRegistry
var server *AccountServiceClient

func Init() error {
	if config, err := registry.ParseConfig("consul://127.0.0.1:8500"); err != nil {
		return err
	} else if plugins, has := registry.GetPlugins(config.Plugin); !has {
		return errors.New("no registry")
	} else {
		if reg, err = plugins.Registry(*config); err != nil {
			return err
		}
	}
	server, _ = NewAccountServiceClient("tenured_store", reg)
	return nil
}

func TestNewAccount(t *testing.T) {
	defer func() {
		commons.ShutdownIfService(reg, true)
		server.Shutdown(true)
	}()
	err := Init()
	assert.Nil(t, err)

	err = server.Start()
	assert.Nil(t, err)

	id, _ := snowflake.NewSnowflake(snowflake.Settings{}).NextID()

	account := &api.Account{}
	account.Id = id
	account.Email = "wo@renzhen.la"

	if err := server.Apply(account); err != nil {
		t.Log(err)
	}

	ac, err := server.Get(account.Id)
	assert.Nil(t, err)

	t.Log(ac)
}
