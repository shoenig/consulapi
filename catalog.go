// Author hoenig

package consulapi

import "sort"

// A Catalog represents the consul catalog feature.
type Catalog interface {
	// Datacenters will return the list of datacenters that
	// are members of the gossip ring known by the consul agent.
	Datacenters() ([]string, error)

	// Nodes will return the list of nodes in dc.
	Nodes(dc string) ([]Node, error)

	// Node will return detailed meta information associated
	// a particulare node in dc.
	Node(dc, name string) (NodeInfo, error)

	// Services will return a list of names of services
	// in dc, along with the associated tags for each service.
	Services(dc string) (map[string][]string, error)

	// Service returns detailed meta information about a particular
	// named service, in dc, which matches all of the listed tags.
	Service(dc, service string, tags ...string) ([]Service, error)
}

func (c *client) Datacenters() ([]string, error) {
	dcs := make([]string, 0, 10)
	if err := c.get("/v1/catalog/datacenters", &dcs); err != nil {
		return nil, err
	}
	return dcs, nil
}

// A Node represents a host on which a consul agent is running.
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

// A NodeInfo contains detailed information about a node,
// including all of the services defined to exist on that
// node.
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

// A Service defines a program that is configured to be
// running somewhere.
//
// By default, a service is defined on
// the same node the service itself is running. However, it
// is possible to specify the ServiceAddress, which 'overrides'
// the Address value (which is always associated with the node),
// which allows for pointing at services that may be running on
// nodes that are not associated with the consul cluster.
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
