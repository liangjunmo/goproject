package v1

import (
	"github.com/liangjunmo/gocode"

	"github.com/liangjunmo/goproject/internal/server/servercode"
)

var i18n = map[string]map[gocode.Code]string{
	"zh_CN": zhCn,
}

var zhCn = map[gocode.Code]string{
	servercode.OK:                  "OK",
	servercode.Unknown:             "未知错误",
	servercode.Timeout:             "请求超时",
	servercode.NotFound:            "资源不存在",
	servercode.InvalidRequest:      "请求错误",
	servercode.InternalServerError: "服务端错误",

	servercode.UserAlreadyExists: "用户已存在",
}
