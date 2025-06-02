package config

import (
	"bytes"
	"fmt"
	"os"
	"time"

	"k8s.io/klog"

	"github.com/fsnotify/fsnotify"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
)

const (
	// DefaultConfigFilePath = "/go/src/w5gc.io/wipro5gcore/configs/amf"
	DefaultConfigFilePath = "/wipro5gc/configs/amf"
	DefaultConfigFileName = "ngap"
	DefaultEtcdConfigKey  = "/w5gc/configs/amf/ngap.json"
	DefaultEtcdServer     = "http://localhost:2379"
	DefaultEtcdConfigType = "json"
	DefaultConfigType     = "json"
)

type NgapConfig struct {
	Version        string
	NodeInfo       AmfNodeInfo
	GrpcClientInfo GrpcClientInfoConfig //for csp
	GrpcServerInfo GrpcServerInfoConfig //for pdusmsp
}

type AmfNodeInfo struct {
	NodeId  string
	ApiPort string
}

type GrpcClientInfoConfig struct {
	ClientIP   string
	ClientPort string
}

type GrpcServerInfoConfig struct {
	ServerIP   string
	ServerPort string
}

type Update struct {
}

var defaultNgapConfig = []byte(`
{
        "Version": "1.0",
        "NodeInfo":{
            "NodeId": "0.0.0.0",
            "ApiPort": "38412"
        },
		"GrpcClientInfo":{
			"ClientIP": "grpcnasamf-service.nasmm.svc.cluster.local",
			"ClientPort": "50052"
		},
		"GrpcServerInfo":{
			"ServerIP": "0.0.0.0",
			"ServerPort": "50055"
		}
}`)

var ConfigChannel chan NgapConfig
var NgapCfg, NgapRuntimeCfg NgapConfig

func InitConfig(cfgFile string, etcdServer string, etcdConfigKey string, resetFlag bool) (*NgapConfig, error) {

	runtime_viper := viper.New()
	runtime_viper.Set("Verbose", true)
	//viper.Set("LogFile", LogFile)
	var etcdConfig bool

	// If Ngap is not reset try to configure from etcd
	if resetFlag != true {
		if etcdServer == "" {
			klog.Errorf("Ngap configuration using default etcd server %s", DefaultEtcdServer)
			etcdServer = DefaultEtcdServer
		}
		if etcdConfigKey == "" {
			klog.Errorf("Ngap configuration using default etcd server key %s", DefaultEtcdConfigKey)
			etcdConfigKey = DefaultEtcdConfigKey
		}

		runtime_viper.AddRemoteProvider("etcd", etcdServer, etcdConfigKey)
		runtime_viper.SetConfigType(DefaultEtcdConfigType)
		err := runtime_viper.ReadRemoteConfig()

		if err == nil {
			klog.Info("Ngap configured using etcd")
			etcdConfig = true
		}

		klog.Error("Ngap configuration using etcd failed")
	}

	if !etcdConfig {
		klog.Info("Ngap configuration using config file")

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
			klog.Info(home)
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
			if err := runtime_viper.ReadConfig(bytes.NewBuffer(defaultNgapConfig)); err != nil {
				klog.Errorf("[Unable to read config file with default values] %s", err.Error())
				os.Exit(1)
			}
		}
		runtime_viper.OnConfigChange(func(in fsnotify.Event) {
			fmt.Println("Config file changed: ", in.Name)
			// unmarshal new config into our runtime config struct
			runtime_viper.Unmarshal(&NgapRuntimeCfg)
			// Notify ngap event handler of the change
			ConfigChannel <- NgapRuntimeCfg
		})
		runtime_viper.WatchConfig()
	}

	// Need to set default vaues TODO GURU

	// unmarshal config in viper to  config
	err := runtime_viper.Unmarshal(&NgapCfg)
	if err != nil {
		klog.Fatalf("Unable to decode ngap config into struct, %v", err)
		os.Exit(1)
	}

	// Write config to remote TODO GURU
	// klog.Info(etcdConfig)
	// etcdConfig = true
	if etcdConfig {
		ConfigChannel = make(chan NgapConfig)

		// Start a goroutine to watch remote config changes forever
		// Watch changes for config file? TODO GURU
		go func() {
			for {
				// delay after each request
				time.Sleep(time.Second * 5)

				// currently, only tested with etcd support
				err := runtime_viper.WatchRemoteConfig()
				if err != nil {
					klog.Errorf("unable to read remote ngap config: %v", err)
					continue
				}

				// unmarshal new config into our runtime config struct
				runtime_viper.Unmarshal(&NgapRuntimeCfg)

				// Notify ngap event handler of the change
				ConfigChannel <- NgapRuntimeCfg
			}
		}()
	}
	// temperory| jugad
	/*go func() {
		ConfigChannel = make(chan NgapConfig)
		klog.Info(NgapCfg)
		ConfigChannel <- NgapCfg
	}()*/
	return &NgapCfg, err
}
