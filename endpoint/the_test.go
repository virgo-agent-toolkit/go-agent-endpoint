package endpoint

import (
	"encoding/json"
	"launchpad.net/gocheck"
	"net"
)

func (s *TestSuite) init(c *gocheck.C, endpoint_addr string) (*Endpoint, <-chan *Requester, net.Conn) {
	hub, requesters := NewHub()
	hub.Authenticator(DumbAuthenticatorDontUseMe(0), 0)

	config := EndpointConfig{}
	config.Hub = hub
	config.ListenAddr = endpoint_addr
	config.UpgradingFileServerAddr = "localhost:8080"

	server, err := NewEndpoint(config)
	c.Assert(err, gocheck.IsNil)

	server.Start()

	conn, err := net.Dial("tcp", endpoint_addr)
	c.Assert(err, gocheck.IsNil)

	_, rsp_test := simpleExpectingTest(c, conn, FIXTURE_handshake_hello_request, "{}")
	c.Assert(rsp_test["error"], gocheck.IsNil)

	return server, requesters, conn
}

func (s *TestSuite) TestHeartbeat(c *gocheck.C) {
	endpoint_addr := "localhost:9876"
	server, _, conn := s.init(c, endpoint_addr)

	rsp_exp, rsp_test := simpleExpectingTest(c, conn, FIXTURE_heartbeat_post_request, FIXTURE_heartbeat_post_response)
	c.Assert(rsp_test["error"], gocheck.IsNil)

	assertEqualMapItem(c, rsp_exp, rsp_test, "v")
	assertEqualMapItem(c, rsp_exp, rsp_test, "id")
	assertEqualMapItem(c, rsp_exp, rsp_test, "source")
	assertEqualMapItem(c, rsp_exp, rsp_test, "target")

	rsp_exp, rsp_test = simpleExpectingTest(c, conn, FIXTURE_heartbeat_post_request_invalid_version, FIXTURE_heartbeat_post_response)
	c.Assert(rsp_test["error"], gocheck.FitsTypeOf, map[string]interface{}{})
	msg := rsp_test["error"].(map[string]interface{})["message"]
	c.Assert(msg, gocheck.NotNil)
	c.Assert(msg, gocheck.Not(gocheck.Equals), "")
	assertEqualMapItem(c, rsp_exp, rsp_test, "v")

	conn.Close()
	server.Destroy()
}

func (s *TestSuite) TestProactive(c *gocheck.C) {
	endpoint_addr := "localhost:9876"
	server, requesters, conn := s.init(c, endpoint_addr)

	requester := <-requesters
	reply, err := requester.Send("proactive_test", FIXTURE_proactive_test_request)
	c.Assert(err, gocheck.IsNil)

	var req *Request
	err = json.NewDecoder(conn).Decode(&req)
	c.Assert(err, gocheck.IsNil)

	c.Assert(req.Method, gocheck.Equals, "proactive_test")

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
