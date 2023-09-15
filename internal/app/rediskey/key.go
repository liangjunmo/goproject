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

func MutexCreateUser(username string) string {
	return fmt.Sprintf("goproject-mutex-create-user-%s", username)
}