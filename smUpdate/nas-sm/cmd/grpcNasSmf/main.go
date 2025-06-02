package main

import (
	grpcSMserver "nasMain/pkg/grpcSmfNas/server"
	"time"
)

func main() {
	//Start the grpc Server
	grpcSMserver.StartSmfNasGrpc()

	<-time.After(time.Second)
}
