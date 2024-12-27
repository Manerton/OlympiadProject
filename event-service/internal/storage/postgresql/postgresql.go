package postgresql

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPosgreSQL(connectStr string) (*gorm.DB, error) {

	const op = "storage.postgresql.NewPostreSQL"

	// try to connect posgresql
	db, err := gorm.Open(postgres.Open(connectStr), &gorm.Config{})
	if err != nil {
		log.Fatalf("%s: %v", op, err)
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// // create enum (types Events)
	// create_enum_type_str := fmt.Sprintf(
	// 	`DO $$ BEGIN IF NOT EXISTS
	// 	(SELECT 1 FROM pg_type WHERE typname = 'event_type')
	// 	THEN CREATE TYPE event_type AS ENUM
	// 	('%s', '%s', '%s', '%s', '%s');
	// 	END IF; END $$;`,
	// 	event.RegionalStage,
	// 	event.Olympiad,
	// 	event.Stage,
	// 	event.ViewWorks,
	// 	event.Appeal,
	// )
	// db.Exec(create_enum_type_str)

	// // migration models
	// err = db.AutoMigrate(&event.Event{})
	// if err != nil {
	// 	log.Fatalf("%s: %v", op, err)
	// 	return nil, fmt.Errorf("%s: %w", op, err)
	// }

	return db, nil
}
