package monitoring

import (
	"encoding/json"
	"launchpad.net/gocheck"
)

func (s *TestSuite) TestCustomJSONMarshals(c *gocheck.C) {
	/* MetricGroup (for CheckMetricsPost request) */
	var mg, mg_ *MetricGroup
	mg = &MetricGroup{Prefix: "this_is_the_prefix", Metrics: make(map[string]*Metric)}
	mg.Metrics["m1"] = &Metric{Type: "int", Value: "123", Unit: "bytes"}
	mg.Metrics["m2"] = &Metric{Type: "float", Value: "0.123", Unit: "percent"}
	js, err := json.Marshal(mg)
	c.Assert(err, gocheck.IsNil)
	err = json.Unmarshal(js, &mg_)
	c.Assert(err, gocheck.IsNil)
	c.Assert(mg.Prefix, gocheck.Equals, mg_.Prefix)
	for k, v := range mg.Metrics {
		c.Assert(v, gocheck.DeepEquals, mg_.Metrics[k])
	}
}
