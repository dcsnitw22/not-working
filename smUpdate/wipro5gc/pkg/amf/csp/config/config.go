package config

import (
	"bytes"
	"os"
	"time"

	"k8s.io/klog"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
)

const (
	// DefaultConfigFilePath = "/go/src/w5gc.io/wipro5gcore/configs/amf"
	DefaultConfigFilePath = "wipro5gc/configs/amf"
	DefaultConfigFileName = "csp"
	DefaultEtcdConfigKey  = "/w5gc/configs/amf/csp/csp.json"
	DefaultEtcdServer     = "http://localhost:2379"
	DefaultEtcdConfigType = "json"
	DefaultConfigType     = "json"
)

type CspConfig struct {
	Version        string
	NodeInfo       AmfNodeInfo
	N11SmfNodes    []N11SmfNodeInfo
	GrpcServerInfo GrpcServerInfoConfig
}

type AmfNodeInfo struct {
	NodeId          string
	ApiPort         string
	LoadControl     bool
	OverloadControl bool
}

type N11SmfNodeInfo struct {
	NodeId string //IP ADDRESS OR DNS NAME (PREFERABLE)
	Port   string
	//ApiRoot string
	//URL string
}

type GrpcServerInfoConfig struct {
	ServerIP   string
	ServerPort string
}

type Update struct {
}

var defaultCspConfig = []byte(`
{
    "Version": "1.0",
    "NodeInfo":{
        "NodeId": "127.0.0.1",
		"ApiPort": ":8083"
    },
    "N11SmfNodes":[
        {
            "NodeId": "pdusmsp-service.pdusmsp.svc.cluster.local",
			"Port": "8080"
        }
    ],
	"GrpcServerInfo":{
		"ServerIP": "0.0.0.0",
		"ServerPort": "50054"
	}
}`)

var ConfigChannel chan CspConfig
var CspCfg, CspRuntimeCfg CspConfig

func InitConfig(cfgFile string, etcdServer string, etcdConfigKey string, resetFlag bool) (*CspConfig, error) {

	runtime_viper := viper.New()
	runtime_viper.Set("Verbose", true)
	//viper.Set("LogFile", LogFile)
	var etcdConfig bool

	// If Csp is not reset try to configure from etcd
	if resetFlag != true {
		if etcdServer == "" {
			klog.Errorf("Csp configuration using default etcd server %s", DefaultEtcdServer)
			etcdServer = DefaultEtcdServer
		}
		if etcdConfigKey == "" {
			klog.Errorf("Csp configuration using default etcd server key %s", DefaultEtcdConfigKey)
			etcdConfigKey = DefaultEtcdConfigKey
		}

		runtime_viper.AddRemoteProvider("etcd", etcdServer, etcdConfigKey)
		runtime_viper.SetConfigType(DefaultEtcdConfigType)
		err := runtime_viper.ReadRemoteConfig()

		if err == nil {
			klog.Info("Csp configured using etcd")
			etcdConfig = true
		}

		klog.Error("Csp configuration using etcd failed")
	}

	if !etcdConfig {
		klog.Info("Csp configuration using config file")

		// Check for config file parameter
		if cfgFile != "" {
			runtime_viper.SetConfigFile(cfgFile)
		} else {
			// Find home directory.
			home, err := homedir.Dir()
			if err != nil {
				klog.Errorf("[Unable to get home directory] %s", err.Error())
			}

			// Search config in home directory with name "go/src/w5gc.io/wipro5gcore/configs/" .
			runtime_viper.AddConfigPath(home + DefaultConfigFilePath)
			runtime_viper.SetConfigName(DefaultConfigFileName)
		}

		runtime_viper.AutomaticEnv()

		if err := runtime_viper.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); ok {
				klog.Error("[Config file not found in viper]")
			} else {
				klog.Error("[Unable to read config file in viper]")
			}
			klog.Errorf("[Unable to read config file, setting default values] %s", err.Error())
			runtime_viper.SetConfigType(DefaultConfigType)
			if err := runtime_viper.ReadConfig(bytes.NewBuffer(defaultCspConfig)); err != nil {
				klog.Errorf("[Unable to read config file with default values] %s", err.Error())
				os.Exit(1)
			}
		}
	}

	// Need to set default vaues TODO GURU

	// unmarshal config in viper to  config
	err := runtime_viper.Unmarshal(&CspCfg)
	if err != nil {
		klog.Fatalf("Unable to decode csp config into struct, %v", err)
		os.Exit(1)
	}

	// Write config to remote TODO GURU

	if etcdConfig {
		ConfigChannel = make(chan CspConfig)

		// Start a goroutine to watch remote config changes forever
		// Watch changes for config file? TODO GURU
		go func() {
			for {
				// delay after each request
				time.Sleep(time.Second * 5)

				// currently, only tested with etcd support
				err := runtime_viper.WatchRemoteConfig()
				if err != nil {
					klog.Errorf("unable to read remote csp config: %v", err)
					continue
				}

				// unmarshal new config into our runtime config struct
				runtime_viper.Unmarshal(&CspRuntimeCfg)

				// Notify csp event handler of the change
				ConfigChannel <- CspRuntimeCfg
			}
		}()
	}

	return &CspCfg, err
}
