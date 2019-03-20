package api

type Instance struct {
	Instance *EurekaInstance `xml:"instance" json:"instance"`
}

type HttpAction struct {
	Method      string `yaml:"method"`
	Url         string `yaml:"url"`
	Body        string `yaml:"body"`
	Accept      string `yaml:"accept"`
	ContentType string `yaml:"contentType"`
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
	LeaseInfo *LeaseInfo `xml:"leaseInfo,omitempty" json:"leaseInfo,omitempty"`
	//optional app specific metadata
	Metadata map[string]string `xml:"metadata,omitempty" json:"metadata,omitempty"`
}

type Port struct {
	Port    int  `xml:",chardata" json:"$"`
	Enabled bool `xml:"enabled,attr" json:"@enabled,omitempty"`
}

type DataCenterInfo struct {
	Name     string              `xml:"name" json:"name"`
	Class    string              `xml:"class,attr" json:"@class"`
	Metadata *DataCenterMetadata `xml:"metadata,omitempty" json:"metadata,omitempty"`
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
	EvictionDurationInSecs uint `xml:"evictionDurationInSecs,omitempty" json:"evictionDurationInSecs,omitempty"`
	RenewalIntervalInSecs  int  `xml:"renewalIntervalInSecs,omitempty" json:"renewalIntervalInSecs,omitempty"`
	DurationInSecs         int  `xml:"durationInSecs,omitempty" json:"durationInSecs,omitempty"`
	RegistrationTimestamp  int  `xml:"registrationTimestamp,omitempty" json:"registrationTimestamp,omitempty"`
	LastRenewalTimestamp   int  `xml:"lastRenewalTimestamp,omitempty" json:"lastRenewalTimestamp,omitempty"`
	EvictionTimestamp      int  `xml:"evictionTimestamp,omitempty" json:"evictionTimestamp,omitempty"`
	ServiceUpTimestamp     int  `xml:"serviceUpTimestamp,omitempty" json:"serviceUpTimestamp,omitempty"`
}

type EurekaServiceResponse struct {
	Application EurekaApplication `json:"application"`
}

type EurekaApplication struct {
	Name     string         `json:"name"`
	Instance []RegistryInfo `json:"instance"`
}

type RegistryInfo struct {
	InstanceId       string            `json:"instanceId"`
	HostName         string            `json:"hostName"`
	App              string            `json:"app"`
	IpAddr           string            `json:"ipAddr"`
	Status           string            `json:"status"`
	Overriddenstatus string            `json:"overriddenstatus"`
	Port             *PortResult       `json:"port,omitempty"`
	SecurePort       *PortResult       `json:"securePort"`
	CountryId        int               `json:"countryId"`
	DataCenterInfo   *DataCenterInfo   `json:"dataCenterInfo"`
	LeaseInfo        *LeaseInfo        `json:"leaseInfo"`
	Metadata         map[string]string `json:"metadata"`
	HomePageUrl      string            `json:"homePageUrl"`
	StatusPageUrl    string            `json:"statusPageUrl"`
	HealthCheckUrl   string            `json:"healthCheckUrl"`
	VipAddress       string            `json:"vipAddress"`
}
type PortResult struct {
	Port    int  `xml:",chardata" json:"$"`
	Enabled bool `xml:"enabled,attr" json:"enabled"`
}

type EurekaApplicationsRootResponse struct {
	Resp EurekaApplicationsResponse `json:"applications"`
}

type EurekaApplicationsResponse struct {
	Version      string              `json:"versions__delta"`
	AppsHashcode string              `json:"versions__delta"`
	Applications []EurekaApplication `json:"application"`
}
