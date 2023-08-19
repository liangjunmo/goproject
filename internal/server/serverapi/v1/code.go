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

	servercode.LoginPasswordWrong:    "密码错误",
	servercode.LoginFailedReachLimit: "登录失败超过5次，请于5分钟后重试",

	servercode.AuthorizeInvalidTicket: "ticket已失效，请重新登录",
	servercode.AuthorizeInvalidToken:  "token已失效，请重新登录",
	servercode.AuthorizeFailed:        "认证失败，请重新登录",

	servercode.UserAlreadyExists: "用户已存在",
	servercode.UserNotFound:      "用户不存在",
}
