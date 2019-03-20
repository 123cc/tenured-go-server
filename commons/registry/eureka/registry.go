package eureka

import (
	"fmt"
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/ihaiker/tenured-go-server/commons/registry"
	"github.com/sirupsen/logrus"
	"github.com/tenured-go-server/commons/registry/eureka/api"
	"net"
	"strconv"
)

//创建订阅详情结构体
type subscribeInfo struct {
	listeners *hashset.Set
	services  map[string]registry.ServerInstance
	closeChan chan struct{}
}

//为结构体添加关闭方法
func (subscribeInfo *subscribeInfo) close() {
	close(subscribeInfo.closeChan)
}

type EurekaServiceRegistry struct {
	client     *api.Client
	config     *EurekaConfig
	subscribes map[string]*subscribeInfo
}

func (this *EurekaServiceRegistry) Start() error {
	return nil
}

func (this *EurekaServiceRegistry) Shutdown(interrupt bool) {
	for name, ch := range this.subscribes {
		ch.close()
		delete(this.subscribes, name)
	}
}

//注册中心注册服务
func (this *EurekaServiceRegistry) Register(serverInstance registry.ServerInstance) error {
	logrus.Infof("register %s(%s) : %s", serverInstance.Name, serverInstance.Address, serverInstance.Id)
	attrs := serverInstance.PluginAttrs.(*EurekaServiceAttrs)
	if host, portStr, err := net.SplitHostPort(serverInstance.Address); err != nil {
		return err
	} else if port, err := strconv.Atoi(portStr); err != nil {
		return err
	} else {
		timeOut, terr := strconv.Atoi(attrs.RequestTimeout)
		renewal, rerr := strconv.Atoi(attrs.Interval)
		eviction, err := strconv.Atoi(attrs.Deregister)
		if err != nil {
			return err
		}
		if rerr != nil {
			return rerr
		}
		if terr != nil {
			return terr
		}
		leaseInfo := &api.LeaseInfo{
			EvictionDurationInSecs: uint(timeOut),
			RenewalIntervalInSecs:  renewal,
			DurationInSecs:         eviction,
		}
		reg := api.NewInstanceInfo(host, port, serverInstance)
		reg.LeaseInfo = leaseInfo
		return this.client.Registry(reg)
	}
	return nil
}

//从注册中心删除注册
func (this *EurekaServiceRegistry) Unregister(appName string) error {
	logrus.Info("Unregister ", appName)
	services, err := this.client.QueryInstancesToAppId(appName)
	if err != nil {
		return err
	}
	return this.client.Deregister(appName, services[0].InstanceId)
}

//订阅服务改变
func (this *EurekaServiceRegistry) Subscribe(serverName string, listener registry.RegistryNotifyListener) error {
	if this.addSubscribe(serverName, listener) {
		//go this.loadSubscribeHealth(serverName)
	}
	return nil
}

//取消服务订阅
func (this *EurekaServiceRegistry) Unsubscribe(serverName string, listener registry.RegistryNotifyListener) error {
	if this.removeSubscribe(serverName, listener) {
		if sub, has := this.subscribes[serverName]; has {
			sub.close()
			delete(this.subscribes, serverName)
		}
	}
	return nil
}

//@return 是否是次服务的最后一个监听器
func (this *EurekaServiceRegistry) removeSubscribe(name string, listener registry.RegistryNotifyListener) bool {
	sets := this.getOrCreateSubscribe(name)
	sets.listeners.Remove(listener)
	return sets.listeners.Size() == 0
}

//发现服务内容
func (this *EurekaServiceRegistry) Lookup(serverName string, tags []string) ([]registry.ServerInstance, error) {
	if services, err := this.client.QueryInstancesToAppId(serverName); err != nil {
		return nil, err
	} else {
		serverInstances := make([]registry.ServerInstance, len(services))
		for i := 0; i < len(services); i++ {
			serverInstances[i] = this.convertService(serverName, services[i])
		}
		return serverInstances, nil
	}
}

func (this *EurekaServiceRegistry) convertService(serverName string, service api.RegistryInfo) registry.ServerInstance {
	var status = service.Status
	if status == api.UP {
		status = "OK"
	}
	tags := []string{service.VipAddress}
	return registry.ServerInstance{
		Id:       service.InstanceId,
		Name:     serverName,
		Metadata: service.Metadata,
		Address:  fmt.Sprintf("%s:%d", service.IpAddr, service.Port.Port),
		Tags:     tags,
		Status:   status,
	}
}

//@return 返回是否是此服务的第一个监听器
func (this *EurekaServiceRegistry) addSubscribe(name string, listener registry.RegistryNotifyListener) bool {
	sets := this.getOrCreateSubscribe(name)
	sets.listeners.Add(listener)
	return sets.listeners.Size() == 1
}
func (this *EurekaServiceRegistry) getOrCreateSubscribe(name string) *subscribeInfo {
	if subInfo, has := this.subscribes[name]; !has {
		subInfo = &subscribeInfo{
			listeners: hashset.New(),
			services:  nil,
			closeChan: make(chan struct{}),
		}
		this.subscribes[name] = subInfo
	}
	return this.subscribes[name]
}

func (this *EurekaServiceRegistry) loadSubscribeHealth(serverName string) {
	defer func() {
		if e := recover(); e != nil {
			logrus.Warnf("close subscribe(%s) error: %v", serverName, e)
		}
	}()
	logrus.Debug("start loop load subscribe server health:", serverName)

	register := make([]registry.ServerInstance, 0)
	deregister := make([]registry.ServerInstance, 0)

	for {
		subInfo, has := this.subscribes[serverName]
		if !has {
			return
		}
		select {
		case <-subInfo.closeChan:
			return
		default:
			services, err := this.client.QueryInstancesToAppId(serverName)
			if err != nil {
				continue
			}

			subInfo, has = this.subscribes[serverName]
			if !has {
				return
			}

			register = register[:0]
			deregister = deregister[:0]
			if subInfo.services == nil {
				subInfo.services = map[string]registry.ServerInstance{}
				for _, s := range services {
					subInfo.services[s.InstanceId] = this.convertService(serverName, s)
				}
			} else {
				currentServices := map[string]registry.ServerInstance{}

				for _, s := range services {
					current := this.convertService(serverName, s)
					if old, has := subInfo.services[s.InstanceId]; !has || current.Status != old.Status {
						register = append(register, current)
					}
					currentServices[s.InstanceId] = current
					delete(subInfo.services, s.InstanceId)
				}

				for _, s := range subInfo.services {
					s.Status = "deregister"
					deregister = append(deregister, s)
				}

				for _, v := range subInfo.listeners.Values() {
					if len(register) > 0 {
						v.(registry.RegistryNotifyListener).
							OnNotify(registry.REGISTER, register)
					}
					if len(deregister) > 0 {
						v.(registry.RegistryNotifyListener).
							OnNotify(registry.UNREGISTER, deregister)
					}
				}
				subInfo.services = currentServices
			}
		}
	}
}

func newRegistry(pluginConfig *registry.PluginConfig) (*EurekaServiceRegistry, error) {
	config := &EurekaConfig{config: pluginConfig}
	serviceRegistry := &EurekaServiceRegistry{
		config:     config,
		subscribes: map[string]*subscribeInfo{},
	}
	serviceRegistry.client = api.NewClient(config.Address())
	return serviceRegistry, nil
}
