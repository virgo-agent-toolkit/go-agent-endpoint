package main

import ()

func main() {
	stgs := new(settings)
	stgs.EndpointAPIAddr = ":8988"
	stgs.FileServerAddr = ":8989"
	stgs.FileServerRoot = "/tmp/upgrading"
	ctrl := &controller{stgs: stgs}
	ctrl.Serve()
}
