package main

import (
	"encoding/json"
	"net/http"
)

const (
	CONTROLLER_API_AUTH = "/ctrl/auth"
)

type settings struct {
	EndpointAPIAddr string
	FileServerAddr  string
	FileServerRoot  string
}

type controller struct {
	stgs *settings
}

func (c *controller) Serve() {
	fileServerHandler := http.FileServer(http.Dir(c.stgs.FileServerRoot))
	go http.ListenAndServe(c.stgs.FileServerAddr, fileServerHandler)
	mux := c.buildMux()
	http.ListenAndServe(c.stgs.EndpointAPIAddr, mux)
}

func (c *controller) buildMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc(CONTROLLER_API_AUTH, c.ctrl_auth)
	return mux
}

func (c *controller) ctrl_auth(rspW http.ResponseWriter, req *http.Request) {
	authReq := c.readReq(req)
	_, ok1 := authReq["agent_name"]
	_, ok2 := authReq["agent_id"]
	_, ok3 := authReq["token"]
	if ok1 && ok2 && ok3 {
		rspW.WriteHeader(200)
	} else {
		rspW.WriteHeader(203)
	}
}

func (c *controller) readReq(req *http.Request) map[string]interface{} {
	defer req.Body.Close()
	ret := make(map[string]interface{})
	json.NewDecoder(req.Body).Decode(&ret)
	return ret
}
