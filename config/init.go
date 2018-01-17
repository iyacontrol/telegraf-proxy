package config

import (
	"log"

	"github.com/BurntSushi/toml"
)

func InitConfig(fConfig string) {
	if fConfig == "" {
		fConfig = "telegtaf-proxy.toml"
	}
	if _, err := toml.DecodeFile(fConfig, &Cfg); err != nil {
		log.Fatalf("配置文件错误:%s\n", err)
	}

}
