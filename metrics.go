package consulapi

type Gauge struct {
	Name   string            `json:"Name"`
	Value  int               `json:"Value"`
	Labels map[string]string `json:"Labels"`
}

type Counter struct {
	Name   string            `json:"Name"`
	Count  int               `json:"Count"`
	Rate   float64           `json:"Rate"`
	Sum    int               `json:"Sum"`
	Min    int               `json:"Min"`
	Max    int               `json:"Max"`
	Mean   float64           `json:"Mean"`
	Stddev float64           `json:"Stddev"`
	Labels map[string]string `json:"Labels"`
}

type Sample struct {
	Name   string            `json:"Name"`
	Count  int               `json:"Count"`
	Rate   float64           `json:"Rate"`
	Sum    float64           `json:"Sum"`
	Min    float64           `json:"Min"`
	Max    float64           `json:"Max"`
	Mean   float64           `json:"Mean"`
	Stddev float64           `json:"Stddev"`
	Labels map[string]string `json:"Labels"`
}

// Metrics contains information about the agent from the /v1/agent/metrics
// endpoint.
type Metrics struct {
	Timestamp string    `json:"Timestamp"`
	Gauges    []Gauge   `json:"Gauges"`
	Counters  []Counter `json:"Counters"`
	Samples   []Sample  `json:"Samples"`
}
