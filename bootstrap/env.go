package bootstrap

import (
	"log"
	"time"

	"github.com/spf13/viper"
)

type Env struct {
	DBHost         string        `mapstructure:"DB_HOST"`
	DBPort         string        `mapstructure:"DB_PORT"`
	DBUser         string        `mapstructure:"DB_USER"`
	DBPass         string        `mapstructure:"DB_PASS"`
	DBName         string        `mapstructure:"DB_NAME"`
	DBSSLMode      string        `mapstructure:"DB_SSLMODE"`
	ContextTimeout time.Duration `mapstructure:"CONTEXT_TIMEOUT"`
}

func NewEnv() *Env {
	env := Env{}
	viper.SetConfigFile(".env")

	_ = viper.ReadInConfig()
	err := viper.Unmarshal(&env)

	if err != nil {
		log.Fatal("Can't find the file .env : ", err)
	}

	return &env
}
