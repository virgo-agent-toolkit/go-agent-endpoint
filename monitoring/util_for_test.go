package monitoring

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"testing"
)

const FIXTURE_SERVER = "localhost:8099"

var PROTOCOL_FIXTURE_PREFIX = "http://" + FIXTURE_SERVER + "/protocol/"

func isEmptyErrorFromJSON(err interface{}) (isEmpty bool, ok bool) {
	if er, ok := err.(map[string]interface{}); er["message"] == "" {
		return true, ok
	} else {
		return false, ok
	}
}

func getFixture(name string) (fixture io.ReadCloser, err error) {
	rsp, err := http.Get(PROTOCOL_FIXTURE_PREFIX + name)
	if err != nil {
		return
	}
	if 200 != rsp.StatusCode {
		return nil, errors.New(fmt.Sprintf("Unexpected HTTP status code. (200 expected, got %v)", rsp.StatusCode))
	}
	return rsp.Body, err
}

func failIfError(t *testing.T, err error, operation string, args ...interface{}) {
	if err != nil {
		err_str := fmt.Sprintf("Error in {%s}: %v\n", operation, err)
		for _, arg := range args {
			err_str = err_str + fmt.Sprintf("%v\n", arg)
		}
		t.Error(err_str)
	}
}

func simpleExpectingTest(t *testing.T, endpoint_server string, fixture_res string, fixture_rsp string) (rsp_exp, rsp_test map[string]interface{}) {
	conn, err := net.Dial("tcp", endpoint_server)
	failIfError(t, err, "net.Dial endpoint server")

	_, err = io.WriteString(conn, fixture_res)
	if err != io.EOF {
		failIfError(t, err, "io.WriteString - write fixture to endpoint server")
	}

	err = json.Unmarshal([]byte(fixture_rsp), &rsp_exp)
	failIfError(t, err, "json Decoding - rsp_exp:", []byte(fixture_rsp))
	err = json.NewDecoder(conn).Decode(&rsp_test)
	failIfError(t, err, "json Decoding - rsp_test:", rsp_test)

	return rsp_exp, rsp_test
}

func assertEqualMapItem(t *testing.T, rsp_exp map[string]interface{}, rsp_test map[string]interface{}, key string) {
	if rsp_exp[key] == nil {
		t.Errorf("Error in test case or fixture. Item [%s] does not exist in fixture.\n", key, rsp_exp[key], rsp_test[key])
	}
	if rsp_exp[key] != rsp_test[key] {
		t.Errorf("Item [%s] assertion failed. Expected %v, got %v\n", key, rsp_exp[key], rsp_test[key])
	}
}
