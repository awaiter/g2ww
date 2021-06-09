package config

/*
当前不支持热更新配置文件
*/

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

// 以下结构体，仅供包内调用
type grafana struct {
	Url string
	Env string
}

type webhook struct {
	Type    string
	Url     string
	Msgtype string
	AtUser  string
}

type conf struct {
	Port    int
	Grafana grafana
	Webhook webhook
	//v       *viper.Viper
}

var Config conf

var _config *viper.Viper

func load() {
	var config *viper.Viper
	config = viper.New()
	config.SetConfigType("yaml")
	config.SetConfigName("config")
	config.AddConfigPath(".")

	err := config.ReadInConfig()
	if err != nil {
		err = fmt.Errorf("配置文件读取失败: %v", err)
		panic(err)
	}

	err = config.Unmarshal(&Config)
	if err != nil {
		err = fmt.Errorf("配置文件序列化失败: %v", err)
		panic(err)
	}

	_config = config

	return
}

func init() {
	fmt.Println("加载配置文件...")
	load()
	//spew.Dump(Config)
}

func GetConfig() *viper.Viper {
	if _config == nil {
		log.Fatalln("config not init...")
	}

	return _config
}
