package postgresql

import (
	"fmt"
	"strings"

	"main/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgreSQL(connectStr string) (*gorm.DB, error) {

	const op = "storage.postgresql.NewPostreSQL"

	// Split the connection string into parts
	params := strings.Fields(connectStr)

	// Initialize variables to hold the modified connection strings
	var dbName string
	var baseDSNParts []string

	// Loop through the parameters to find dbname and construct the base DSN
	for _, param := range params {
		if strings.HasPrefix(param, "dbname=") {
			// Extract the database name
			dbName = strings.TrimPrefix(param, "dbname=")
		} else {
			// Append other parameters to baseDSNParts
			baseDSNParts = append(baseDSNParts, param)
		}
	}

	// Create the base connection string without dbname
	baseDSN := strings.Join(baseDSNParts, " ")

	// Connect to PostgreSQL server without specifying a database
	db, err := gorm.Open(postgres.Open(baseDSN), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// Check if the database exists; if not, create it
	db.Exec(fmt.Sprintf("CREATE DATABASE \"%s\"", dbName))

	// try to connect posgresql
	db, err = gorm.Open(postgres.Open(connectStr), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// migration models
	err = db.AutoMigrate(&models.Application{})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return db, nil
}
