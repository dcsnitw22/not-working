package main

import (
        "os"
        "k8s.io/klog"

        "w5gc.io/wipro5gcore/cmd/smf/pdusmsp/app"
)

func main() {
        rootCmd := app.NewPdusmspRootCommand()
        err := rootCmd.Execute()

        if err != nil {
                os.Exit(1)
        }
        klog.Info("SMF PDU SMS Stopped")

}

