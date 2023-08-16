package servercode

import "github.com/liangjunmo/gocode"

const (
	OK                  gocode.Code = "OK"
	Unknown             gocode.Code = "Unknown"
	Timeout             gocode.Code = "Timeout"
	NotFound            gocode.Code = "NotFound"
	InvalidRequest      gocode.Code = "InvalidRequest"
	InternalServerError gocode.Code = "InternalServerError"

	UserAlreadyExists gocode.Code = "UserAlreadyExists"
)
