package consulapi

import (
	"bytes"
	"context"
	"crypto/tls"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"gophers.dev/pkgs/loggy"
)

type responder struct {
	t *testing.T // our test controller

	code int    // respond with http status code
	body string // respond with this body

	hasMethod  string              // assert request has this HTTP method type
	hasPath    string              // assert request has this path
	hasQuery   map[string][]string // assert request as this query
	hasHeaders map[string]string   // assert request has these headers
	hasBody    string              // assert request has this body
}

func (rs *responder) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 0) check responder is configured all the way
	rs.isConfigured()

	// 1) check request method
	rs.checkMethod(r)

	// 2) check request path
	rs.checkPath(r)

	// 3) check request query
	rs.checkQuery(r)

	// 4) check request header
	rs.checkHeaders(r)

	// 5) check request body
	rs.checkBody(r)

	// 6) okay now we can write the response
	w.Header().Set(headerContentType, mimeJSON)
	w.WriteHeader(rs.code)
	_, _ = w.Write([]byte(rs.body))
}

func (rs *responder) isConfigured() {
	if rs.hasMethod == "" {
		rs.t.Fatal("responder: hasMethod not specified")
	}

	if rs.hasPath == "" {
		rs.t.Fatal("responder: hasPath not specified")
	}

	if rs.code == 0 {
		rs.t.Fatal("responder: code not set")
	}
}

func (rs *responder) checkMethod(r *http.Request) {
	rMethod := r.Method
	require.Equal(rs.t, rs.hasMethod, rMethod)
}

func (rs *responder) checkPath(r *http.Request) {
	rPath := r.URL.Path
	require.Equal(rs.t, rs.hasPath, rPath)
}

func (rs *responder) checkQuery(r *http.Request) {
	values := r.URL.Query()

	expN := len(rs.hasQuery)
	N := len(values)
	require.Equal(rs.t, expN, N, "expected %d query params, got %d", expN, N)

	for key, expValues := range rs.hasQuery {
		list, exists := values[key]
		require.True(rs.t, exists, "expected query %s:%s", key, expValues)
		require.Equal(rs.t, len(expValues), len(list))
		// for each expVal, make sure list contains it
		for _, expVal := range expValues {
			require.Contains(rs.t, list, expVal)
		}
	}
}

func (rs *responder) checkHeaders(r *http.Request) {
	if rs.hasHeaders == nil {
		rs.hasHeaders = make(map[string]string)
	}

	// insert check for Content-Type:application/json if not set
	if _, exists := rs.hasHeaders[headerContentType]; !exists {
		rs.hasHeaders[headerContentType] = mimeJSON
	}

	for key, v := range rs.hasHeaders {
		value := r.Header.Get(key)
		require.Equal(rs.t, v, value)
	}
}

func (rs *responder) checkBody(r *http.Request) {
	bs, err := ioutil.ReadAll(r.Body)
	require.NoError(rs.t, err)
	require.Equal(rs.t, rs.hasBody, string(bs))
	r.Body = ioutil.NopCloser(bytes.NewReader(bs))
}

func load(t *testing.T, file string) string {
	filePath := filepath.Join("hack/resources", file)
	bs, err := ioutil.ReadFile(filePath)
	require.NoError(t, err)
	s := strings.TrimSpace(string(bs))
	return s
}

func testClient(h http.Handler) (Ctx, *httptest.Server, Client) {
	ts := httptest.NewServer(h)

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	httpClient := &http.Client{
		Transport: transport,
		Timeout:   1 * time.Second,
	}

	client := New(ClientOptions{
		Address:    ts.URL,
		Token:      "abc132",
		HTTPClient: httpClient,
		Logger:     loggy.New("test-client"),
	})

	ctx := context.Background()

	return ctx, ts, client
}
