package eureka

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

type Client struct {
	config Config
}

type HttpAction struct {
	Method      string `yaml:"method"`
	Url         string `yaml:"url"`
	Body        string `yaml:"body"`
	Template    string `yaml:"template"`
	Accept      string `yaml:"accept"`
	ContentType string `yaml:"contentType"`
	Title       string `yaml:"title"`
	StoreCookie string `yaml:"storeCookie"`
}

type Config struct {
	Scheme     string
	Address    string
	HttpClient *http.Client
}

type request struct {
	config *Config
	method string
	url    *url.URL
	params url.Values
	body   io.Reader
	header http.Header
	obj    interface{}
	ctx    context.Context
}

// newRequest is used to create a new request
func (client *Client) newRequest(method, path string) *request {
	r := &request{
		config: &client.config,
		method: method,
		url: &url.URL{
			Scheme: client.config.Scheme,
			Host:   client.config.Address,
			Path:   path,
		},
		params: make(map[string][]string),
		header: make(http.Header),
	}
	return r
}

// toHTTP converts the request to an HTTP request
func (r *request) toHTTP() (*http.Request, error) {
	// Encode the query parameters
	r.url.RawQuery = r.params.Encode()

	// Check if we should encode the body
	if r.body == nil && r.obj != nil {
		b, err := encodeBody(r.obj)
		if err != nil {
			return nil, err
		}
		r.body = b
	}

	// Create the HTTP request
	req, err := http.NewRequest(r.method, r.url.RequestURI(), r.body)
	if err != nil {
		return nil, err
	}
	req.URL.Host = r.url.Host
	req.URL.Scheme = r.url.Scheme
	req.Host = r.url.Host
	req.Header = r.header
	return req, nil
}

// encodeBody is used to encode a request body
func encodeBody(obj interface{}) (io.Reader, error) {
	buf := bytes.NewBuffer(nil)
	enc := json.NewEncoder(buf)
	if err := enc.Encode(obj); err != nil {
		return nil, err
	}
	return buf, nil
}
