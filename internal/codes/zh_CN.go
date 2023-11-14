package codes

import (
	"github.com/liangjunmo/gocode"
)

var zhCn = map[gocode.Code]string{
	OK:                  "OK",
	Unknown:             "未知错误",
	Timeout:             "请求超时",
	NotFound:            "资源不存在",
	InvalidRequest:      "请求错误",
	InternalServerError: "服务端错误",

	LoginFailedPasswordWrong: "密码错误",
	LoginFailedReachLimit:    "登录失败次数超过5次，请于5分钟后重试",

	AuthorizeFailed:              "认证失败，请重新登录",
	AuthorizeFailedInvalidTicket: "ticket已失效，请重新登录",
	AuthorizeFailedInvalidToken:  "token已失效，请重新登录",

	UserNotFound:      "用户不存在",
	UserAlreadyExists: "用户已存在",
}
