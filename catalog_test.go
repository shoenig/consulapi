package consulapi

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Client_v1_catalog_datacenters_ok(t *testing.T) {
	ctx, ts, client := testClient(&responder{
		t:         t,
		code:      http.StatusOK,
		body:      load(t, "v1_catalog_datacenters.json"),
		hasPath:   "/v1/catalog/datacenters",
		hasMethod: http.MethodGet,
	})
	defer ts.Close()

	dcs, err := client.DataCenters(ctx)
	require.NoError(t, err)

	require.Equal(t, []string{"mydc1", "mydc2", "mydc3"}, dcs)
}

func Test_Client_v1_catalog_datacenters_err(t *testing.T) {
	ctx, ts, client := testClient(&responder{
		t:         t,
		code:      http.StatusInternalServerError,
		body:      "malfunction",
		hasPath:   "/v1/catalog/datacenters",
		hasMethod: http.MethodGet,
	})
	defer ts.Close()

	_, err := client.DataCenters(ctx)
	require.EqualError(t, err, "status code (500)")
}

func node(name, address, wan string) Node {
	return Node{
		Name:    name,
		Address: address,
		TaggedAddresses: map[string]string{
			"lan": address,
			"wan": wan,
		},
	}
}

func Test_Client_v1_catalog_nodes_defaults(t *testing.T) {
	ctx, ts, client := testClient(&responder{
		t:         t,
		code:      http.StatusOK,
		body:      load(t, "v1_catalog_nodes.json"),
		hasPath:   "/v1/catalog/nodes",
		hasMethod: http.MethodGet,
		hasQuery:  map[string][]string{},
	})
	defer ts.Close()

	nodes, err := client.Nodes(ctx, NodesQuery{
		// empty
	})
	require.NoError(t, err)

	require.Equal(t, []Node{
		node("dc1-node1", "10.0.0.1", "1.1.1.1"),
		node("dc1-node2", "10.0.0.2", "1.1.1.2"),
		node("dc2-node1", "10.0.1.1", "1.1.1.3"),
	}, nodes)
}

func Test_Client_v1_catalog_nodes_defaults_err(t *testing.T) {
	ctx, ts, client := testClient(&responder{
		t:         t,
		code:      http.StatusInternalServerError,
		body:      "malfunction",
		hasPath:   "/v1/catalog/nodes",
		hasMethod: http.MethodGet,
		hasQuery:  map[string][]string{},
	})
	defer ts.Close()

	_, err := client.Nodes(ctx, NodesQuery{
		// empty
	})
	require.EqualError(t, err, "status code (500)")
}

func Test_Client_v1_catalog_nodes_dc(t *testing.T) {
	ctx, ts, client := testClient(&responder{
		t:         t,
		code:      http.StatusOK,
		body:      load(t, "v1_catalog_nodes-dc.json"),
		hasPath:   "/v1/catalog/nodes",
		hasMethod: http.MethodGet,
		hasQuery: map[string][]string{
			"dc": {"dc2"},
		},
	})
	defer ts.Close()

	nodes, err := client.Nodes(ctx, NodesQuery{
		DC: "dc2",
	})
	require.NoError(t, err)
	require.Equal(t, []Node{
		node("dc2-node1", "10.0.1.1", "1.1.1.3"),
	}, nodes)
}

func Test_Client_v1_catalog_nodes_near(t *testing.T) {
	ctx, ts, client := testClient(&responder{
		t:         t,
		code:      http.StatusOK,
		body:      load(t, "v1_catalog_nodes-near.json"),
		hasPath:   "/v1/catalog/nodes",
		hasMethod: http.MethodGet,
		hasQuery: map[string][]string{
			"near": {"dc1-node1"},
		},
	})
	defer ts.Close()

	nodes, err := client.Nodes(ctx, NodesQuery{
		Near: "dc1-node1",
	})
	require.NoError(t, err)
	require.Equal(t, []Node{
		node("dc1-node1", "10.0.0.1", "1.1.1.1"),
		node("dc1-node2", "10.0.0.2", "1.1.1.2"),
		node("dc2-node1", "10.0.1.1", "1.1.1.3"),
	}, nodes)
}

func Test_Client_v1_catalog_nodes_nodeMeta(t *testing.T) {
	ctx, ts, client := testClient(&responder{
		t:         t,
		code:      http.StatusOK,
		body:      load(t, "v1_catalog_nodes-nodemeta.json"),
		hasPath:   "/v1/catalog/nodes",
		hasMethod: http.MethodGet,
		hasQuery: map[string][]string{
			"node-meta": {
				"instance_type:t2.medium",
				"instance_type:t2.tiny",
			},
		},
	})
	defer ts.Close()

	nodes, err := client.Nodes(ctx, NodesQuery{
		NodeMeta: []Pair{
			{Key: "instance_type", Value: "t2.medium"},
			{Key: "instance_type", Value: "t2.tiny"},
		},
	})
	require.NoError(t, err)
	require.Equal(t, []Node{
		node("dc1-node1", "10.0.0.1", "1.1.1.1"),
	}, nodes)
}

func Test_Client_v1_catalog_nodes_filter(t *testing.T) {
	ctx, ts, client := testClient(&responder{
		t:         t,
		code:      http.StatusOK,
		body:      load(t, "v1_catalog_nodes-filter.json"),
		hasPath:   "/v1/catalog/nodes",
		hasMethod: http.MethodGet,
		hasQuery: map[string][]string{
			"filter": {url.QueryEscape("Meta.env == qa")},
		},
	})
	defer ts.Close()

	nodes, err := client.Nodes(ctx, NodesQuery{
		Filter: "Meta.env == qa",
	})
	require.NoError(t, err)
	require.Equal(t, []Node{
		node("dc1-node1", "10.0.0.1", "1.1.1.1"),
	}, nodes)
}

func Test_Client_v1_catalog_nodes_mix(t *testing.T) {
	ctx, ts, client := testClient(&responder{
		t:         t,
		code:      http.StatusOK,
		body:      load(t, "v1_catalog_nodes-filter.json"),
		hasPath:   "/v1/catalog/nodes",
		hasMethod: http.MethodGet,
		hasQuery: map[string][]string{
			"dc":        {"dc1"},
			"near":      {"dc1-node1"},
			"node-meta": {"instance_type:t2.tiny"},
			"filter":    {url.QueryEscape("Meta.env == qa")},
		},
	})
	defer ts.Close()

	nodes, err := client.Nodes(ctx, NodesQuery{
		DC:       "dc1",
		Near:     "dc1-node1",
		NodeMeta: []Pair{{Key: "instance_type", Value: "t2.tiny"}},
		Filter:   "Meta.env == qa",
	})
	require.NoError(t, err)
	require.Equal(t, []Node{
		node("dc1-node1", "10.0.0.1", "1.1.1.1"),
	}, nodes)
}

func Test_Client_v1_catalog_node(t *testing.T) {
	ctx, ts, client := testClient(&responder{
		t:         t,
		code:      http.StatusOK,
		body:      load(t, "v1_catalog_node.json"),
		hasPath:   "/v1/catalog/node/foobar",
		hasMethod: http.MethodGet,
		hasQuery:  map[string][]string{},
	})
	defer ts.Close()

	nodeInfo, err := client.Node(ctx, "foobar", NodeQuery{})
	require.NoError(t, err)
	require.Equal(t, "foobar", nodeInfo.Node.Name)
}

func Test_Client_v1_catalog_node_err(t *testing.T) {
	ctx, ts, client := testClient(&responder{
		t:         t,
		code:      http.StatusInternalServerError,
		body:      load(t, "v1_catalog_node.json"),
		hasPath:   "/v1/catalog/node/foobar",
		hasMethod: http.MethodGet,
		hasQuery:  map[string][]string{},
	})
	defer ts.Close()

	_, err := client.Node(ctx, "foobar", NodeQuery{})
	require.EqualError(t, err, "status code (500)")
}

func Test_Client_v1_catalog_node_dc(t *testing.T) {
	ctx, ts, client := testClient(&responder{
		t:         t,
		code:      http.StatusOK,
		body:      load(t, "v1_catalog_node-dc.json"),
		hasPath:   "/v1/catalog/node/foobar",
		hasMethod: http.MethodGet,
		hasQuery: map[string][]string{
			"dc": {"dc1"},
		},
	})
	defer ts.Close()

	nodeInfo, err := client.Node(ctx, "foobar", NodeQuery{
		DC: "dc1",
	})
	require.NoError(t, err)
	require.Equal(t, "foobar", nodeInfo.Node.Name)
}

func Test_Client_v1_catalog_node_filter(t *testing.T) {
	ctx, ts, client := testClient(&responder{
		t:         t,
		code:      http.StatusOK,
		body:      load(t, "v1_catalog_node-filter.json"),
		hasPath:   "/v1/catalog/node/foobar",
		hasMethod: http.MethodGet,
		hasQuery: map[string][]string{
			"filter": {url.QueryEscape("Meta.redis_version == 4.0")},
		},
	})
	defer ts.Close()

	nodeInfo, err := client.Node(ctx, "foobar", NodeQuery{
		Filter: "Meta.redis_version == 4.0",
	})
	require.NoError(t, err)
	require.Equal(t, "foobar", nodeInfo.Node.Name)
}

func Test_Client_v1_catalog_services(t *testing.T) {
	ctx, ts, client := testClient(&responder{
		t:         t,
		code:      http.StatusOK,
		body:      load(t, "v1_catalog_services.json"),
		hasPath:   "/v1/catalog/services",
		hasMethod: http.MethodGet,
		hasQuery:  map[string][]string{},
	})
	defer ts.Close()

	services, err := client.Services(ctx, ServicesQuery{
		// empty
	})
	require.NoError(t, err)
	require.Equal(t, 4, len(services))
	require.Equal(t, []string{"active", "standby"}, services["vault"])
}

func Test_Client_v1_catalog_services_err(t *testing.T) {
	ctx, ts, client := testClient(&responder{
		t:         t,
		code:      http.StatusInternalServerError,
		body:      "malfunction",
		hasPath:   "/v1/catalog/services",
		hasMethod: http.MethodGet,
		hasQuery:  map[string][]string{},
	})
	defer ts.Close()

	_, err := client.Services(ctx, ServicesQuery{
		// empty
	})
	require.EqualError(t, err, "status code (500)")
}

func Test_Client_v1_catalog_services_dc(t *testing.T) {
	ctx, ts, client := testClient(&responder{
		t:         t,
		code:      http.StatusOK,
		body:      load(t, "v1_catalog_services-dc.json"),
		hasPath:   "/v1/catalog/services",
		hasMethod: http.MethodGet,
		hasQuery: map[string][]string{
			"dc": {"dc1"},
		},
	})
	defer ts.Close()

	services, err := client.Services(ctx, ServicesQuery{
		DC: "dc1",
	})
	require.NoError(t, err)
	require.Equal(t, 1, len(services))
}

func Test_Client_v1_catalog_services_nodeMeta(t *testing.T) {
	ctx, ts, client := testClient(&responder{
		t:         t,
		code:      http.StatusOK,
		body:      load(t, "v1_catalog_services-nodemeta.json"),
		hasPath:   "/v1/catalog/services",
		hasMethod: http.MethodGet,
		hasQuery: map[string][]string{
			"node-meta": {"a:1", "b:2"},
		},
	})
	defer ts.Close()

	services, err := client.Services(ctx, ServicesQuery{
		NodeMeta: []Pair{
			{Key: "a", Value: "1"},
			{Key: "b", Value: "2"},
		},
	})

	require.NoError(t, err)
	require.Equal(t, 1, len(services))
}

func Test_Client_v1_catalog_service(t *testing.T) {
	ctx, ts, client := testClient(&responder{
		t:         t,
		code:      http.StatusOK,
		body:      load(t, "v1_catalog_service.json"),
		hasPath:   "/v1/catalog/service/myapp",
		hasMethod: http.MethodGet,
		hasQuery:  map[string][]string{},
	})
	defer ts.Close()

	instances, err := client.Service(ctx, "myapp", ServiceQuery{
		// empty
	})
	require.NoError(t, err)
	require.Equal(t, 2, len(instances))
	require.Equal(t, "myapp", instances[0].ServiceName)
	require.Equal(t, "myapp", instances[1].ServiceName)
}

func Test_Client_v1_catalog_service_err(t *testing.T) {
	ctx, ts, client := testClient(&responder{
		t:         t,
		code:      http.StatusInternalServerError,
		body:      "malfunction",
		hasPath:   "/v1/catalog/service/myapp",
		hasMethod: http.MethodGet,
		hasQuery:  map[string][]string{},
	})
	defer ts.Close()

	_, err := client.Service(ctx, "myapp", ServiceQuery{
		// empty
	})
	require.EqualError(t, err, "status code (500)")
}

func Test_Client_v1_catalog_service_mix(t *testing.T) {
	ctx, ts, client := testClient(&responder{
		t:         t,
		code:      http.StatusOK,
		body:      load(t, "v1_catalog_service.json"),
		hasPath:   "/v1/catalog/service/myapp",
		hasMethod: http.MethodGet,
		hasQuery: map[string][]string{
			"dc":        {"dc1"},
			"tag":       {"tag1", "tag2"},
			"near":      {"dc1-node7"},
			"node-meta": {"k1:v1"},
			"filter":    {url.QueryEscape("Meta.env == qa")},
		},
	})
	defer ts.Close()

	instances, err := client.Service(ctx, "myapp", ServiceQuery{
		DC:   "dc1",
		Tags: []string{"tag1", "tag2"},
		Near: "dc1-node7",
		NodeMeta: []Pair{{
			Key: "k1", Value: "v1",
		}},
		Filter: "Meta.env == qa",
	})
	require.NoError(t, err)
	require.Equal(t, 2, len(instances))
}

func Test_Client_v1_catalog_connect(t *testing.T) {
	ctx, ts, client := testClient(&responder{
		t:         t,
		code:      http.StatusOK,
		body:      load(t, "v1_catalog_service.json"),
		hasPath:   "/v1/catalog/connect/myapp",
		hasMethod: http.MethodGet,
		hasQuery:  map[string][]string{},
	})
	defer ts.Close()

	instances, err := client.Connect(ctx, "myapp", ServiceQuery{
		// empty
	})
	require.NoError(t, err)
	require.Equal(t, 2, len(instances))
	require.Equal(t, "myapp", instances[0].ServiceName)
	require.Equal(t, "myapp", instances[1].ServiceName)
}

func Test_Client_v1_catalog_connect_err(t *testing.T) {
	ctx, ts, client := testClient(&responder{
		t:         t,
		code:      http.StatusInternalServerError,
		body:      "malfunction",
		hasPath:   "/v1/catalog/connect/myapp",
		hasMethod: http.MethodGet,
		hasQuery:  map[string][]string{},
	})
	defer ts.Close()

	_, err := client.Connect(ctx, "myapp", ServiceQuery{
		// empty
	})
	require.EqualError(t, err, "status code (500)")
}

func Test_Client_v1_catalog_connect_mix(t *testing.T) {
	ctx, ts, client := testClient(&responder{
		t:         t,
		code:      http.StatusOK,
		body:      load(t, "v1_catalog_service.json"),
		hasPath:   "/v1/catalog/connect/myapp",
		hasMethod: http.MethodGet,
		hasQuery: map[string][]string{
			"dc":        {"dc1"},
			"tag":       {"tag1", "tag2"},
			"near":      {"dc1-node7"},
			"node-meta": {"k1:v1"},
			"filter":    {url.QueryEscape("Meta.env == qa")},
		},
	})
	defer ts.Close()

	instances, err := client.Connect(ctx, "myapp", ServiceQuery{
		DC:   "dc1",
		Tags: []string{"tag1", "tag2"},
		Near: "dc1-node7",
		NodeMeta: []Pair{{
			Key: "k1", Value: "v1",
		}},
		Filter: "Meta.env == qa",
	})
	require.NoError(t, err)
	require.Equal(t, 2, len(instances))
}
