VIRGO_DIR=/data/gopath/src/github.com/racker/go-agent-endpoint/testing/virgo
ENDPOINT_PKG=/data/gopath/src/github.com/racker/go-agent-endpoint/endpoint
MONITORING_PKG=/data/gopath/src/github.com/racker/go-agent-endpoint/examples/monitoring

default:
	@echo Target: integration test_agent test_server

integration: build_endpoint_bin build_agent integration/test.pem
	@echo
	@echo -------------------------  Integration Test  -------------------------
	integration/test.sh $(MONITORING_PKG)/monitoring $(VIRGO_DIR)/out/Debug

integration/test.pem:
	openssl req -x509 -newkey rsa:2048 -keyout /tmp/key.pem -out /tmp/cert.pem -days 3653 -nodes -subj '/C=/ST=/L=/O=/OU=/CN=/emailAddress=email@example.com'
	cat /tmp/key.pem /tmp/cert.pem > integration/test.pem

build_endpoint_bin:
	cd $(MONITORING_PKG) && go build

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
