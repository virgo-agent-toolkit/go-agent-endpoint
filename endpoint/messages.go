package endpoint

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

var (
	// UnmarshalFieldError is an error occured when unmarshaling json object,
	// mostly due to missing of field(s).
	UnmarshalFieldError = errors.New("Unmarshal field error")

	// NilPointerError is an error returned when UnmarshalJSON is called on a nil
	// pointer.
	NilPointerError = errors.New("Unmarshal function called on a nil pointer")
)

// Heartbeat is the Params and Result for heartbeat message.
type Heartbeat struct {
	Timestamp time.Time `json:"timestamp"`
}

// MarshalJSON marshals h into json bytes
func (h Heartbeat) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`{"timestamp": %d}`, h.Timestamp.UnixNano()/1e6)), nil
}

// UnmarshalJSON unmarshals json bytes into h
func (h *Heartbeat) UnmarshalJSON(data []byte) error {
	if h == nil {
		return NilPointerError
	}
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

// HelloParams is the Params for handshake.hello request, used for
// authentication
type HelloParams struct {
	Token          string `json:"token"`
	AgentID        string `json:"agent_id"`
	AgentName      string `json:"agent_name"`
	ProcessVersion string `json:"process_version"`
	BundleVersion  string `json:"bundle_version"`
}

// HelloResult is the Result for handshake.hello response
type HelloResult struct {
	HeartbeatInterval string `json:"heartbeat_interval"`
}

// GetVersionResult is the Result for get_version responses
type GetVersionResult struct {
	Version string `json:"version"`
}
