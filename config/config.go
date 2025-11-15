package config

import (
	"github.com/spf13/viper"
	"fmt"
)

func LoadConfig() {
	viper.SetConfigFile("config.yaml")
	viper.ReadInConfig()
	
	fmt.Println(viper.GetString("loadbalancer.strategy"))
}
