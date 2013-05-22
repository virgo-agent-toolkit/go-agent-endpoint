go-agent-endpoint
=================

An experiment in creating an endpoint for the Virgo agent framework in Go

## How to test
In `go-agent-endpoint/virgovm`, you can vagrantup a VM that has an environment set up to test virgo and the endpoint.

### build VM
* Make sure [virgo](https://github.com/racker/virgo) is in the same directory where `go-agent-endpoint` is.
```
repo/
├── go-agent-endpoint
└── virgo
```

* Create the VM and initialize:
```
$ cd go-agent-endpoint/virgovm
$ vagrant up
```
Checkout `go-agent-endpoint/virgovm/provision.sh` for provision process (it runs in VM). It basically does following stuff:

* install necessary packages
* change permissions
* installs go if it's absent
* bind /data/O_O
* modify .profile of vagrant user to import /data/O_O/conf/rc, which sets up various environmental variables

### filesystem
The vagrant VM creates following NFS shared folder:
```
  name                          VM path                                 host path
--------  ------------------------------------------------------      -------------
endpoint: "/data/gopath/src/github.com/racker/go-agent-endpoint"  --> ".." (go-agent-endpoint)
virgo"    "/data/virgo                                         "  --> "../../virgo"
```
During provision process, `provision.sh` mounts `/data/gopath/src/github.com/racker/go-agent-endpoint/virgovm/O_O` to `/data/O_O`.

### test
```shell
$ vagrant ssh
```
In VM:
```shell
$ cd /data/O_O
$ make
```

### develop
Just hack in `go-agent-endpoint` or `virgo` on your host machine. Since they are mounted through NFS into VM, changes are gonna reflect in VM.
