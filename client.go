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

const (
	defaultAddress = "http://localhost:8500"
	defaultTimeout = 10 * time.Second
)

type Client interface {
	Catalog
	KV
}

type ClientOptions struct {
	Address             string
	HTTPTimeout         time.Duration
	SkipTLSVerification bool
}

func New(opts ClientOptions) Client {
	if opts.Address == "" {
		opts.Address = defaultAddress
	}

	if opts.HTTPTimeout == 0 {
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
			values.Set(param[0], param[1])
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
	response, err := c.httpClient.Get(url)
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
