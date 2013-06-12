package main

import (
	"encoding/json"
	"github.com/racker/go-agent-endpoint/endpoint"
)

type check struct {
	ID       string      `json:"id"`
	Type     string      `json:"type"`
	Details  interface{} `json:"details"`
	Period   int         `json:"period"`
	Timeout  int         `json:"timeout"`
	Disabled bool        `json:"disabled"`
}

type checkScheduleGetResult struct {
	Checks []check `json:"checks"`
}

type metric struct {
	Type  string `json:"t"`
	Value string `json:"v"`
	Unit  string `json:"u"`
}

type metricGroup struct {
	Prefix  string
	Metrics map[string]*metric
}

type checkMetricsPostParams struct {
	CheckID   string        `json:"check_id"`
	CheckType string        `json:"check_type"`
	State     string        `json:"state"`
	Status    string        `json:"status"`
	Metrics   []metricGroup `json:"metrics"`
}

// MarshalJSON marshals a metricGroup (g) into valid json bytes
func (g metricGroup) MarshalJSON() (js []byte, err error) {
	js, err = json.Marshal([2]interface{}{g.Prefix, g.Metrics})
	return
}

// UnmarshalJSON unmarshals valid json bytes into g.
func (g *metricGroup) UnmarshalJSON(data []byte) error {
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
