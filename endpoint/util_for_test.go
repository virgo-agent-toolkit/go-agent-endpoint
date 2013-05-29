package endpoint

import (
	"encoding/json"
	"io"
	"launchpad.net/gocheck"
	"net"
)

const FIXTURE_SERVER = "localhost:8099"

var PROTOCOL_FIXTURE_PREFIX = "http://" + FIXTURE_SERVER + "/protocol/"

func simpleExpectingTest(c *gocheck.C, endpoint_server string, fixture_res string, fixture_rsp string) (rsp_exp, rsp_test map[string]interface{}) {
	conn, err := net.Dial("tcp", endpoint_server)
	c.Assert(err, gocheck.IsNil)
	defer conn.Close()

	_, err = io.WriteString(conn, fixture_res)
	if err != io.EOF {
		c.Assert(err, gocheck.IsNil)
	}

	err = json.Unmarshal([]byte(fixture_rsp), &rsp_exp)
	c.Assert(err, gocheck.IsNil)
	err = json.NewDecoder(conn).Decode(&rsp_test)
	c.Assert(err, gocheck.IsNil)

	return rsp_exp, rsp_test
}

func assertEqualMapItem(c *gocheck.C, rsp_exp map[string]interface{}, rsp_test map[string]interface{}, key string) {
	c.Assert(rsp_exp[key], gocheck.NotNil)
	c.Assert(rsp_exp[key], gocheck.Equals, rsp_test[key])
}
