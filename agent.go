package consulapi

import (
	"encoding/json"
	"strconv"

	"github.com/pkg/errors"
)

//go:generate go run github.com/gojuno/minimock/v3/cmd/minimock -g -i Agent -s _mock.go

// Agent provides an interface to information about the consul agent that is
// being communicated with.
type Agent interface {

	// Self reports internal metrics specific to this running instance of the
	// consul service.
	//
	// https://www.consul.io/api/agent.html#read-configuration
	Self(ctx Ctx) (AgentInfo, error)

	// Members reports what instances of the consul service belong to the
	// consul cluster. If wan is true, the
	//
	// https://www.consul.io/api/agent.html#list-members
	Members(ctx Ctx, wan bool) ([]AgentInfo, error)

	// Reload the consul service config files.
	//
	// https://www.consul.io/api/agent.html#reload-agent
	Reload(ctx Ctx) error

	// MaintenanceMode puts a consul instance into a mode where the node is
	// marked unavailable, and will not be present in DNS or API queries.
	//
	// https://www.consul.io/api/agent.html#enable-maintenance-mode
	MaintenanceMode(ctx Ctx, enabled bool, reason string) error

	// Metrics will report about the consul instance.
	//
	// https://www.consul.io/api/agent.html#view-metrics
	Metrics(ctx Ctx) (Metrics, error)

	// Join the consul instance with an existing consul cluster.
	//
	// https://www.consul.io/api/agent.html#join-agent
	Join(ctx Ctx, address string, wan bool) error

	// Leave the consul cluster.
	//
	// https://www.consul.io/api/agent.html#graceful-leave-and-shutdown
	Leave(ctx Ctx) error

	// ForceLeave will purge the named node from the consul cluster.
	//
	// https://www.consul.io/api/agent.html#force-leave-and-shutdown
	ForceLeave(ctx Ctx, node string) error

	// SetACLToken will set the given kind of token to the value.
	//
	// https://www.consul.io/api/agent.html#update-acl-tokens
	SetACLToken(ctx Ctx, kind, token string) error

	// Monitor(loglevel string) // log stream, maybe someday
}

// An assertions that client satisfies Agent
var _ Agent = (*client)(nil)

// An AgentInfo contains information about a particular
// consul agent.
type AgentInfo struct {
	Name    string            `json:"Name"`
	Address string            `json:"Addr"`
	Port    int               `json:"Port"`
	Tags    map[string]string `json:"Tags"`
}

type selfResponse struct {
	Config struct {
		Datacenter string `json:"Datacenter"`
		NodeName   string `json:"NodeName"`
		Server     bool   `json:"Server"`
		Version    string `json:"Version"`
	} `json:"Config"`
	Member struct {
		Addr string            `json:"Addr"`
		Port int               `json:"Port"`
		Tags map[string]string `json:"Tags"`
	} `json:"Member"`
}

func (c *client) Self(ctx Ctx) (AgentInfo, error) {
	rPath := fixup("/v1/agent/", "self")

	var response selfResponse

	if err := c.get(ctx, rPath, &response); err != nil {
		return AgentInfo{}, err
	}

	return AgentInfo{
		Name:    response.Config.NodeName,
		Address: response.Member.Addr,
		Port:    response.Member.Port,
		Tags:    response.Member.Tags,
	}, nil
}

func (c *client) Members(ctx Ctx, wan bool) ([]AgentInfo, error) {
	rPath := fixup("/v1/agent", "/members")
	if wan {
		// we cannot simply do wan=false,
		// as consul then appends the dc to the hostname (!)
		wanS := strconv.FormatBool(wan)
		rPath = fixup("/v1/agent", "/members", [2]string{"wan", wanS})
	}

	var agentInfos []AgentInfo

	if err := c.get(ctx, rPath, &agentInfos); err != nil {
		return nil, err
	}

	return agentInfos, nil
}

func (c *client) Reload(ctx Ctx) error {
	rPath := fixup("/v1/agent", "/reload")
	if err := c.put(ctx, rPath, "", nil); err != nil {
		return err
	}
	return nil
}

func (c *client) MaintenanceMode(ctx Ctx, enabled bool, reason string) error {
	enableS := strconv.FormatBool(enabled)
	rPath := fixup("/v1/agent", "/maintenance", [2]string{"enable", enableS}, [2]string{"reason", reason})
	if err := c.put(ctx, rPath, "", nil); err != nil {
		return err
	}
	return nil
}

func (c *client) Metrics(ctx Ctx) (Metrics, error) {
	rPath := fixup("/v1/agent", "/metrics")

	var metrics Metrics
	if err := c.get(ctx, rPath, &metrics); err != nil {
		return Metrics{}, err
	}

	return metrics, nil
}

func (c *client) Join(ctx Ctx, address string, wan bool) error {
	rPath := fixup(
		"/v1/agent/join", address,
		[2]string{"wan", strconv.FormatBool(wan)},
	)

	if err := c.put(ctx, rPath, "", nil); err != nil {
		return err
	}

	return nil
}

func (c *client) Leave(ctx Ctx) error {
	rPath := fixup("/v1/agent", "/leave")
	if err := c.put(ctx, rPath, "", nil); err != nil {
		return err
	}
	return nil
}

func (c *client) ForceLeave(ctx Ctx, node string) error {
	rPath := fixup("/v1/agent/force-leave", node)
	if err := c.put(ctx, rPath, "", nil); err != nil {
		return err
	}
	return nil
}

type setToken struct {
	Token string `json:"Token"`
}

func (c *client) SetACLToken(ctx Ctx, kind, token string) error {

	switch kind {
	case "default", "agent", "agent_master", "replication":
	default:
		return errors.Errorf("unrecognized kind of token %q", kind)
	}

	rPath := fixup("/v1/agent/token", kind)

	bs, err := json.Marshal(setToken{Token: token})
	if err != nil {
		return errors.Wrap(err, "unable to create token payload")
	}

	if err := c.put(ctx, rPath, string(bs), nil); err != nil {
		return err
	}

	return nil
}
