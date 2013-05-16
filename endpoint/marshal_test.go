package endpoint

import (
	"encoding/json"
	"testing"
	"time"
)

func TestCustomJSONMarshals(t *testing.T) {
	/* Heartbeat (for Heartbeat request/response) */
	var hb, hb_ *Heartbeat
	hb = &Heartbeat{time.Now()}
	js, err := hb.MarshalJSON()
	failIfError(t, err, "hb.MarshalJSON", hb)
	err = json.Unmarshal(js, &hb_)
	failIfError(t, err, "json.Unmarshal(hg_)", string(js))
	if hb.Timestamp.UnixNano()/1e6 != hb_.Timestamp.UnixNano()/1e6 {
		t.Errorf("Timestamp changed during Marshal/Unmarshal")
	}
}
