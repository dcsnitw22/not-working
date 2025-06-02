package app

import (
        "os"
        "sync"

        "k8s.io/klog"
        "github.com/benbjohnson/clock"

        "w5gc.io/wipro5gcore/pkg/smf/pdusmsp"
        "w5gc.io/wipro5gcore/pkg/smf/pdusmsp/config"
)


func RunPdusmsp(cfgFile string, etcdServer string, etcdConfigKey string, reset bool){
        var wg sync.WaitGroup
        appClock := clock.New()
        start := appClock.Now()
        // Configure PDU SMS
        cfg, err := config.InitConfig(cfgFile, etcdServer, etcdConfigKey, reset)
        if err != nil {
                klog.Fatal("Unable to configure PDU SMS")
                os.Exit(1)
        }

        // Initialize PDU SMS
        u, ok := pdusmsp.NewPdusmsp(cfg, start)
        if !ok {
                klog.Fatal("Unable to Initialize PDU SMS")
                os.Exit(1)
        }

        // Start PDU SMS
        u.Run(config.ConfigChannel)

        wg.Wait()
}

