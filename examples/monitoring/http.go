package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path"
	"regexp"
)

var upgradeChannels = [...]string{"stable", "beta", "master"}
var upgradeURLRe, _ = regexp.Compile("^/upgrades/(.*)/(.*)$")
var upgradeHTTPDir = http.Dir(path.Join(os.Getenv("GOPATH"), "src", "github.com", "racker", "go-agent-endpoint", "examples", "monitoring", "httpForAgent"))

func httpServer(laddr string) {
	mux := http.NewServeMux()

	mux.Handle("/upgrades/", upgradeMux())
	mux.HandleFunc("/agent-crash-report", handleCrashReport)
	mux.HandleFunc("/", http.NotFound)

	http.ListenAndServe(laddr, mux)
}

func upgradeMux() *http.ServeMux {
	fileServer := http.FileServer(upgradeHTTPDir)
	mux := http.NewServeMux()
	for _, channel := range upgradeChannels {
		mux.HandleFunc(fmt.Sprintf("/upgrades/%s/", channel),
			func(rsp http.ResponseWriter, request *http.Request) {
				fmt.Printf("Upgrading request: %s\n", request.RequestURI)
				fileServer.ServeHTTP(rsp, request)
			})
	}
	mux.HandleFunc("/", http.NotFound)
	return mux
}

func getChannelAndFileName(requestURI string) (channel string, filename string, err error) {
	matches := upgradeURLRe.FindStringSubmatch(requestURI)
	if len(matches) < 3 {
		return "", "", errors.New("Bad upgrading request")
	}
	return matches[1], matches[2], nil
}

/*
func handleUpgrade(rsp http.ResponseWriter, request *http.Request) {
  channel, filename, err := getChannelAndFileName(request.RequestURI)
  if err != nil {
    fmt.Printf("Sending 400 in response to requestURI %s; %v\n", request.RequestURI, err)
    rsp.WriteHeader(400)
    rsp.Write([]byte(err.Error()))
    return
  }
}
*/

func handleCrashReport(rsp http.ResponseWriter, request *http.Request) {
	fmt.Printf("got a crash report %s\n", request.RequestURI)
	rsp.WriteHeader(404)
	rsp.Write(nil)
}
