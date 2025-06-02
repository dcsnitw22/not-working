package config

import (
	"bytes"
	"os"

	"github.com/spf13/viper"
	"k8s.io/klog"
)

const (
	DefaultConfigType = "json"
)

type GrpcCspServerConfig struct {
	CspServerIP   string
	CspServerPort string
}

type GrpcNgapServerConfig struct {
	NgapServerIP   string
	NgapServerPort string
}

type GrpcNgapDLServerConfig struct {
	DLNgapServerIP   string
	DLNgapServerPort string
}

type GrpcSeverConfig struct {
	CspGrpc GrpcCspServerConfig
	NgapNas GrpcNgapServerConfig
	DLNgap  GrpcNgapDLServerConfig
}

var defaultConfig = []byte(`
{
    "CspGrpc":{
        "CspServerIP": "csp-service.csp.svc.cluster.local",
        "CspServerPort": "50054"
    },
    "NgapNas":{
        "NgapServerIP": "0.0.0.0",
        "NgapServerPort": "50052"

    },
	"DLNgap":{
        "DLNgapServerIP":"ngap-svc.ngap.svc.cluster.local",
        "DLNgapServerPort":"50056"
	
	}
}
`)

func InitConfig() (GrpcSeverConfig, error) {
	var grpcConfig GrpcSeverConfig
	runtime_viper := viper.New()
	runtime_viper.Set("Verbose", true)
	runtime_viper.SetConfigType(DefaultConfigType)
	if err := runtime_viper.ReadConfig(bytes.NewBuffer(defaultConfig)); err != nil {
		klog.Errorf("[Unable to read config file with default values] %s", err.Error())
		os.Exit(1)
		return grpcConfig, err
	}

	// unmarshal config in viper to  config
	err := runtime_viper.Unmarshal(&grpcConfig)
	if err != nil {
		klog.Fatalf("Unable to decode pdusmsp config into struct, %v", err)
		os.Exit(1)
		return grpcConfig, err
	}

	return grpcConfig, nil

}
