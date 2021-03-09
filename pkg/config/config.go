package config

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
)

type Config struct {
	Addr          string `json:"addr"`
	Cluster       string `json:"cluster"`
	HeartbeatTime int64  `json:"heartbeat_time"`
}

var (
	Conf Config
)

func ReadConfig() {
	viper.SetConfigFile("./res/config.yaml") // 指定配置文件路径
	viper.SetConfigName("config")            // 配置文件名称(无扩展名)
	viper.SetConfigType("yaml")              // 如果配置文件的名称中没有扩展名，则需要配置此项
	viper.AddConfigPath("./res/")            // 查找配置文件所在的路径
	err := viper.ReadInConfig()              // 查找并读取配置文件
	if err != nil {                          // 处理读取配置文件的错误
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	err = viper.Unmarshal(&Conf)
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	if addr := os.Getenv("ADDRESS"); addr != "" {
		Conf.Addr = addr
	}
}

func init() {
	ReadConfig()
}
