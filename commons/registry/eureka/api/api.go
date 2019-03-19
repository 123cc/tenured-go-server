package api

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/ihaiker/tenured-go-server/commons/registry"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

const (
	UP             = "UP"
	DOWN           = "DOWN"
	STARTING       = "STARTING"
	OUT_OF_SERVICE = "OUT_OF_SERVICE"
	UNKNOWN        = "UNKNOWN"
)

type Client struct {
	config Config
}

type Config struct {
	// Address is the address of the eureka server
	Address string
}

//Register new application instance
func (this *Client) Registry(eu *EurekaInstance) error {
	httpAction := toHttpAction("POST", this.config.Address, eu.App, "", "", false)
	httpAction.Body = toRegistry(eu)
	_, resp, err := doHttpRequest(httpAction)
	if err != nil {
		return err
	}
	if resp.StatusCode != 204 {
		var buf bytes.Buffer
		io.Copy(&buf, resp.Body)
		resp.Body.Close()
		return fmt.Errorf("eureka Registry unexpected response code: %d (%s)", resp.StatusCode, buf.Bytes())
	}
	return nil
}

//De-register application instance
func (this *Client) Deregister(appName string, instanceId string) error {
	httpAction := toHttpAction("DELETE", this.config.Address, appName, instanceId, "", false)
	_, resp, err := doHttpRequest(httpAction)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		var buf bytes.Buffer
		io.Copy(&buf, resp.Body)
		resp.Body.Close()
		return fmt.Errorf("eureka Deregister unexpected response code: %d (%s)", resp.StatusCode, buf.Bytes())
	}
	resp.Body.Close()
	return nil
}

//Send application instance heartbeat
func (this *Client) Heartbeat(appName string, instanceId string) error {
	httpAction := toHttpAction("PUT", this.config.Address, appName, instanceId, "", false)
	_, resp, err := doHttpRequest(httpAction)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		var buf bytes.Buffer
		io.Copy(&buf, resp.Body)
		resp.Body.Close()
		return fmt.Errorf("eureka Heartbeat unexpected response code: %d (%s)", resp.StatusCode, buf.Bytes())
	}
	resp.Body.Close()
	return nil
}

//Query for all instances
func (this *Client) QueryInstancesAll() ([]EurekaApplication, error) {
	var m EurekaApplicationsRootResponse
	httpAction := toHttpAction("GET", this.config.Address, "", "", "", true)
	bytes, err := executeQuery(httpAction)
	if err != nil {
		return nil, err
	} else {
		err := json.Unmarshal(bytes, &m)
		if err != nil {
			return nil, err
		}
		return m.Resp.Applications, nil
	}
}

//Query for all appID instances
func (this *Client) QueryInstancesToAppId(appName string) ([]RegistryInfo, error) {
	var m EurekaServiceResponse
	httpAction := toHttpAction("GET", this.config.Address, appName, "", "", true)
	bytes, err := executeQuery(httpAction)
	if err != nil {
		return nil, err
	} else {
		err := json.Unmarshal(bytes, &m)
		if err != nil {
			return nil, err
		}
		return m.Application.Instance, nil
	}
}

//Query for a specific appID/instanceID
func (this *Client) QueryInstancesToAppIdAndInstanceId(appName string, instanceId string) (string, error) {
	httpAction := toHttpAction("GET", this.config.Address, appName, instanceId, "", true)
	_, resp, err := doHttpRequest(httpAction)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	io.Copy(&buf, resp.Body)
	if resp.StatusCode != 200 {
		resp.Body.Close()
		return "", fmt.Errorf("eureka QueryInstancesToAppIdAndInstanceId unexpected response code: %d (%s)", resp.StatusCode, buf.Bytes())
	}
	return buf.String(), nil
}

//Query for a specific instanceID
func (this *Client) QueryInstancesToInstanceId(instanceId string) (string, error) {
	httpAction := toHttpAction("GET", this.config.Address, "", instanceId, "", true)
	_, resp, err := doHttpRequest(httpAction)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	io.Copy(&buf, resp.Body)
	if resp.StatusCode != 200 {
		resp.Body.Close()
		return "", fmt.Errorf("eureka QueryInstancesToInstanceId unexpected response code: %d (%s)", resp.StatusCode, buf.Bytes())
	}
	return buf.String(), nil
}

//Take instance out of service
func (this *Client) OutOfInstances(appName string, instanceId string) error {
	httpAction := toHttpAction("PUT", this.config.Address, appName, instanceId, "status?value=OUT_OF_SERVICE", false)
	_, resp, err := doHttpRequest(httpAction)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		var buf bytes.Buffer
		io.Copy(&buf, resp.Body)
		resp.Body.Close()
		return fmt.Errorf("eureka QueryInstancesToInstanceId unexpected response code: %d (%s)", resp.StatusCode, buf.Bytes())
	} else {
		return nil
	}
}

//Move instance back into service (remove override)
func (this *Client) RefreshInstances(appName string, instanceId string) error {
	httpAction := toHttpAction("DELETE", this.config.Address, appName, instanceId, "status?value=UP", false)
	_, resp, err := doHttpRequest(httpAction)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		var buf bytes.Buffer
		io.Copy(&buf, resp.Body)
		resp.Body.Close()
		return fmt.Errorf("eureka QueryInstancesToInstanceId unexpected response code: %d (%s)", resp.StatusCode, buf.Bytes())
	} else {
		return nil
	}
}

//Update metadata
func (this *Client) UpdateMetadata(appName string, instanceId string) error {
	httpAction := toHttpAction("PUT", this.config.Address, appName, instanceId, "metadata?key=value", false)
	_, resp, err := doHttpRequest(httpAction)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		var buf bytes.Buffer
		io.Copy(&buf, resp.Body)
		resp.Body.Close()
		return fmt.Errorf("eureka QueryInstancesToInstanceId unexpected response code: %d (%s)", resp.StatusCode, buf.Bytes())
	} else {
		return nil
	}
}

func doHttpRequest(request HttpAction) (time.Duration, *http.Response, error) {
	req := buildHttpRequest(request)
	var DefaultTransport http.RoundTripper = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	start := time.Now()
	resp, err := DefaultTransport.RoundTrip(req)
	diff := time.Since(start)
	if err != nil {
		log.Printf("HTTP request failed: %s", err.Error())
		return diff, nil, err
	}
	return diff, resp, err
}

func executeQuery(httpAction HttpAction) ([]byte, error) {
	req := buildHttpRequest(httpAction)

	var DefaultTransport http.RoundTripper = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	resp, err := DefaultTransport.RoundTrip(req)
	if err != nil {
		return []byte(nil), err
	} else {
		defer resp.Body.Close()
		responseBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return []byte(nil), err
		}
		return responseBody, nil
	}
}

func buildHttpRequest(httpAction HttpAction) *http.Request {
	var req *http.Request
	var err error
	if httpAction.Body != "" {
		reader := strings.NewReader(httpAction.Body)
		req, err = http.NewRequest(httpAction.Method, httpAction.Url, reader)
	} else {
		req, err = http.NewRequest(httpAction.Method, httpAction.Url, nil)
	}
	if err != nil {
		log.Fatal(err)
	}

	// Add headers
	req.Header = map[string][]string{
		"Accept":       {httpAction.Accept},
		"Content-Type": {httpAction.ContentType},
	}
	return req
}

func toRegistry(eu *EurekaInstance) string {
	instance := Instance{
		Instance: eu,
	}
	jsonBytes, err := json.Marshal(instance)
	if err != nil {
		return ""
	}
	return string(jsonBytes)
}

func toHttpAction(method string, ipAddr string, appName string, instanceId string, param string, isJson bool) HttpAction {
	httpAction := HttpAction{
		Method:      method,
		ContentType: "application/json;charset=UTF-8",
	}
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("http://%s/eureka/apps", ipAddr))
	if appName != "" {
		buffer.WriteString("/" + appName)
	}
	if instanceId != "" {
		buffer.WriteString("/" + instanceId)
	}
	if param != "" {
		buffer.WriteString("/" + param)
	}
	if isJson {
		httpAction.Accept = "application/json;charset=UTF-8"
	}
	httpAction.Url = buffer.String()
	return httpAction
}

func NewClient(ipAddr string) *Client {
	client := &Client{
		config: Config{
			Address: ipAddr,
		},
	}
	return client
}

//创建请求实体
func NewInstanceInfo(ipAddrs string, port int, serverInstance registry.ServerInstance) *EurekaInstance {
	Port := &Port{
		Port:    port,
		Enabled: true,
	}
	dataConterInfo := &DataCenterInfo{
		Name:  "MyOwn",
		Class: "com.netflix.appinfo.InstanceInfo$DefaultDataCenterInfo",
	}
	instanceInfo := &EurekaInstance{
		HostName:         ipAddrs,             // 主机名称ip地址
		App:              serverInstance.Name, // 服务名称
		IpAddr:           ipAddrs,
		Status:           UP,
		VipAddress:       serverInstance.Name,
		secureVipAddress: serverInstance.Name,
		Port:             Port, // 服务 IP:端口
		SecurePort:       Port,
		HomePageUrl:      fmt.Sprintf("http://%s:%d", ipAddrs, port),
		StatusPageUrl:    fmt.Sprintf("http://%s:%d/info", ipAddrs, port),
		HealthCheckUrl:   fmt.Sprintf("http://%s:%d/health", ipAddrs, port),
		DataCenterInfo:   dataConterInfo,
		Metadata:         serverInstance.Metadata,
	}
	return instanceInfo
}
