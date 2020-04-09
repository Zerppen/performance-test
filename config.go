package main

import (
	"fmt"

	"path/filepath"

	"github.com/spf13/viper"
	//"yonghui.cn/blockchain/turbot/backend/service"
)

type PeerConfig struct {
	IpAddress string
	Port      int
}

//var FabClient service.FabricClient

/**
function：InitConfig 初始化yaml配置文件
params：
	- string yamlConfigPath 配置文件config.yaml
return:
	- error 返回错误信息
*/
func InitConfig(yamlConfig string) error {
	fmt.Printf("InitYamlConfig configpath:%s \n", yamlConfig)
	fullPath, err := filepath.Abs(yamlConfig)
	if err != nil {
		fmt.Printf("The file path is error: %s", err.Error())
	}
	viper.SetConfigFile(fullPath)
	err = viper.ReadInConfig()
	if err != nil {
		return err
	}
	return nil
}

func GetConfigFile() string {
	return viper.ConfigFileUsed()
}

func GetRestfulListenAddress() string {
	return viper.GetString("server.restful.listenAddress")
}

func GetPeers() (peers []PeerConfig, err error) {
	err = viper.UnmarshalKey("server.fabric.peers", &peers)
	return
}
