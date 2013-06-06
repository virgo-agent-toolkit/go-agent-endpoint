package endpoint

import (
	"launchpad.net/gocheck"
	"net"
)

func (s *TestSuite) init(c *gocheck.C, endpoint_addr string) (*Endpoint, net.Conn) {
	hub := NewSimpleHub()
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

	return server, conn
}

func (s *TestSuite) TestHeartbeat(c *gocheck.C) {
	endpoint_addr := "localhost:9876"
	server, conn := s.init(c, endpoint_addr)

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
