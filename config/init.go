package config

import (
	"log"

	"github.com/BurntSushi/toml"
)

func InitConfig(fConfig string) {
	if fConfig == "" {
		fConfig = "telegraf-proxy.toml"
	}
	if _, err := toml.DecodeFile(fConfig, &Cfg); err != nil {
		log.Fatalf("config err:%s\n", err)
	}

}
