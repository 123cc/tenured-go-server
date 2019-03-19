package eureka

import "github.com/ihaiker/tenured-go-server/commons/registry"

type EurekaServiceAttrs struct {
	//检查类型
	Health string `json:"health" yaml:"health"` //http url

	//检查频率，
	Interval string `json:"interval" yaml:"interval"` //10s

	//当前节点出现异常多长时间下线
	Deregister string `json:"deregister" yaml:"deregister"`

	//请求处理超时时间
	RequestTimeout string `json:"request_timeout" yaml:"request_timeout"`
}

func (this *EurekaServiceAttrs) Config(attrs map[string]string) {
	registry.LoadModel(this, attrs)
}

func newInstance() *EurekaServiceAttrs {
	return &EurekaServiceAttrs{
		Health:         "/health",
		Interval:       "10",
		Deregister:     "90",
		RequestTimeout: "3",
	}
}
