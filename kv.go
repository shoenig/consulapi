// Author hoenig

package consulapi

import (
	"encoding/base64"
	"sort"

	"github.com/pkg/errors"
)

// A KV represents the key-value store built into consul.
//
// Although consul supports arbitrary bytes as keys and values,
// this library assumes all keys and values are strings. This
// helps simplify code for clients, the 99% use case for which
// is reading and writing small configuration values and other
// string-y information.
//
// Each method allows for specifying a particular dc from which
// to set or retrieve information. If left unset, the dc defaults
// to the dc associated with the consul agent being communicated
// with.
type KV interface {
	// Get will return the value defined at path, for dc.
	Get(dc, path string) (string, error)
	// Put will set value at path, in dc.
	Put(dc, path, value string) error
	// Delete will remove the value at path, in dc.
	Delete(dc, path string) error
	// Keys will list all subpaths in asciibetical order.
	// The returned paths may be terminal (ie, the value is
	// stored content) or they may be further traversable like
	// a directory listing, in dc.
	Keys(dc, path string) ([]string, error)
	// Recurse will recursively descend through path, collecting
	// all KV pairs along the way, in dc.
	Recurse(dc, path string) ([][2]string, error)
}

type value struct {
	Key   string `json:"Key"`
	Value string `json:"Value"`
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
