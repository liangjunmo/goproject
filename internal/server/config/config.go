package config

import (
	"github.com/liangjunmo/goproject/internal/configtemplate"
)

type Template struct {
	Environment configtemplate.Environment `mapstructure:"environment"`
	Debug       bool                       `mapstructure:"debug"`
	API         struct {
		Addr           string `mapstructure:"addr"`
		JWTKey         string `mapstructure:"jwtKey"`
		UserCenterAddr string `mapstructure:"userCenterAddr"`
	} `mapstructure:"api"`
	UserCenter struct {
		Addr string `mapstructure:"addr"`
	} `mapstructure:"userCenter"`
	DB struct {
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

var (
	Config Template

	ProjectDir string
)

const (
	TraceIDKey string = "TraceID"

	GinCtxUserKey = "user_claims"
)
