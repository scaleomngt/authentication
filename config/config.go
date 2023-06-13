package config

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log"
	"os"
)

var Config *viper.Viper

func init() {
	//监听改变动态跟新配置
	go watchConfig()
	//加载配置
	loadConfig()
}

// 监听配置改变
func watchConfig() {
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Println("Config file changed:", e.Name)
		//改变重新加载
		loadConfig()
	})
}

// 加载配置
func loadConfig() {
	viper.SetConfigName("config")
	viper.AddConfigPath("./config/")
	//viper.AddConfigPath("/etc/env")
	err := viper.ReadInConfig()
	if err != nil {
		log.Println("Fatal error config file: ", err)
		os.Exit(-1)
	}
	//全局配置
	Config = viper.GetViper()
	log.Println(Config.AllSettings())
}
