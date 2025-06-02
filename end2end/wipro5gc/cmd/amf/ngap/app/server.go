package app

import (
	"os"
	"sync"

	"github.com/benbjohnson/clock"
	"k8s.io/klog"
	ngap "w5gc.io/wipro5gcore/pkg/amf/ngap"
	"w5gc.io/wipro5gcore/pkg/amf/ngap/config"
)

func RunNgap(cfgFile string, etcdServer string, etcdConfigKey string, reset bool, requestType string) {
	var wg sync.WaitGroup
	appClock := clock.New()
	start := appClock.Now()
	// Configure NGAP
	cfg, err := config.InitConfig(cfgFile, etcdServer, etcdConfigKey, reset)
	if err != nil {
		klog.Fatal("Unable to configure ngap")
		os.Exit(1)
	}

	// Initialize NGAP
	u, ok := ngap.NewNgap(cfg, start)
	if !ok {
		klog.Fatal("Unable to Initialize NGAP")
		os.Exit(1)
	}
	klog.Info("Initialised NGAP")

	// Start NGAP
	u.Run()
	klog.Info("Running NGAP")
	wg.Wait()
}
