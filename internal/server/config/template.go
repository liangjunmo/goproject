package config

import (
	"github.com/liangjunmo/goproject/internal/configtemplate"
)

type Template struct {
	Environment configtemplate.Environment `mapstructure:"environment"`
	Debug       bool                       `mapstructure:"debug"`
	API         struct {
		Addr   string `mapstructure:"addr"`
		JwtKey string `mapstructure:"jwt_key"`
	} `mapstructure:"api"`
	Db struct {
		Addr     string `mapstructure:"addr"`
		User     string `mapstructure:"user"`
		Password string `mapstructure:"password"`
		Database string `mapstructure:"database"`
	} `mapstructure:"db"`
	Redis struct {
		Addr     string `mapstructure:"addr"`
		Password string `mapstructure:"password"`
	} `mapstructure:"redis"`
}
