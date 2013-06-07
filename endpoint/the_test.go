package endpoint

import (
	"encoding/json"
	"launchpad.net/gocheck"
	"net"
)

func (s *TestSuite) init(c *gocheck.C, endpointAddr string) (*Endpoint, <-chan *Requester, net.Conn) {
	hub, requesters := NewHub()
	hub.Authenticator(DumbAuthenticatorDontUseMe(0), 0)

	config := EndpointConfig{}
	config.Hub = hub
	config.ListenAddr = endpointAddr
	config.UpgradingFileServerAddr = "localhost:8080"

	server, err := NewEndpoint(config)
	c.Assert(err, gocheck.IsNil)

	server.Start()

	conn, err := net.Dial("tcp", endpointAddr)
	c.Assert(err, gocheck.IsNil)

	_, rspTest := simpleExpectingTest(c, conn, FIXTUREHandshakeHelloRequest, "{}")
	c.Assert(rspTest["error"], gocheck.IsNil)

	return server, requesters, conn
}

func (s *TestSuite) TestHeartbeat(c *gocheck.C) {
	endpointAddr := "localhost:9876"
	server, _, conn := s.init(c, endpointAddr)

	rspExp, rspTest := simpleExpectingTest(c, conn, FIXTUREHeartbeatPostRequest, FIXTUREHeartbeatPostResponse)
	c.Assert(rspTest["error"], gocheck.IsNil)

	assertEqualMapItem(c, rspExp, rspTest, "v")
	assertEqualMapItem(c, rspExp, rspTest, "id")
	assertEqualMapItem(c, rspExp, rspTest, "source")
	assertEqualMapItem(c, rspExp, rspTest, "target")

	rspExp, rspTest = simpleExpectingTest(c, conn, FIXTUREHeartbeatPostRequestInvalidVersion, FIXTUREHeartbeatPostResponse)
	c.Assert(rspTest["error"], gocheck.FitsTypeOf, map[string]interface{}{})
	msg := rspTest["error"].(map[string]interface{})["message"]
	c.Assert(msg, gocheck.NotNil)
	c.Assert(msg, gocheck.Not(gocheck.Equals), "")
	assertEqualMapItem(c, rspExp, rspTest, "v")

	conn.Close()
	server.Destroy()
}

func (s *TestSuite) TestProactive(c *gocheck.C) {
	endpointAddr := "localhost:9876"
	server, requesters, conn := s.init(c, endpointAddr)

	requester := <-requesters
	reply, err := requester.Send("proactiveTest", FIXTUREProactiveTestRequest)
	c.Assert(err, gocheck.IsNil)

	var req *Request
	err = json.NewDecoder(conn).Decode(&req)
	c.Assert(err, gocheck.IsNil)

	c.Assert(req.Method, gocheck.Equals, "proactiveTest")

	var params map[string]interface{}
	err = json.Unmarshal(req.Params, &params)
	c.Assert(err, gocheck.IsNil)
	c.Assert(params["hello"], gocheck.Equals, "world")

	connRsp := respondingTo(req)
	connRsp.Result = json.RawMessage(`{"world": "hello"}`)
	err = json.NewEncoder(conn).Encode(connRsp)
	c.Assert(err, gocheck.IsNil)

	rsp := <-reply
	var rspResult map[string]interface{}
	json.Unmarshal(rsp.Result, &rspResult)
	c.Assert(rspResult["world"], gocheck.Equals, "hello")

	conn.Close()
	server.Destroy()
}
