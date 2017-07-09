// Author hoenig

package consulapi

import (
	"crypto/tls"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/shoenig/toolkit"
)

// mocks generated with github.com/vektra/mockery
//go:generate mockery -name Client -case=underscore -outpkg consulapitest -output consulapitest

const (
	defaultAddress    = "http://localhost:8500"
	defaultTimeout    = 10 * time.Second
	consulTokenHeader = "X-Consul-Token"
)

// A Client is used to communicate with consul. The interface is composed of
// other interfaces, which reflect the different categories of API supported by
// the consul agent.
type Client interface {
	Agent
	Catalog
	KV
}

// ClientOptions are used to configure options of a client upon creation.
type ClientOptions struct {
	// Address of the consul agent to communicate with. This value will
	// default to http://localhost:8500 if left unset. This is likely
	// the desired value, as consul is designed to run with an agent on
	// every node.
	Address string

	// HTTPTimeout configures how long underlying HTTP requests should wait
	// before giving up and returning a timeout error. By default, this value
	// is 10 seconds.
	HTTPTimeout time.Duration

	// SkipTLSVerification configures the underlying HTTP client to ignore
	// any TLS certificate validation errors. This is a hacky option that can
	// be useful for working in environments that are using self-signed
	// certificates. For best security practices, this option should never
	// be used in a production environment.
	SkipTLSVerification bool

	// If consul is configured to authenticate requests with a token,
	// set the value of that token here.
	Token string
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
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	values := make(url.Values)

	for _, param := range params {
		if param[1] != "" {
			values.Add(param[0], param[1])
		}
	}

	query := values.Encode()

	// there is probably a better way to build url queries
	url := prefix + path
	if len(query) > 0 {
		url += "?" + query
	}
	return url
}

func (c *client) get(path string, i interface{}) error {
	url := c.opts.Address + path

	request, err := http.NewRequest(http.MethodGet, url, nil)
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
		return errors.Errorf("bad status code: %d", response.StatusCode)
	}

	return json.NewDecoder(response.Body).Decode(i)
}

func (c *client) put(path, body string) error {
	url := c.opts.Address + path

	request, err := http.NewRequest(http.MethodPut, url, strings.NewReader(body))
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
		return errors.Errorf("bad status code: %d", response.StatusCode)
	}
	return nil
}

func (c *client) delete(path string) error {
	url := c.opts.Address + path

	request, err := http.NewRequest(http.MethodDelete, url, nil)
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
		return errors.Errorf("bad status code: %d", response.StatusCode)
	}
	return nil
}

func (c *client) maybeSetToken(request *http.Request) {
	if c.opts.Token != "" {
		request.Header.Set(consulTokenHeader, c.opts.Token)
	}
}
