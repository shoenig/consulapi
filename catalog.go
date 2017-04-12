// Author hoenig

package consulapi

import "sort"

type Catalog interface {
	Datacenters() ([]string, error)
	Nodes(dc string) ([]Node, error)
	Node(dc, name string) (NodeInfo, error)
	Services(dc string) (map[string][]string, error)
	Service(dc, service string, tags ...string) ([]Service, error)
}

func (c *client) Datacenters() ([]string, error) {
	dcs := make([]string, 0, 10)
	if err := c.get("/v1/catalog/datacenters", &dcs); err != nil {
		return nil, err
	}
	return dcs, nil
}

type Node struct {
	Name            string            `json:"Node"`
	Address         string            `json:"Address"`
	TaggedAddresses map[string]string `json:"TaggedAddresses"`
}

func (c *client) Nodes(dc string) ([]Node, error) {
	path := fixup("/v1/catalog", "/nodes", [2]string{"dc", dc})
	nodes := make([]Node, 0, 100)
	err := c.get(path, &nodes)
	return nodes, err
}

type NodeInfo struct {
	Node struct {
		Name            string            `json:"Node"`
		Address         string            `json:"Address"`
		TaggedAddresses map[string]string `json:"TaggedAddresses"`
	} `json:"Node"`
	Services map[string]struct {
		ID      string   `json:"ID"`
		Service string   `json:"Service"`
		Tags    []string `json:"Tags"`
		Port    int      `json:"Port"`
	} `json:"Services"`
}

func (c *client) Node(dc, name string) (NodeInfo, error) {
	path := fixup("/v1/catalog", "/node/"+name)
	var info NodeInfo
	err := c.get(path, &info)
	return info, err
}

func (c *client) Services(dc string) (map[string][]string, error) {
	path := fixup("/v1/catalog", "/services", [2]string{"dc", dc})
	services := make(map[string][]string, 1024)
	err := c.get(path, &services)
	// sort all the tags because i said so
	for _, tags := range services {
		sort.Strings(tags)
	}
	return services, err
}

type Service struct {
	Node            string            `json:"Node"`
	Address         string            `json:"Address"`
	TaggedAddresses map[string]string `json:"TaggedAddresses"`
	ServiceID       string            `json:"ServiceID"`
	ServiceName     string            `json:"ServiceName"`
	ServiceTags     []string          `json:"ServiceTags"`
	ServiceAddress  string            `json:"ServiceAddress"`
	ServicePort     int               `json:"ServicePort"`
}

func (c *client) Service(dc, service string, tags ...string) ([]Service, error) {
	params := [][2]string{{"dc", dc}}
	for _, tag := range tags {
		params = append(params, [2]string{"tag", tag})
	}
	path := fixup("/v1/catalog", "/service/"+service, params...)
	listing := make([]Service, 0, 100)
	err := c.get(path, &listing)
	return listing, err
}
