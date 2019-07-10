package consulapi

import (
	"fmt"
	"net/url"
	"sort"
)

//go:generate go run github.com/gojuno/minimock/v3/cmd/minimock -g -i Catalog -s _mock.go

// A Catalog represents the consul catalog feature.
//
// Note that the register and deregister endpoints are not implemented. Per the
// consul documentation, it is preferable to make use of the Agent endpoints
// for service registrations.
type Catalog interface {

	// DataCenters returns the list of all known DCs. The order of the
	// list of DCs is in estimated round trip distance (i.e. the dcs
	// furthest away are listed last).
	//
	// https://www.consul.io/api/catalog.html#list-datacenters
	DataCenters(Ctx) ([]string, error)

	// Nodes will return the list of nodes in dc.
	//
	// https://www.consul.io/api/catalog.html#list-nodes
	Nodes(Ctx, NodesQuery) ([]Node, error)

	// Node will return detailed meta information associated
	// a particular node in dc.
	//
	// https://www.consul.io/api/catalog.html#list-services-for-node
	Node(Ctx, string, NodeQuery) (NodeInfo, error)

	// Services will return a list of names of services
	// in dc, along with the associated tags for each service.
	//
	// https://www.consul.io/api/catalog.html#list-services
	Services(Ctx, ServicesQuery) (map[string][]string, error)

	// Service returns detailed meta information about a particular
	// named service, in dc, which matches all of the listed tags.
	//
	// https://www.consul.io/api/catalog.html#list-nodes-for-service
	Service(Ctx, string, ServiceQuery) ([]Instance, error)

	// Connect returns the detailed meta information about a particular
	// consul CONNECT enabled service in a given DC.
	//
	// https://www.consul.io/api/catalog.html#list-nodes-for-connect-capable-service
	Connect(Ctx, string, ServiceQuery) ([]Instance, error)
}

func (c *client) DataCenters(ctx Ctx) ([]string, error) {
	dcs := make([]string, 0, 10)
	if err := c.get(ctx, "/v1/catalog/datacenters", &dcs); err != nil {
		return nil, err
	}
	return dcs, nil
}

func (p Pair) String() string {
	return fmt.Sprintf("%s:%s", p.Key, p.Value)
}

// NodesQuery is used to define values for each of the optional parameters
// to the catalog nodes endpoint.
//
// https://www.consul.io/api/catalog.html#parameters-2
type NodesQuery struct {
	// DC indicates the datacenter to query.
	//
	// If blank, this will default to the datacenter that the queried agent is in.
	DC string

	// Near specifies which node should be treated as the "center" in terms of
	// round-trip time for ordering the returned nodes. Nodes at the beginning
	// of the list will have the shortest round-trip times to the given node.
	//
	// The value "_agent" will cause the agent's node to be used for the sort.
	//
	// If blank, no default behavior is defined.
	Near string

	// NodeMeta creates a filter based on node metadata, using the given list
	// of key:value pairs.
	//
	// If blank, no node metadata filter is applied.
	NodeMeta []Pair

	// Filter specifies an advanced filtering expression.
	//
	// https://www.consul.io/api/features/filtering.html
	Filter string
}

// A Node represents a host on which a consul agent is running.
type Node struct {
	Name            string            `json:"Node"`
	Address         string            `json:"Address"`
	TaggedAddresses map[string]string `json:"TaggedAddresses"`
}

func (c *client) Nodes(ctx Ctx, nq NodesQuery) ([]Node, error) {
	var params [][2]string

	if nq.DC != "" {
		params = append(params, [2]string{"dc", nq.DC})
	}

	if nq.Near != "" {
		params = append(params, [2]string{"near", nq.Near})
	}

	for _, pair := range nq.NodeMeta {
		params = append(params, [2]string{"node-meta", pair.String()})
	}

	if nq.Filter != "" {
		params = append(params, [2]string{"filter", url.QueryEscape(nq.Filter)})
	}

	path := fixup("/v1/catalog", "/nodes", params...)
	nodes := make([]Node, 0, 100)

	if err := c.get(ctx, path, &nodes); err != nil {
		return nil, err
	}

	return nodes, nil
}

// A NodeInfo contains detailed information about a node,
// including all of the services defined to exist on that
// node.
type NodeInfo struct {
	Node     Node `json:"Node"`
	Services map[string]struct {
		ID      string   `json:"ID"`
		Service string   `json:"Service"`
		Tags    []string `json:"Tags"`
		Port    int      `json:"Port"`
	} `json:"Services"`
}

// A NodeQuery is used to define values for each of the optional parameters
// to the catalog node endpoint.
type NodeQuery struct {
	// DC indicates the datacenter to query.
	//
	// If blank, this will default to the datacenter that the queried agent is in.
	DC string

	// Filter specifies an advanced filtering expression.
	//
	// https://www.consul.io/api/features/filtering.html
	Filter string
}

func (c *client) Node(ctx Ctx, name string, nq NodeQuery) (NodeInfo, error) {
	var params [][2]string

	if nq.DC != "" {
		params = append(params, [2]string{"dc", nq.DC})
	}

	if nq.Filter != "" {
		params = append(params, [2]string{"filter", url.QueryEscape(nq.Filter)})
	}

	path := fixup("/v1/catalog", "/node/"+name, params...)

	var info NodeInfo
	if err := c.get(ctx, path, &info); err != nil {
		return NodeInfo{}, err
	}

	return info, nil
}

type ServicesQuery struct {
	// DC indicates the datacenter to query.
	//
	// If blank, this will default to the datacenter that the queried agent is in.
	DC string

	// NodeMeta creates a filter based on node metadata, using the given list
	// of key:value pairs.
	//
	// If blank, no node metadata filter is applied.
	NodeMeta []Pair
}

func (c *client) Services(ctx Ctx, sq ServicesQuery) (map[string][]string, error) {
	var params [][2]string

	if sq.DC != "" {
		params = append(params, [2]string{"dc", sq.DC})
	}

	for _, pair := range sq.NodeMeta {
		params = append(params, [2]string{"node-meta", pair.String()})
	}

	path := fixup("/v1/catalog", "/services", params...)

	services := make(map[string][]string, 1024)
	if err := c.get(ctx, path, &services); err != nil {
		return nil, err
	}

	// Sort the list of tags for each returned service, so that the response
	// is deterministic, fixing an annoyance with the raw reply.
	for _, tags := range services {
		sort.Strings(tags)
	}

	return services, nil
}

type ServiceQuery struct {
	// DC indicates the datacenter to query.
	//
	// If blank, this will default to the datacenter that the queried agent is in.
	DC string

	// Tags specifies a list of tags to filter on. Only instances matching all
	// given tags will be returned.
	Tags []string

	// Near specifies which node should be treated as the "center" in terms of
	// round-trip time for ordering the returned nodes. Nodes at the beginning
	// of the list will have the shortest round-trip times to the given node.
	//
	// The value "_agent" will cause the agent's node to be used for the sort.
	//
	// If blank, no default behavior is defined.
	Near string

	// NodeMeta creates a filter based on node metadata, using the given list
	// of key:value pairs.
	//
	// If blank, no node metadata filter is applied.
	NodeMeta []Pair

	// Filter specifies an advanced filtering expression.
	//
	// https://www.consul.io/api/features/filtering.html
	Filter string
}

func (c *client) Service(ctx Ctx, service string, sq ServiceQuery) ([]Instance, error) {
	serviceEP := "/v1/catalog/service/"
	return c.service(ctx, serviceEP, service, sq)
}

func (c *client) Connect(ctx Ctx, service string, sq ServiceQuery) ([]Instance, error) {
	connectEP := "/v1/catalog/connect/"
	return c.service(ctx, connectEP, service, sq)
}

func (c *client) service(ctx Ctx, ep, service string, sq ServiceQuery) ([]Instance, error) {
	var params [][2]string

	if sq.DC != "" {
		params = append(params, [2]string{"dc", sq.DC})
	}

	for _, tag := range sq.Tags {
		params = append(params, [2]string{"tag", tag})
	}

	if sq.Near != "" {
		params = append(params, [2]string{"near", sq.Near})
	}

	for _, pair := range sq.NodeMeta {
		params = append(params, [2]string{"node-meta", pair.String()})
	}

	if sq.Filter != "" {
		params = append(params, [2]string{"filter", url.QueryEscape(sq.Filter)})
	}

	path := fixup(ep, service, params...)
	instances := make([]Instance, 0, 100)

	if err := c.get(ctx, path, &instances); err != nil {
		return nil, err
	}

	return instances, nil
}
