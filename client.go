// Author hoenig

package consulapi

import (
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/shoenig/toolkit"
)

const (
	defaultAddress = "http://localhost:8500"
	defaultTimeout = 10 * time.Second
)

type KV interface {
	Get(dc, path string) (string, error)
	Put(dc, path, value string) error
	Delete(dc, path string) error
	Keys(dc, path string) ([]string, error)
	Recurse(dc, path string) ([][2]string, error)
}

type Catalog interface {
	Datacenters() ([]string, error)
}

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

func (c *client) Datacenters() ([]string, error) {
	dcs := make([]string, 0, 10)
	if err := c.get("/v1/catalog/datacenters", &dcs); err != nil {
		return nil, err
	}
	return dcs, nil
}

func (c *client) Keys(dc, path string) ([]string, error) {
	path = fixup("/v1/kv", path, [2]string{"dc", dc}, [2]string{"keys", "true"})
	var keys []string
	err := c.get(path, &keys)
	sort.Strings(keys)
	return keys, err
}

func (c *client) Recurse(dc, path string) ([][2]string, error) {
	path = fixup("/v1/kv", path, [2]string{"dc", dc}, [2]string{"recurse", "true"})
	var values []value
	if err := c.get(path, &values); err != nil {
		return nil, err
	}
	kvs := make([][2]string, 0, len(values))
	for _, value := range values {
		decoded, err := base64.StdEncoding.DecodeString(value.Value)
		if err != nil {
			return nil, err
		}
		kvs = append(kvs, [2]string{
			value.Key,
			string(decoded),
		})
	}
	sort.Slice(kvs, func(i, j int) bool {
		return kvs[i][0] < kvs[j][0]
	})
	return kvs, nil
}

func (c *client) Get(dc, path string) (string, error) {
	path = fixup("/v1/kv", path, [2]string{"dc", dc})

	var values []value
	if err := c.get(path, &values); err != nil {
		return "", err
	}
	if len(values) == 0 {
		return "", errors.Errorf("key %q does not exist", path)
	}

	bs, err := base64.StdEncoding.DecodeString(values[0].Value)
	return string(bs), err
}

func (c *client) Put(dc, path, value string) error {
	path = fixup("/v1/kv", path, [2]string{"dc", dc})
	return c.put(path, value)
}

func (c *client) Delete(dc, path string) error {
	path = fixup("/v1/kv", path, [2]string{"dc", dc})
	return c.delete(path)
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

type value struct {
	Key   string `json:"Key"`
	Value string `json:"Value"`
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
