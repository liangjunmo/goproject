package helper

import (
	"github.com/liangjunmo/goproject/internal/types"
)

func UserCenterUserToMap(users []types.UserCenterUser) map[uint32]types.UserCenterUser {
	userMap := make(map[uint32]types.UserCenterUser)

	for _, user := range users {
		userMap[user.UID] = user
	}

	return userMap
}

func FetchUserCenterUserUids(users []types.UserCenterUser) []uint32 {
	uids := make([]uint32, 0, len(users))

	for _, user := range users {
		uids = append(uids, user.UID)
	}

	return uids
}

func UserToMap(users []types.User) map[uint32]types.User {
	userMap := make(map[uint32]types.User)

	for _, user := range users {
		userMap[user.UID] = user
	}

	return userMap
}

func FetchUserUids(users []types.User) []uint32 {
	uids := make([]uint32, 0, len(users))

	for _, user := range users {
		uids = append(uids, user.UID)
	}

	return uids
}
