package consulapi

import (
	"encoding/base64"
	"net/http"
	"sort"

	"github.com/pkg/errors"
)

type Query struct {
	DC string
}

//go:generate go run github.com/gojuno/minimock/v3/cmd/minimock -g -i KV -s _mock.go

// A KV can access the key-value store of consul.
//
// The consul KV store is useful for storing small values of information. It is
// intended to be used to store things like service configuration and other
// meta-data. The maximum size of a value is 512KiB.
//
// Each DC contains its own KV store. The data in a KV store is not replicated
// across DCs.
//
// https://www.consul.io/api/kv.html
type KV interface {

	// Get will return the value defined at path, for dc.
	Get(Ctx, string, Query) (string, error)

	// Put will set value at path, in dc.
	Put(Ctx, string, string, Query) error

	// Delete will remove the value at path, in dc.
	Delete(Ctx, string, Query) error

	// Keys will list all subpaths in asciibetical order.
	// The returned paths may be terminal (ie, the value is
	// stored content) or they may be further traversable like
	// a directory listing, in dc.
	Keys(Ctx, string, Query) ([]string, error)

	// Recurse will recursively descend through path, collecting
	// all KV pairs along the way, in dc.
	Recurse(Ctx, string, Query) ([]Pair, error)
}

func (c *client) Get(ctx Ctx, path string, query Query) (string, error) {
	var params [][2]string

	if query.DC != "" {
		params = append(params, [2]string{"dc", query.DC})
	}

	path = fixup("/v1/kv", path, params...)

	var values []Pair

	if err := c.get(ctx, path, &values); err != nil {
		if re, ok := err.(*RequestError); ok {
			if re.StatusCode() == http.StatusNotFound {
				return "", errors.Errorf("key %q does not exist", path)
			}
		}
		return "", err
	}

	bs, err := base64.StdEncoding.DecodeString(values[0].Value)
	if err != nil {
		return "", err
	}

	return string(bs), nil
}

func (c *client) Put(ctx Ctx, path, value string, query Query) error {
	var params [][2]string

	if query.DC != "" {
		params = append(params, [2]string{"dc", query.DC})
	}

	path = fixup("/v1/kv", path, params...)

	if err := c.put(ctx, path, value, nil); err != nil {
		return err
	}

	return nil
}

func (c *client) Delete(ctx Ctx, path string, query Query) error {
	var params [][2]string

	if query.DC != "" {
		params = append(params, [2]string{"dc", query.DC})
	}

	path = fixup("/v1/kv", path, params...)

	if err := c.delete(ctx, path); err != nil {
		return err
	}

	return nil
}

func (c *client) Keys(ctx Ctx, path string, query Query) ([]string, error) {
	var params [][2]string

	if query.DC != "" {
		params = append(params, [2]string{"dc", query.DC})
	}

	params = append(params, [2]string{"keys", "true"})

	path = fixup("/v1/kv", path, params...)

	var keys []string
	if err := c.get(ctx, path, &keys); err != nil {
		return nil, err
	}

	sort.Strings(keys)
	return keys, nil
}

func (c *client) Recurse(ctx Ctx, path string, query Query) ([]Pair, error) {
	var params [][2]string

	if query.DC != "" {
		params = append(params, [2]string{"dc", query.DC})
	}

	params = append(params, [2]string{"recurse", "true"})

	rPath := fixup("/v1/kv", path, params...)

	var values []Pair

	if err := c.get(ctx, rPath, &values); err != nil {
		if re, ok := err.(*RequestError); ok {
			if re.StatusCode() == http.StatusNotFound {
				return nil, errors.Errorf("key-space %q does not exist", path)
			}
		}
		return nil, err
	}

	kvPairs := make([]Pair, 0, len(values))

	for _, value := range values {
		decoded, err := base64.StdEncoding.DecodeString(value.Value)
		if err != nil {
			return nil, err
		}

		kvPairs = append(kvPairs, Pair{
			Key:   value.Key,
			Value: string(decoded),
		})
	}

	sort.Slice(kvPairs, func(i, j int) bool {
		return kvPairs[i].Key < kvPairs[j].Key
	})

	return kvPairs, nil
}
