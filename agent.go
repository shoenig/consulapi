// Author hoenig

package consulapi

import "strconv"

type Agent interface {
	Members(wan bool) ([]AgentInfo, error)
	Reload() error
	MaintenanceMode(enabled bool, reason string) error
}

type AgentInfo struct {
	Name    string            `json:"Name"`
	Address string            `json:"Addr"`
	Port    int               `json:"Port"`
	Tags    map[string]string `json:"Tags"`
}

func (c *client) Members(wan bool) ([]AgentInfo, error) {
	path := fixup("/v1/agent", "/members")
	if wan {
		// we cannot simply do wan=false, as consul then appends the dc to the hostname (!)
		wanS := strconv.FormatBool(wan)
		path = fixup("/v1/agent", "/members", [2]string{"wan", wanS})
	}

	infos := make([]AgentInfo, 0, 1000)
	err := c.get(path, &infos)
	return infos, err
}

func (c *client) Reload() error {
	path := fixup("/v1/agent", "/reload")
	return c.put(path, "")
}

func (c *client) MaintenanceMode(enabled bool, reason string) error {
	enableS := strconv.FormatBool(enabled)
	path := fixup("/v1/agent", "/maintenance", [2]string{"enable", enableS}, [2]string{"reason", reason})
	return c.put(path, "")
}
