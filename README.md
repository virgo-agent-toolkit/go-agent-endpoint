go-agent-endpoint
=================

An experiment in creating an endpoint for the Virgo agent framework in Go

## How to test
* Clone the `go-agent-endpiont` repo; You don't have to worry about init/update submodules since Makefile ensures it.
* `make test` does pretty much everything -- create VM, provision VM, update code, and run the test.
* Alternatively, use `make ssh` to get into VM.
* `make reload` or `make clean` cleans up everything in VM and reloads them

## develop
Just hack in `go-agent-endpoint` or `virgo` on your host machine. Since they are mounted through NFS into VM, changes are gonna reflect in VM.
