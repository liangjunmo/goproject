package userservice

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"github.com/liangjunmo/goproject/internal/goproject/testutil"
	"github.com/liangjunmo/goproject/internal/pkg/dbutil"
	"github.com/liangjunmo/goproject/internal/pkg/pagination"
)

func TestDefaultRepository(t *testing.T) {
	db := testutil.InitDB()

	beforeTest := func(t *testing.T) {
		err := dbutil.TruncateTable(db, []interface{}{&User{}})
		require.Nil(t, err)
	}

	t.Run("Commit", func(t *testing.T) {
		beforeTest(t)

		tx := db.Begin()

		tx.Create(&User{UID: 1})

		repository := newDefaultRepository(tx)

		err := repository.Commit()
		require.Nil(t, err)

		err = db.Take(&User{}, 1).Error
		require.Nil(t, err)
	})

	t.Run("Rollback", func(t *testing.T) {
		beforeTest(t)

		tx := db.Begin()

		tx.Create(&User{UID: 1})

		repository := newDefaultRepository(tx)

		err := repository.Rollback()
		require.Nil(t, err)

		err = db.Take(&User{}, 1).Error
		require.ErrorIs(t, err, gorm.ErrRecordNotFound)
	})

	t.Run("List", func(t *testing.T) {
		beforeTest(t)

		db.Create(&User{UID: 1})

		repository := newDefaultRepository(db)

		_, users, err := repository.List(context.Background(), criteria{
			Request: pagination.DefaultRequest{
				Page:     1,
				Capacity: 10,
			},
			uids: []uint32{1},
		})
		require.Nil(t, err)
		require.Equal(t, uint32(1), users[0].UID)
	})

	t.Run("Search", func(t *testing.T) {
		beforeTest(t)

		db.Create(&User{UID: 1})

		repository := newDefaultRepository(db)

		users, err := repository.Search(context.Background(), criteria{
			uids: []uint32{1},
		})
		require.Nil(t, err)
		require.Equal(t, uint32(1), users[1].UID)
	})

	t.Run("Get", func(t *testing.T) {
		beforeTest(t)

		db.Create(&User{UID: 1})

		repository := newDefaultRepository(db)

		user, exist, err := repository.Get(context.Background(), 1)
		require.Nil(t, err)
		require.True(t, exist)
		require.Equal(t, uint32(1), user.UID)
	})

	t.Run("Create", func(t *testing.T) {
		beforeTest(t)

		repository := newDefaultRepository(db)

		err := repository.Create(context.Background(), &User{UID: 1})
		require.Nil(t, err)

		var user User

		err = db.Take(&user, 1).Error
		require.Nil(t, err)
		require.Equal(t, uint32(1), user.UID)
	})

	{
		db, _ := db.DB()
		db.Close()
	}
}