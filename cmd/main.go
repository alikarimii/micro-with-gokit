package main

import (
	"github.com/alikarimii/micro-with-gokit/services/rest/cmd"
)

func main() {
	//grpc
	// @TODO

	// http
	s := cmd.RunHttp()
	go s.Start()
	s.WaitForStopSignal()
}
