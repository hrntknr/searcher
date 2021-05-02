package main

import (
	"github.com/spf13/viper"
)

type config struct {
	Listen string
	Dsn    string
}

func loadConfig(fileName string, path []string) (*config, error) {
	viper.SetConfigName(fileName)
	for _, path := range path {
		viper.AddConfigPath(path)
	}
	viper.SetDefault("Listen", "0.0.0.0:8000")
	viper.SetDefault("Dsn", "user:pass@tcp(127.0.0.1:3306)/searcher?charset=utf8&parseTime=True&loc=Local")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	var config config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}
	return &config, nil
}
