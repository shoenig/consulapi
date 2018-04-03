// Author hoenig

package consulapi

import (
	"strconv"

	"github.com/pkg/errors"
)

//go:generate mockery -interface Agent -package consulapitest

// Agent provides an interface to information about the
// consul agent that is being communicated with.
type Agent interface {
	Self() (AgentInfo, error)
	Members(wan bool) ([]AgentInfo, error)
	Reload() error
	MaintenanceMode(enabled bool, reason string) error
}

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

func (c *client) Self() (AgentInfo, error) {
	path := fixup("/v1/agent/", "self")
	var response selfResponse
	if err := c.get(path, &response); err != nil {
		return AgentInfo{}, errors.Wrap(err, "failed to get agent self info")
	}

	return AgentInfo{
		Name:    response.Config.NodeName,
		Address: response.Member.Addr,
		Port:    response.Member.Port,
		Tags:    response.Member.Tags,
	}, nil
}

func (c *client) Members(wan bool) ([]AgentInfo, error) {
	path := fixup("/v1/agent", "/members")
	if wan {
		// we cannot simply do wan=false, as consul then appends the dc to the hostname (!)
		wanS := strconv.FormatBool(wan)
		path = fixup("/v1/agent", "/members", [2]string{"wan", wanS})
	}

	agentInfos := make([]AgentInfo, 0, 1000)
	err := c.get(path, &agentInfos)
	return agentInfos, err
}

func (c *client) Reload() error {
	path := fixup("/v1/agent", "/reload")
	return c.put(path, "", nil)
}

func (c *client) MaintenanceMode(enabled bool, reason string) error {
	enableS := strconv.FormatBool(enabled)
	path := fixup("/v1/agent", "/maintenance", [2]string{"enable", enableS}, [2]string{"reason", reason})
	return c.put(path, "", nil)
}
