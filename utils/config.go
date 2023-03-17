package utils

import (
	"fmt"
	"github.com/spf13/viper"
)

var Conf *viper.Viper

func LoadConf() {

	Conf = viper.New()
	Conf.SetConfigFile("config.json")
	Conf.SetConfigType("json")
	Conf.AddConfigPath(".")

	err := Conf.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	return
}
