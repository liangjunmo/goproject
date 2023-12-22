package codes

import "github.com/liangjunmo/gocode"

const (
	OK                  gocode.Code = "OK"
	Unknown             gocode.Code = "Unknown"
	Timeout             gocode.Code = "Timeout"
	NotFound            gocode.Code = "NotFound"
	InvalidRequest      gocode.Code = "InvalidRequest"
	InternalServerError gocode.Code = "InternalServerError"

	LoginFailedPasswordWrong gocode.Code = "LoginFailedPasswordWrong"
	LoginFailedReachLimit    gocode.Code = "LoginFailedReachLimit"

	AuthorizeFailed              gocode.Code = "AuthorizeFailed"
	AuthorizeFailedInvalidTicket gocode.Code = "AuthorizeFailedInvalidTicket"
	AuthorizeFailedInvalidToken  gocode.Code = "AuthorizeFailedInvalidToken"

	UserNotFound      gocode.Code = "UserNotFound"
	UserAlreadyExists gocode.Code = "UserAlreadyExists"
)
