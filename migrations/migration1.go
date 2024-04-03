package migrations 

import (
	"gorm.io/gorm"
)

type Migration1617448756 struct {}

func (Migration1617448756) Migrate(db *gorm.DB) error {
	err := db.Exec(`
		DO $$ BEGIN 
		IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'access_type')	THEN
			CREATE TYPE access_type AS ENUM ('accepted', 'pending');
		END IF;
		END $$;
	`).Error
	if err != nil {
		return err
	}
	return nil
}