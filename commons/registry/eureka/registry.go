package eureka

import (
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/ihaiker/tenured-go-server/commons/registry"
	"github.com/sirupsen/logrus"
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
	client *Client
	config *EurekaConfig

	subscribes map[string]*subscribeInfo
}

func (this *EurekaServiceRegistry) Start() error {
	return nil
}

func (this *EurekaServiceRegistry) Shutdown() {
	for name, ch := range this.subscribes {
		ch.close()
		delete(this.subscribes, name)
	}
}

//注册中心注册服务
func (this *EurekaServiceRegistry) Register(serverInstance registry.ServerInstance) error {
	logrus.Infof("register %s(%s) : %s", serverInstance.Name, serverInstance.Address, serverInstance.Id)
	if _, portStr, err := net.SplitHostPort(serverInstance.Address); err != nil {
		return err
	} else if port, err := strconv.Atoi(portStr); err != nil {
		return err
	} else {
		reg := NewInstanceInfo("169.254.67.217", port, serverInstance)
		return this.client.Agent().Registry(reg)
	}
	return nil
}

func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	panic("Unable to determine local IP address (non loopback). Exiting.")
}

//从注册中心删除注册
func (this *EurekaServiceRegistry) Unregister(appName string) error {
	logrus.Info("Unregister ", appName)
	return this.client.Agent().Deregister(appName, appName)
}

//订阅服务改变
func (this *EurekaServiceRegistry) Subscribe(serverName string, listener registry.RegistryNotifyListener) error {
	if this.addSubscribe(serverName, listener) {
		go this.loadSubscribeHealth(serverName)
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
	if services, err := this.client.Agent().QueryInstancesToAppId(serverName); err != nil {
		return nil, err
	} else {
		serverInstances := make([]registry.ServerInstance, len(services))
		for i := 0; i < len(services); i++ {
			serverInstances[i] = this.convertService(serverName, services[i])
		}
		return serverInstances, nil
	}
}

func (this *EurekaServiceRegistry) convertService(serverName string, service RegistryInfo) registry.ServerInstance {
	return registry.ServerInstance{
		Name:    serverName,
		Address: service.HostName,
		Status:  "OK",
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
}

func newRegistry(pluginConfig *registry.PluginConfig) (*EurekaServiceRegistry, error) {
	config := &EurekaConfig{config: pluginConfig}
	serviceRegistry := &EurekaServiceRegistry{
		config:     config,
		subscribes: map[string]*subscribeInfo{},
	}
	defConfig := &Config{
		Address: "127.0.0.1:8761",
		Scheme:  "http",
	}
	defConfig.Scheme = config.Scheme()
	defConfig.Address = config.Address()
	client := &Client{
		config: *defConfig,
	}
	serviceRegistry.client = client
	return serviceRegistry, nil
}
