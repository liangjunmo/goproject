package rediskey

import (
	"fmt"
)

func MutexCreateUser(username string) string {
	return fmt.Sprintf("goproject-mutex-create-user-%s", username)
}
