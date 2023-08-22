package servercode

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

	LoginPasswordWrong:    "密码错误",
	LoginFailedReachLimit: "登录失败超过5次，请于5分钟后重试",

	AuthorizeInvalidTicket: "ticket已失效，请重新登录",
	AuthorizeInvalidToken:  "token已失效，请重新登录",
	AuthorizeFailed:        "认证失败，请重新登录",

	UserAlreadyExists: "用户已存在",
	UserNotFound:      "用户不存在",
}
