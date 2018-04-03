// Author hoenig

package consulapi

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/shoenig/toolkit"
)

const (
	defaultAddress    = "http://localhost:8500"
	defaultTimeout    = 10 * time.Second
	consulTokenHeader = "X-Consul-Token"
)

//go:generate mockery -interface Client -package consulapitest

// A Client is used to communicate with consul. The interface is composed of
// other interfaces, which reflect the different categories of API supported by
// the consul agent.
type Client interface {
	Agent
	Catalog
	KV
	Session
	Candidate
}

// ClientOptions are used to configure options of a client upon creation.
type ClientOptions struct {
	// Address (optional) of the consul agent to communicate with. This value
	// will default to http://localhost:8500 if left unset. This is likely
	// the desired value, as consul is designed to run with an agent on
	// every node.
	Address string

	// HTTPTimeout (optional) configures how long underlying HTTP requests should
	// wait before giving up and returning a timeout error. By default, this value
	// is 10 seconds.
	HTTPTimeout time.Duration

	// SkipTLSVerification configures the underlying HTTP client to ignore
	// any TLS certificate validation errors. This is a hacky option that can
	// be useful for working in environments that are using self-signed
	// certificates. For best security practices, this option should never
	// be used in a production environment.
	SkipTLSVerification bool

	// Token (optional) will be used to authenticate requests to consul.
	Token string

	// Logger may be optionally configured as an output for trace level logging
	// produced internally by the Client. This can be helpful for debugging logic
	// errors in client code.
	Logger *log.Logger
}

// RequestError exposes the status code of a http request error
type RequestError struct {
	statusCode int
}

func (h *RequestError) Error() string {
	return fmt.Sprintf("bad status code: %d", h.statusCode)
}

func (h *RequestError) StatusCode() int {
	return h.statusCode
}

// New creates a new Client that will connect to the configured consul
// agent.
func New(opts ClientOptions) Client {
	if opts.Address == "" {
		opts.Address = defaultAddress
	}

	if opts.HTTPTimeout <= 0 {
		opts.HTTPTimeout = defaultTimeout
	}

	if opts.Logger == nil {
		opts.Logger = log.New(ioutil.Discard, "", 0)
	}

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: opts.SkipTLSVerification,
		},
	}

	return &client{
		opts: opts,
		httpClient: &http.Client{
			Transport: transport,
			Timeout:   opts.HTTPTimeout,
		},
	}
}

type client struct {
	opts       ClientOptions
	httpClient *http.Client
}

// params are url param kv pairs
func fixup(prefix, path string, params ...[2]string) string {
	prefix = strings.TrimSuffix(prefix, "/")
	path = strings.TrimPrefix(path, "/")

	values := make(url.Values)
	for _, param := range params {
		if param[1] != "" {
			values.Add(param[0], param[1])
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

func (c *client) get(path string, i interface{}) error {
	completeURL := c.opts.Address + path

	request, err := http.NewRequest(http.MethodGet, completeURL, nil)
	if err != nil {
		return err
	}

	c.maybeSetToken(request)

	response, err := c.httpClient.Do(request)
	if err != nil {
		return err
	}
	defer toolkit.Drain(response.Body)

	if response.StatusCode >= 400 {
		return &RequestError{statusCode: response.StatusCode}
	}

	return json.NewDecoder(response.Body).Decode(i)
}

func (c *client) put(path, body string, i interface{}) error {
	completeURL := c.opts.Address + path

	request, err := http.NewRequest(http.MethodPut, completeURL, strings.NewReader(body))
	if err != nil {
		return err
	}

	c.maybeSetToken(request)

	response, err := c.httpClient.Do(request)
	if err != nil {
		return err
	}

	if response.StatusCode >= 400 {
		return &RequestError{statusCode: response.StatusCode}
	}

	if i != nil {
		defer toolkit.Drain(response.Body)
		return json.NewDecoder(response.Body).Decode(i)
	}

	return nil
}

func (c *client) delete(path string) error {
	completeURL := c.opts.Address + path

	request, err := http.NewRequest(http.MethodDelete, completeURL, nil)
	if err != nil {
		return err
	}

	c.maybeSetToken(request)

	response, err := c.httpClient.Do(request)
	if err != nil {
		return err
	}
	// do not read response

	if response.StatusCode >= 400 {
		return &RequestError{statusCode: response.StatusCode}
	}
	return nil
}

func (c *client) maybeSetToken(request *http.Request) {
	if c.opts.Token != "" {
		request.Header.Set(consulTokenHeader, c.opts.Token)
	}
}
