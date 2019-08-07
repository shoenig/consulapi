package consulapi

type Pair struct {
	Key   string `json:"Key"`
	Value string `json:"Value"`
}

type Address struct {
	Address string `json:"address"`
	Port    int    `json:"port"`
}

type Proxy struct {
	DestinationServiceName string `json:"DestinationServiceName"`
	DestinationServiceID   string `json:"DestinationServiceID"`
	LocalServiceAddress    string `json:"LocalServiceAddress"`
	LocalServicePort       int    `json:"LocalServicePort"`
	// Config
	// Upstreams
}

type ServiceConnect struct {
	Native bool `json:"Native"`
	// Proxy
}

type Instance struct {
	ID                       string             `json:"ID"`
	Node                     string             `json:"Node"`
	Address                  string             `json:"Address"`
	Datacenter               string             `json:"Datacenter"`
	TaggedAddresses          map[string]string  `json:"TaggedAddresses"`
	NodeMeta                 map[string]string  `json:"NodeMeta"`
	ServiceAddress           string             `json:"ServiceAddress"`
	ServiceEnableTagOverride bool               `json:"ServiceEnableTagOverride"`
	ServiceID                string             `json:"ServiceID"`
	ServiceName              string             `json:"ServiceName"`
	ServicePort              int                `json:"ServicePort"`
	ServiceMeta              map[string]string  `json:"ServiceMeta"`
	ServiceTaggedAddresses   map[string]Address `json:"ServiceTaggedAddresses"`
	ServiceTags              []string           `json:"ServiceTags"`
	ServiceProxyDestination  string             `json:"ServiceProxyDestination"`
	ServiceProxy             Proxy              `json:"ServiceProxy"`
	ServiceConnect           ServiceConnect     `json:"ServiceConnect"`
}
