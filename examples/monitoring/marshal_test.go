package main

import (
	"encoding/json"
	"launchpad.net/gocheck"
)

func (s *TestSuite) TestCustomJSONMarshals(c *gocheck.C) {
	/* MetricGroup (for CheckMetricsPost request) */
	var mg, mG *metricGroup
	mg = &metricGroup{Prefix: "this_is_the_prefix", Metrics: make(map[string]*metric)}
	mg.Metrics["m1"] = &metric{Type: "int", Value: "123", Unit: "bytes"}
	mg.Metrics["m2"] = &metric{Type: "float", Value: "0.123", Unit: "percent"}
	js, err := json.Marshal(mg)
	c.Assert(err, gocheck.IsNil)
	err = json.Unmarshal(js, &mG)
	c.Assert(err, gocheck.IsNil)
	c.Assert(mg.Prefix, gocheck.Equals, mG.Prefix)
	for k, v := range mg.Metrics {
		c.Assert(v, gocheck.DeepEquals, mG.Metrics[k])
	}
}
