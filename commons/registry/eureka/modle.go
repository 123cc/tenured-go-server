package eureka

type Instance struct {
	Instance *EurekaInstance `xml:"instance" json:"instance"`
}

type EurekaInstance struct {
	HostName         string          `xml:"hostName" json:"hostName"`
	App              string          `xml:"app" json:"app"`
	IpAddr           string          `xml:"ipAddr" json:"ipAddr"`
	VipAddress       string          `xml:"vipAddress" json:"vipAddress"`
	secureVipAddress string          `xml:"secureVipAddress" json:"secureVipAddress"`
	Status           string          `xml:"status" json:"status"`
	Port             *Port           `xml:"port" json:"port"`
	SecurePort       *Port           `xml:"securePort" json:"securePort"`
	HomePageUrl      string          `xml:"homePageUrl" json:"homePageUrl"`
	StatusPageUrl    string          `xml:"statusPageUrl" json:"statusPageUrl"`
	HealthCheckUrl   string          `xml:"healthCheckUrl" json:"healthCheckUrl"`
	DataCenterInfo   *DataCenterInfo `xml:"dataCenterInfo" json:"dataCenterInfo"`
	//optional
	LeaseInfo *LeaseInfo `xml:"leaseInfo" json:"leaseInfo"`
	//optional app specific metadata
	Metadata map[string]string `xml:"metadata" json:"metadata"`
}

type Port struct {
	Port    int  `xml:",chardata" json:"$"`
	Enabled bool `xml:"enabled,attr" json:"@enabled"`
}

type DataCenterInfo struct {
	Name  string `xml:"name" json:"name"`
	Class string `xml:"class,attr" json:"@class"`
	//Metadata *DataCenterMetadata `xml:"metadata" json:"metadata"`
}

type DataCenterMetadata struct {
	//ami发布索引
	AmiLaunchIndex string `xml:"ami-launch-index" json:"ami-launch-index"`
	//本地主机名
	LocalHostname string `xml:"local-hostname" json:"local-hostname"`
	//有效区间
	AvailabilityZone string `xml:"availability-zone" json:"availability-zone"`
	//实例化ID
	InstanceId     string `xml:"instance-id" json:"instance-id"`
	PublicIpv4     string `xml:"public-ipv4" json:"public-ipv4"`
	PublicHostname string `xml:"public-hostname" json:"public-hostname"`
	//AMI清单路径
	AmiManifestPath string `xml:"ami-manifest-path" json:"ami-manifest-path"`
	LocalIpv4       string `xml:"local-ipv4" json:"local-ipv4"`
	Hostname        string `xml:"hostname" json:"hostname"`
	AmiId           string `xml:"ami-id" json:"ami-id"`
	//实例化类型
	InstanceType string `xml:"instance-type" json:"instance-type"`
}

type LeaseInfo struct {
	//收回持续时间插入
	EvictionDurationInSecs uint `xml:"evictionDurationInSecs" json:"evictionDurationInSecs"`
}

type EurekaServiceResponse struct {
	Application EurekaApplication `json:"application"`
}

type EurekaApplication struct {
	Name     string         `json:"name"`
	Instance []RegistryInfo `json:"instance"`
}

type RegistryInfo struct {
	HostName string     `json:"hostName"`
	Port     EurekaPort `json:"port"`
}

type EurekaPort struct {
	Port int `json:"$"`
}

type EurekaApplicationsRootResponse struct {
	Resp EurekaApplicationsResponse `json:"applications"`
}

type EurekaApplicationsResponse struct {
	Version      string              `json:"versions__delta"`
	AppsHashcode string              `json:"versions__delta"`
	Applications []EurekaApplication `json:"application"`
}
