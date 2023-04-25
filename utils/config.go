package utils

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
)

var Conf *viper.Viper

func LoadConf() {

	Conf = viper.New()
	Conf.SetConfigName("config.json")
	Conf.SetConfigType("json")
	Conf.AddConfigPath(".")
	Conf.AddConfigPath("./../../") //for testing purposes
	Conf.AddConfigPath("./../")    //for testing purposes
	err := Conf.ReadInConfig()
	if err != nil {
		log.Panic(fmt.Errorf("fatal error config file: %w", err))
	}

	return
}
