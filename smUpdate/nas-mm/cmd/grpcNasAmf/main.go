package main

import (
	"nasMM/pkg/config"
	grpcserver "nasMM/pkg/ngapNas/server"
	"time"

	"k8s.io/klog"
)

func main() {
	//Start the grpc Server
	config, err := config.InitConfig()
	if err != nil {
		klog.Info(err)
	}

	grpcServer := grpcserver.NewGrpcServer(config.NgapNas)

	grpcServer.Start()

	<-time.After(time.Second)
}
