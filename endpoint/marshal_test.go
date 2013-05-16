package endpoint

import (
	"encoding/json"
	"launchpad.net/gocheck"
	"time"
)

func (s *TestSuite) TestCustomJSONMarshals(c *gocheck.C) {
	/* Heartbeat (for Heartbeat request/response) */
	var hb, hb_ *Heartbeat
	hb = &Heartbeat{time.Now()}
	js, err := hb.MarshalJSON()
	c.Assert(err, gocheck.IsNil)
	err = json.Unmarshal(js, &hb_)
	c.Assert(err, gocheck.IsNil)
	if hb.Timestamp.UnixNano()/1e6 != hb_.Timestamp.UnixNano()/1e6 {
		c.Fatal("Timestamp changed during Marshal/Unmarshal")
	}
}
