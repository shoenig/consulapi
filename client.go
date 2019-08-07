package consulapi // import "gophers.dev/pkgs/consulapi"

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	clean "github.com/hashicorp/go-cleanhttp"

	"gophers.dev/pkgs/ignore"
	"gophers.dev/pkgs/loggy"
)

const (
	defaultAddress    = "http://localhost:8500"
	defaultTimeout    = 10 * time.Second
	consulTokenHeader = "X-Consul-Token"
)

//go:generate go run github.com/gojuno/minimock/v3/cmd/minimock -g -i Client -s _mock.go

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

	// Token (optional) will be used to authenticate requests to consul.
	Token string

	// HTTPClient (optional) is the underlying HTTP client to use for making
	// requests to consul agents and servers. If not set, a default HTTP client
	// is used with a default timeout of 10 seconds, and will keep connections
	// open.
	HTTPClient *http.Client

	// Logger may be optionally configured as an output for trace level logging
	// produced internally by the Client. This can be helpful for debugging logic
	// errors in client code.
	Logger loggy.Logger
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

// New creates a new Client that will use the provided ClientOptions for
// making requests to a configured consul agent.
func New(opts ClientOptions) Client {
	address := opts.Address
	if address == "" {
		address = defaultAddress
	}

	httpClient := opts.HTTPClient
	if httpClient == nil {
		httpClient = clean.DefaultPooledClient()
		httpClient.Timeout = defaultTimeout
	}

	logger := opts.Logger
	if logger == nil {
		logger = loggy.Discard()
	}

	return &client{
		address:    address,
		token:      opts.Token,
		httpClient: httpClient,
		log:        logger,
	}
}

type client struct {
	address    string
	token      string
	httpClient *http.Client
	log        loggy.Logger
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

func (c *client) newRequest(ctx Ctx, method, fullURL string, body io.Reader) (*http.Request, error) {
	request, err := http.NewRequest(method, fullURL, body)
	if err != nil {
		return nil, err
	}
	rCtx := request.WithContext(ctx)
	c.maybeSetToken(rCtx)
	c.setHeaders(rCtx)
	return rCtx, nil
}

func (c *client) get(ctx Ctx, path string, i interface{}) error {
	completeURL := c.address + path

	request, err := c.newRequest(ctx, http.MethodGet, completeURL, nil)
	if err != nil {
		return err
	}

	response, err := c.httpClient.Do(request)
	if err != nil {
		return err
	}
	defer ignore.Drain(response.Body)

	if response.StatusCode >= 400 {
		return &RequestError{statusCode: response.StatusCode}
	}

	return json.NewDecoder(response.Body).Decode(i)
}

func (c *client) put(ctx Ctx, path, body string, i interface{}) error {
	completeURL := c.address + path

	r := strings.NewReader(body)
	request, err := c.newRequest(ctx, http.MethodPut, completeURL, r)
	if err != nil {
		return err
	}

	response, err := c.httpClient.Do(request)
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

func (c *client) delete(ctx Ctx, path string) error {
	completeURL := c.address + path

	request, err := c.newRequest(ctx, http.MethodDelete, completeURL, nil)
	if err != nil {
		return err
	}

	response, err := c.httpClient.Do(request)
	if err != nil {
		return err
	}
	defer ignore.Drain(response.Body)

	if response.StatusCode >= 400 {
		return &RequestError{statusCode: response.StatusCode}
	}

	return nil
}

func (c *client) maybeSetToken(request *http.Request) {
	if c.token != "" {
		request.Header.Set(consulTokenHeader, c.token)
	}
}

const (
	headerContentType = "Content-Type"
	headerUserAgent   = "User-Agent"
	mimeJSON          = "application/json"
	userAgent         = "consulapi/1.0"
)

func (c *client) setHeaders(request *http.Request) {
	request.Header.Set(headerContentType, mimeJSON)
	request.Header.Set(headerUserAgent, userAgent)
}
