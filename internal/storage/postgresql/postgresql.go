package postgresql

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPosgreSQL(connectStr string) (*gorm.DB, error) {

	const op = "storage.postgresql.NewPostreSQL"

	// try to connect posgresql
	db, err := gorm.Open(postgres.Open(connectStr), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// migration models
	// TODO: need model
	// err = db.AutoMigrate(&event.Event{})
	// if err != nil {
	// 	return nil, fmt.Errorf("%s: %w", op, err)
	// }

	return db, nil
}
