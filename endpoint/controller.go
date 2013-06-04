package endpoint

import (
	"encoding/json"
	"io"
	"net/http"
)

const (
	CONTROLLER_API_AUTH = "/ctrl/auth"
)

type controller struct {
	host string
}

func newController(host string) *controller {
	return &controller{host: host}
}

func (c *controller) call(api string, req map[string]interface{}) (resp *http.Response, err error) {
	reader, writer := io.Pipe()
	go func() {
		json.NewEncoder(writer).Encode(req)
		writer.Close()
	}()
	resp, err = http.Post("http://"+c.host+api, "application/json", reader)
	return
}

func (c *controller) Authenticate(agentName string, agentId string, token string) bool {
	resp, err := c.call(CONTROLLER_API_AUTH, map[string]interface{}{"agent_name": agentName, "agent_id": agentId, "token": token})
	return nil == err && 200 == resp.StatusCode
}
