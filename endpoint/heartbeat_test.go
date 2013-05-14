package endpoint

import (
	"testing"
)

func TestHeartbeat(t *testing.T) {
	endpoint_addr := "localhost:9876"
	server, err := NewServer(endpoint_addr)
	failIfError(t, err, "creating server")
	server.Start()

	rsp_exp, rsp_test := simpleExpectingTest(t, endpoint_addr, FIXTURE_heartbeat_post_request, FIXTURE_heartbeat_post_response)
	if noErr, ok := isEmptyErrorFromJSON(rsp_test["error"]); ok && !noErr {
		t.Errorf("error in parsing error, or error is expected to be null but is not null in response: %v", rsp_test["error"])
	}
	assertEqualMapItem(t, rsp_exp, rsp_test, "v")
	assertEqualMapItem(t, rsp_exp, rsp_test, "id")
	assertEqualMapItem(t, rsp_exp, rsp_test, "source")
	assertEqualMapItem(t, rsp_exp, rsp_test, "target")

	rsp_exp, rsp_test = simpleExpectingTest(t, endpoint_addr, FIXTURE_heartbeat_post_request_invalid_version, FIXTURE_heartbeat_post_response)
	if noErr, ok := isEmptyErrorFromJSON(rsp_test["error"]); ok && noErr {
		t.Errorf("error in parsing error, or error is expected but is null in response: %v", rsp_test)
	}
	assertEqualMapItem(t, rsp_exp, rsp_test, "v")

	server.Destroy()
}
