// Author hoenig

package consulapi

import (
	"encoding/base64"
	"sort"

	"github.com/pkg/errors"
)

type KV interface {
	Get(dc, path string) (string, error)
	Put(dc, path, value string) error
	Delete(dc, path string) error
	Keys(dc, path string) ([]string, error)
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
