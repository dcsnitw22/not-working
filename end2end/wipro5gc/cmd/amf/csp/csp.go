package main

import (
	"os"

	"k8s.io/klog"
	"w5gc.io/wipro5gcore/cmd/amf/csp/app"
)

func main() {
	rootCmd := app.NewCspRootCommand()
	err := rootCmd.Execute()

	if err != nil {
		os.Exit(1)
	}
	klog.Info("AMF PDU SMS Stopped")

}
