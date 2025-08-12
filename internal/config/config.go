package config

import (
	"main/pkg/logger"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Listen struct {
		Addr   string
		BindIP string `env:"BIND_IP" env-default:"127.0.0.1"`
		Port   string `env:"LISTEN_PORT" env-default:"8000"`
	}
	Postgresql struct {
		Host     string `env:"PSQL_HOST"`
		Port     string `env:"PSQL_PORT"`
		Database string `env:"PSQL_NAME"`
		Username string `env:"PSQL_USER"`
		Password string `env:"PSQL_PASSWORD"`
	}
}

var instance *Config
var once sync.Once

func GetConfig() *Config {

	once.Do(func() {
		l := logger.GetLogger()
		l.Infoln("reading app configuration")
		instance = &Config{}
		err := cleanenv.ReadConfig(".env", instance)
		if err != nil {
			l.Fatalln("read app configuration error")
		}
		instance.Listen.Addr = instance.Listen.BindIP + ":" + instance.Listen.Port
		l.Infoln("reading config OK")
	})
	return instance
}
