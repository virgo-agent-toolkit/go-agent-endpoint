VIRGO_DIR=/data/virgo
GO_AGENT_ENDPOINT_PKG=/data/gopath/src/github.com/racker/go-agent-endpoint
ENDPOINT_PKG=/data/gopath/src/github.com/racker/go-agent-endpoint/endpoint
MONITORING_PKG=/data/gopath/src/github.com/racker/go-agent-endpoint/monitoring

default:
	@echo Target: integration test_agent test_server

integration: build_endpoint_bin build_agent
	@echo
	@echo -------------------------  Integration Test  -------------------------
	integration/test.sh $(GO_AGENT_ENDPOINT_PKG)/go-agent-endpoint $(VIRGO_DIR)/out/Debug

build_endpoint_bin:
	cd $(GO_AGENT_ENDPOINT_PKG) && go build

test_agent: build_agent
	@echo
	@echo -------------------------   Testing  agent   -------------------------
	cd $(VIRGO_DIR) && make test
	@echo

build_agent:
	@echo
	@echo -------------------------   Building agent   -------------------------
	cd $(VIRGO_DIR) && ./configure && make
	@echo

test_server: test_endpoint test_monitoring

test_endpoint:
	@echo
	@echo -------------------------  Testing endpoint  -------------------------
	cd $(ENDPOINT_PKG) && go test
	@echo

test_monitoring:
	@echo
	@echo ------------------------- Testing monitoring -------------------------
	cd $(MONITORING_PKG) && go test
	@echo
