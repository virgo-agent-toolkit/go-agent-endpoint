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
