package monitoring

import (
	"encoding/json"
	"testing"
)

func TestCustomJSONMarshals(t *testing.T) {
	/* MetricGroup (for CheckMetricsPost request) */
	var mg, mg_ *MetricGroup
	metrics := make(map[string]Metric)
	metrics["m1"] = Metric{Type: "int", Value: "123", Unit: "bytes"}
	metrics["m2"] = Metric{Type: "float", Value: "0.123", Unit: "percent"}
	mg = &MetricGroup{Prefix: "this_is_the_prefix", Metrics: metrics}
	js, err := json.Marshal(mg)
	failIfError(t, err, "json.Marshal(mg)", mg)
	err = json.Unmarshal(js, &mg_)
	failIfError(t, err, "json.Unmarshal(mg_)", string(js))
	if mg.Prefix != mg_.Prefix {
		t.Errorf("Prefix changed during Marshal/Unmarshal")
	}
	compare := func(m1, m2 *Metric) bool {
		return m1.Type == m2.Type && m1.Value == m2.Value && m1.Unit == m2.Unit
	}
	for k, v := range mg.Metrics {
		if v_, ok := mg_.Metrics[k]; !(ok && compare(&v, &v_)) {
			t.Errorf("Item [%s] changed during Marshal/Unmarshal (<%v> != <%v>)", k, v, v_)
		}
	}
}
