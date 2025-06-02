package main

import (
	"os"

	"k8s.io/klog"
	"w5gc.io/wipro5gcore/cmd/amf/ngap/app"
)

func main() {
	rootCmd := app.NewNgapRootCommand()
	err := rootCmd.Execute()

	if err != nil {
		klog.Error(err)
		os.Exit(1)
	}
	klog.Info("AMF NGAP Stopped")

}
