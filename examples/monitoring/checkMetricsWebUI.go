package main

import (
    "crypto/md5"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "os"
    "path"
    "strconv"
    "sync"
    "time"

    "github.com/virgo-agent-toolkit/go-agent-endpoint/endpoint"
)

// checkMetricsWebUIHandler presents metrics (from only first connected agent)
// in web browser.
type checkMetricsWebUIHandler struct {
    webDir       string
    webClients   map[string]chan webRsp
    webClientsMu *sync.RWMutex
}

type webRsp struct {
    ClientID  string
    Data      float64
    Name      string
    AgentName string
}

func newCheckMetricsWebUIHandler(laddr string) *checkMetricsWebUIHandler {
    ret := new(checkMetricsWebUIHandler)
    ret.webClients = make(map[string]chan webRsp)
    ret.webClientsMu = new(sync.RWMutex)
    ret.webDir = path.Join(os.Getenv("GOPATH"), "src", "github.com", "virgo-agent-toolkit", "go-agent-endpoint", "examples", "monitoring", "web")

    mux := http.NewServeMux()
    mux.Handle("/", http.FileServer(http.Dir(ret.webDir)))
    mux.HandleFunc("/data", ret.webHandleData)
    mux.HandleFunc("/config", func(writer http.ResponseWriter, request *http.Request) {
        json.NewEncoder(writer).Encode(map[string]interface{}{
            "interval": scheduleInterval * 1000,
        })
    })

    go http.ListenAndServe(laddr, mux)

    return ret
}

func checkMetricsWebUIGetWebRsp(param *checkMetricsPostParams) ([]webRsp, bool) {
    switch param.CheckID {
    case "check-00":
        // agent.memory
        if len(param.Metrics) == 0 {
            return nil, false
        }
        v, err := strconv.ParseFloat(param.Metrics[0].Metrics["actual_used"].Value, 64)
        if err != nil {
            return nil, false
        }
        return []webRsp{{Name: "memory_used", Data: v}}, true
    case "check-03":
        // agent.network - eth0
        if len(param.Metrics) == 0 {
            return nil, false
        }
        rx, err := strconv.ParseFloat(param.Metrics[0].Metrics["rx_bytes"].Value, 64)
        if err != nil {
            return nil, false
        }
        tx, err := strconv.ParseFloat(param.Metrics[0].Metrics["tx_bytes"].Value, 64)
        if err != nil {
            return nil, false
        }
        return []webRsp{{Name: "eth0_rx", Data: rx}, {Name: "eth0_tx", Data: tx}}, true
    case "check-04":
        // agent.network - eth1
        if len(param.Metrics) == 0 {
            return nil, false
        }
        rx, err := strconv.ParseFloat(param.Metrics[0].Metrics["rx_bytes"].Value, 64)
        if err != nil {
            return nil, false
        }
        tx, err := strconv.ParseFloat(param.Metrics[0].Metrics["tx_bytes"].Value, 64)
        if err != nil {
            return nil, false
        }
        return []webRsp{{Name: "eth1_rx", Data: rx}, {Name: "eth1_tx", Data: tx}}, true
    case "check-05":
        // agent.cpu
        if len(param.Metrics) == 0 {
            return nil, false
        }
        v, err := strconv.ParseFloat(param.Metrics[0].Metrics["usage_average"].Value, 64)
        if err != nil {
            return nil, false
        }
        return []webRsp{{Name: "cpu_usage", Data: v}}, true
    default:
        return nil, false
    }
}

func (c checkMetricsWebUIHandler) multiplex(param *checkMetricsPostParams, source string) {
    rsps, ok := checkMetricsWebUIGetWebRsp(param)
    if !ok {
        return
    }
    c.webClientsMu.RLock()
    for _, rsp := range rsps {
        rsp.AgentName = source
        for clientID, client := range c.webClients {
            rsp.ClientID = clientID
            select {
            case client <- rsp:
            default:
            }
        }
    }
    c.webClientsMu.RUnlock()
}

func (c *checkMetricsWebUIHandler) Handle(req *endpoint.Request, responder *endpoint.Responder, connCtx endpoint.ConnContext) endpoint.HandleCode {
    params := new(checkMetricsPostParams)
    err := json.Unmarshal(req.Params, params)
    if err != nil { // parsing failed, should not go on
        fmt.Printf("parsing check_metrics.post Params failed: %v\n", err)
        return endpoint.FAIL
    }
    c.multiplex(params, req.Source)
    return endpoint.DECLINED
}

func (c *checkMetricsWebUIHandler) webHandleData(writer http.ResponseWriter, request *http.Request) {
    fmt.Println("got a web request")
    clientID := request.URL.Query().Get("clientID")
    if clientID == "" {
        fmt.Println("generating client id")
        hash := md5.New()
        hash.Write([]byte(time.Now().String()))
        clientID = fmt.Sprintf("%x", hash.Sum(nil))
        c.webClientsMu.Lock()
        c.webClients[clientID] = make(chan webRsp, 64)
        c.webClientsMu.Unlock()
    }
    fmt.Println(clientID)
    c.webClientsMu.RLock()
    fmt.Println(c.webClients)
    ch := c.webClients[clientID]
    c.webClientsMu.RUnlock()
    if ch != nil {
        fmt.Println("encoding")
        json.NewEncoder(io.MultiWriter(os.Stdout, writer)).Encode(<-ch)
    }
}
