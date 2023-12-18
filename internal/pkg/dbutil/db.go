package dbutil

import (
	"gorm.io/gorm"
)

func TruncateTable(db *gorm.DB, models []interface{}) error {
	stmt := &gorm.Statement{DB: db}

	for _, model := range models {
		err := stmt.Parse(model)
		if err != nil {
			return err
		}

		err = db.Exec("truncate table " + stmt.Schema.Table + ";").Error
		if err != nil {
			return err
		}
	}

	return nil
}
