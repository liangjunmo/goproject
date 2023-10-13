package rediskey

import (
	"fmt"
)

func LoginFailedCount(username string) string {
	return fmt.Sprintf("goproject-login-failed-count-%s", username)
}

func LoginTicket(ticket string) string {
	return fmt.Sprintf("goproject-login-ticket-%s", ticket)
}

func MutexCreateUserCenterUser(username string) string {
	return fmt.Sprintf("goproject-mutex-create-usercenter-user-%s", username)
}

func MutexCreateUser(uid uint32) string {
	return fmt.Sprintf("goproject-mutex-create-user-%d", uid)
}
