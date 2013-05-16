package endpoint

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

var (
	UnmarshalFieldError = errors.New("Unmarshal field error")
)

type Heartbeat struct {
	Timestamp time.Time `json:"timestamp"`
}

func (h Heartbeat) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`{"timestamp": %d}`, h.Timestamp.UnixNano()/1e6)), nil
}

func (h *Heartbeat) UnmarshalJSON(data []byte) error {
	var tmp map[string]float64
	err := json.Unmarshal(data, &tmp)
	if err != nil {
		return err
	}
	msec, ok := tmp["timestamp"]
	if !ok {
		return UnmarshalFieldError
	}
	h.Timestamp = time.Unix(0, int64(msec)*1e6)
	return nil
}

type HelloParams struct {
	Token          string `json:"token"`
	AgentId        string `json:"agent_id"`
	AgentName      string `json:"agent_name"`
	ProcessVersion string `json:"process_version"`
	BundleVersion  string `json:"bundle_version"`
}

type HelloResult struct {
	HeartbeatInterval string `json:"heartbeat_interval"`
}

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

type GetVersionResult struct {
	Version string `json:"version"`
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
	var tmp []json.RawMessage
	err := json.Unmarshal(data, &tmp)
	if err != nil {
		return err
	}
	err = json.Unmarshal(tmp[0], &g.Prefix)
	if err != nil {
		return UnmarshalFieldError
	}
	err = json.Unmarshal(tmp[1], g.Metrics)
	if err != nil {
		return UnmarshalFieldError
	}
	return nil
}
