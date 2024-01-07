package userserviceimpl

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"github.com/liangjunmo/goproject/internal/pkg/dbutil"
	"github.com/liangjunmo/goproject/internal/testutil"
	"github.com/liangjunmo/goproject/internal/usercenter/service/userservice"
)

func TestDefaultRepository(t *testing.T) {
	db := testutil.InitDB()
	defer func() {
		db, _ := db.DB()
		db.Close()
	}()

	var (
		repository *defaultRepository
		ctx        context.Context
	)

	beforeTest := func(t *testing.T) {
		err := dbutil.TruncateTable(db, []interface{}{&userservice.User{}})
		require.Nil(t, err)

		repository = newDefaultRepository(db)

		ctx = context.Background()
	}

	t.Run("Commit", func(t *testing.T) {
		beforeTest(t)

		repository, err := repository.Begin()
		require.Nil(t, err)

		err = repository.Create(ctx, &userservice.User{
			UID:      1,
			Username: "user",
			Password: "pass",
		})
		require.Nil(t, err)

		err = repository.Commit()
		require.Nil(t, err)

		err = db.Take(&userservice.User{}, 1).Error
		require.Nil(t, err)
	})

	t.Run("Rollback", func(t *testing.T) {
		beforeTest(t)

		repository, err := repository.Begin()
		require.Nil(t, err)

		err = repository.Create(ctx, &userservice.User{
			UID:      1,
			Username: "user",
			Password: "pass",
		})
		require.Nil(t, err)

		err = repository.Rollback()
		require.Nil(t, err)

		err = db.Take(&userservice.User{}, 1).Error
		require.ErrorIs(t, err, gorm.ErrRecordNotFound)
	})

	t.Run("Search", func(t *testing.T) {
		beforeTest(t)

		db.Create(&userservice.User{
			UID:      1,
			Username: "user",
		})

		users, err := repository.Search(ctx, criteria{
			uids:     []uint32{1},
			username: "user",
		})
		require.Nil(t, err)
		require.Equal(t, uint32(1), users[1].UID)
		require.Equal(t, "user", users[1].Username)
	})

	t.Run("Get", func(t *testing.T) {
		beforeTest(t)

		db.Create(&userservice.User{UID: 1})

		user, exist, err := repository.Get(ctx, 1)
		require.Nil(t, err)
		require.True(t, exist)
		require.Equal(t, uint32(1), user.UID)
	})

	t.Run("GetByUsername", func(t *testing.T) {
		beforeTest(t)

		db.Create(&userservice.User{Username: "user"})

		user, exist, err := repository.GetByUsername(ctx, "user")
		require.Nil(t, err)
		require.True(t, exist)
		require.Equal(t, "user", user.Username)
	})

	t.Run("Create", func(t *testing.T) {
		beforeTest(t)

		err := repository.Create(ctx, &userservice.User{
			UID:      1,
			Username: "user",
			Password: "pass",
		})
		require.Nil(t, err)

		var user userservice.User

		err = db.Take(&user, 1).Error
		require.Nil(t, err)
		require.Equal(t, uint32(1), user.UID)
		require.Equal(t, "user", user.Username)
		require.Equal(t, "pass", user.Password)
	})
}
