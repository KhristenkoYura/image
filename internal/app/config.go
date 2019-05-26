package app

import (
	"github.com/spf13/viper"
	"os"
)

// Config - основной
type Config struct {
	App ConfigApp `mapstructure:"app"`
}

type ConfigApp struct {
	Name     string `mapstructure:"name"`
	AddrBind string `mapstructure:"addr_bind"`
}

func NewConfig() *Config {
	var c = &Config{}
	var v = viper.New()
	v.AddConfigPath("./config")

	v.SetConfigName("default")
	// Основной конфиг
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}

	// Локальный конфиг
	v.SetConfigName("config.local")
	if err := v.MergeInConfig(); err != nil {
		panic(err)
	}

	//Получаем тип окружения
	mode := v.GetString("app.mode")
	if len(mode) == 0 {
		mode = os.Getenv("MODE")
	}

	// Конфиг окружение
	switch mode {
	case "PRODUCTION":
		v.SetConfigName("production")
	default:
		v.SetConfigName("developer")
	}

	if err := v.MergeInConfig(); err != nil {
		panic(err)
	}

	if err := v.Unmarshal(c); err != nil {
		panic(err)
	}

	return c
}
