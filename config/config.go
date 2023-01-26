package config

import (
	"log"
	"react-go-mybackend/utils"

	"gopkg.in/ini.v1"
)

type ConfigList struct {
	Port    string
	LogFile string
}

var Config ConfigList

func init() {
	LoadConfig()
	utils.Logging(Config.LogFile)
}

func LoadConfig() {
	cfg, err := ini.Load("config.ini")
	if err != nil {
		log.Fatalln(err)
	}
	Config = ConfigList{
		Port:    cfg.Section("web").Key("port").MustString("8080"),
		LogFile: cfg.Section("web").Key("logfile").String(),
	}
}
