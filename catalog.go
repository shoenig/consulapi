// Author hoenig

package consulapi

type Catalog interface {
	Datacenters() ([]string, error)
	Nodes(dc string) ([]Node, error)
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
