VIRGO_DIR=/data/virgo
ENDPOINT_PKG=/data/gopath/src/github.com/racker/go-agent-endpoint/endpoint
MONITORING_PKG=/data/gopath/src/github.com/racker/go-agent-endpoint/monitoring

integration:
	@echo integration test not implemented...
	@echo try individual tests: \`make test_agent\` or \`make test_server\`

test_agent:
	@echo
	@echo -------------------------   Testing  agent   -------------------------
	(cd $(VIRGO_DIR); make test)
	@echo

test_server: test_endpoint test_monitoring

test_endpoint:
	@echo
	@echo -------------------------  Testing endpoint  -------------------------
	(cd $(ENDPOINT_PKG); go test)
	@echo

test_monitoring:
	@echo
	@echo ------------------------- Testing monitoring -------------------------
	(cd $(MONITORING_PKG); go test)
	@echo
