package common

import (
	"github.com/kayac/go-config"
	"sync"
)

var (
	envOnce = new(sync.Once)
	Conf    *Config // Singleton
)

type Config struct {
	Env  string   `yaml:"env"`
	Notice  NoticeConf  `yaml:"notice"`
	Aero AeroConf `yaml:"aero"`
}

type NoticeConf struct {
	Addr string `yaml:"addr"`
}

type AeroConf struct {
	Addr string `yaml:"addr"`
	Port string `yaml:"port"`
}

func SetupEnv() {
	envOnce.Do(func() {
		tmp := &Config{}
		err := config.LoadWithEnv(tmp, "resources/config.yaml")
		if err != nil {
			Logger.Error(err)
		}
		Conf = tmp
	})
}
