package endpoint

import (
	"encoding/json"
	"io"
	"launchpad.net/gocheck"
	"net"
)

const FIXTUREServer = "localhost:8099"

var FIXTUREProtocolPrefix = "http://" + FIXTUREServer + "/protocol/"

func simpleExpectingTest(c *gocheck.C, conn net.Conn, fixtureRes string, fixtureRsp string) (rspExp, rspTest map[string]interface{}) {
	_, err := io.WriteString(conn, fixtureRes)
	if err != io.EOF {
		c.Assert(err, gocheck.IsNil)
	}

	err = json.Unmarshal([]byte(fixtureRsp), &rspExp)
	c.Assert(err, gocheck.IsNil)
	err = json.NewDecoder(conn).Decode(&rspTest)
	c.Assert(err, gocheck.IsNil)

	return rspExp, rspTest
}

func assertEqualMapItem(c *gocheck.C, rspExp map[string]interface{}, rspTest map[string]interface{}, key string) {
	c.Assert(rspExp[key], gocheck.NotNil)
	c.Assert(rspExp[key], gocheck.Equals, rspTest[key])
}
