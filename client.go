package consulapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/hashicorp/go-hclog"
	"gophers.dev/pkgs/ignore"
)

const (
	defaultScheme = "http"
	defaultHost   = "localhost"
	defaultPort   = 8500
)

//go:generate go run github.com/gojuno/minimock/v3/cmd/minimock -g -i Client -s _mock.go

// A Client is the abstraction of ConsulClient, a tool used to communicate with consul.
// The interface is composed of other interfaces, which reflect the different categories
// of API supported by the consul agent.
type Client interface {
	Agent
	// Catalog
	// KV
	// Session
	// Candidate
}

// ConsulClient is used to communicate with Consul.
type ConsulClient struct {
	scheme     string
	host       string
	port       int
	userAgent  string
	headers    http.Header
	httpClient *http.Client
	log        hclog.Logger
}

// ClientOption implements functional options for ConsulClient.
type ClientOption func(*ConsulClient)

func WithDefaultScheme(scheme string) ClientOption {
	return func(c *ConsulClient) {
		c.scheme = scheme
	}
}

func WithDefaultHost(host string) ClientOption {
	return func(c *ConsulClient) {
		c.host = host
	}
}

func WithDefaultPort(port int) ClientOption {
	return func(c *ConsulClient) {
		c.port = port
	}
}

func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *ConsulClient) {
		c.httpClient = httpClient
	}
}

func WithLogger(log hclog.Logger) ClientOption {
	return func(c *ConsulClient) {
		c.log = log
	}
}

func WithUserAgent(ua string) ClientOption {
	return func(c *ConsulClient) {
		c.userAgent = ua
	}
}

func WithHeaders(h http.Header) ClientOption {
	return func(c *ConsulClient) {
		c.headers = h
	}
}

const (
	headerContentType = "Content-Type"
	mimeJSON          = "application/json"

	headerUserAgent = "User-Agent"
	userAgent       = "consulapi/v2"

	headerConsulToken = "X-Consul-Token"
)

func defaultHeaders() http.Header {
	var h http.Header
	h.Set(headerContentType, mimeJSON)
	h.Set(headerUserAgent, userAgent)
	return h
}

// New creates a new Client that will use the provided ClientOptions for
// making requests to a configured consul agent.
func New(opts ...ClientOption) Client {
	c := &ConsulClient{
		scheme:     defaultScheme,
		host:       defaultHost,
		port:       defaultPort,
		headers:    defaultHeaders(),
		httpClient: http.DefaultClient,
		log:        hclog.NewNullLogger(),
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// params are url param kv pairs
func fixup(prefix, path string, params ...[2]string) string {
	prefix = strings.TrimSuffix(prefix, "/")
	path = strings.TrimPrefix(path, "/")

	values := make(url.Values)
	for _, parameter := range params {
		if parameter[1] != "" {
			values.Add(parameter[0], parameter[1])
		}
	}

	query := values.Encode()

	// there is a better way to build url queries
	completeURL := prefix + "/" + path
	if len(query) > 0 {
		completeURL += "?" + query
	}
	return completeURL
}

func param(key, value string) [2]string {
	return [2]string{key, value}
}

func (cc *ConsulClient) newRequest(ctx Ctx, method, uri, token string, body io.Reader) (*http.Request, error) {
	request, err := http.NewRequest(method, uri, body)
	if err != nil {
		return nil, err
	}
	rCtx := request.WithContext(ctx)
	cc.setHeaders(rCtx, token)
	return rCtx, nil
}

func (cc *ConsulClient) get(path, i interface{}, opts *Optional) error {
	// completeURL := cc.address + path
	// create url

	host := cc.address
	if opts.host != "" {
		host = opts.host
	}

	u := url.URL{
		Scheme:      "",
		Host:        host,
		Path:        "",
		RawPath:     "",
		ForceQuery:  false,
		RawQuery:    "",
		Fragment:    "",
		RawFragment: "",
	}

	request, err := cc.newRequest(ctx, http.MethodGet, completeURL, token, nil)
	if err != nil {
		return err
	}

	response, err := cc.httpClient.Do(request)
	if err != nil {
		return err
	}
	defer ignore.Drain(response.Body)

	if response.StatusCode >= 400 {
		return &RequestError{statusCode: response.StatusCode}
	}

	return json.NewDecoder(response.Body).Decode(i)
}

func (cc *ConsulClient) put(ctx Ctx, path, token, body string, i interface{}) error {
	completeURL := cc.address + path

	r := strings.NewReader(body)
	request, err := cc.newRequest(ctx, http.MethodPut, completeURL, token, r)
	if err != nil {
		return err
	}

	response, err := cc.httpClient.Do(request)
	if err != nil {
		return err
	}
	defer ignore.Drain(response.Body)

	if response.StatusCode >= 400 {
		return &RequestError{statusCode: response.StatusCode}
	}

	if i != nil {
		return json.NewDecoder(response.Body).Decode(i)
	}

	return nil
}

func (cc *ConsulClient) delete(ctx Ctx, path, token string) error {
	completeURL := cc.address + path

	request, err := cc.newRequest(ctx, http.MethodDelete, completeURL, token, nil)
	if err != nil {
		return err
	}

	response, err := cc.httpClient.Do(request)
	if err != nil {
		return err
	}
	defer ignore.Drain(response.Body)

	if response.StatusCode >= 400 {
		return &RequestError{statusCode: response.StatusCode}
	}

	return nil
}

func (cc *ConsulClient) setHeaders(request *http.Request, token string) {
	// set common default headers
	request.Header.Set(headerContentType, mimeJSON)
	request.Header.Set(headerUserAgent, userAgent)

	// set acl token header if it exists
	if token != "" {
		request.Header.Set(consulTokenHeader, token)
	}

	// override using specified headers
	for k, values := range cc.headers {
		for _, v := range values {
			request.Header.Add(k, v)
		}
	}
}

// RequestError exposes the status code of a http request error
type RequestError struct {
	statusCode int
}

func (h *RequestError) Error() string {
	return fmt.Sprintf("status code (%d)", h.statusCode)
}

func (h *RequestError) StatusCode() int {
	return h.statusCode
}
