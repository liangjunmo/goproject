package servercode

import "github.com/liangjunmo/gocode"

const (
	OK                  gocode.Code = "OK"
	Unknown             gocode.Code = "Unknown"
	Timeout             gocode.Code = "Timeout"
	NotFound            gocode.Code = "NotFound"
	InvalidRequest      gocode.Code = "InvalidRequest"
	InternalServerError gocode.Code = "InternalServerError"

	LoginPasswordWrong    gocode.Code = "LoginPasswordWrong"
	LoginFailedReachLimit gocode.Code = "LoginFailedReachLimit"

	AuthorizeInvalidTicket gocode.Code = "AuthorizeInvalidTicket"
	AuthorizeInvalidToken  gocode.Code = "AuthorizeInvalidToken"
	AuthorizeFailed        gocode.Code = "AuthorizeFailed"

	UserAlreadyExists gocode.Code = "UserAlreadyExists"
	UserNotFound      gocode.Code = "UserNotFound"
)
