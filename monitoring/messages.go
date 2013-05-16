package monitoring

import (
	"encoding/json"
	"github.com/racker/go-agent-endpoint/endpoint"
)

type Check struct {
	Id       string            `json:"id"`
	Type     string            `json:"type"`
	Details  map[string]string `json:"details"`
	Period   int               `json:"period"`
	Timeout  int               `json:"timeout"`
	Disabled bool              `json:"disabled"`
}

type CheckScheduleGetResult struct {
	Checks []Check `json:"checks"`
}

type Metric struct {
	Type  string `json:"t"`
	Value string `json:"v"`
	Unit  string `json:"u"`
}

type MetricGroup struct {
	Prefix  string
	Metrics map[string]Metric
}

type CheckMetricsPostParams struct {
	CheckId   string        `json:"check_id"`
	CheckType string        `json:"check_type"`
	State     string        `json:"state"`
	Status    string        `json:"status"`
	Metrics   []MetricGroup `json:"metrics"`
}

func (g MetricGroup) MarshalJSON() (js []byte, err error) {
	js, err = json.Marshal([2]interface{}{g.Prefix, g.Metrics})
	return
}

func (g *MetricGroup) UnmarshalJSON(data []byte) error {
	if g == nil {
		return endpoint.NilPointerError
	}
	var tmp []json.RawMessage
	err := json.Unmarshal(data, &tmp)
	if err != nil {
		return err
	}
	err = json.Unmarshal(tmp[0], &g.Prefix)
	if err != nil {
		return endpoint.UnmarshalFieldError
	}
	err = json.Unmarshal(tmp[1], &g.Metrics)
	if err != nil {
		return endpoint.UnmarshalFieldError
	}
	return nil
}
