package postgresql

import (
	"fmt"
	"strings"
	"time"

	"OlimpiadPortal/ApplicationService/internal/models"

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

	// Add fixed data after migration
	seedData(db)

	return db, nil
}

func seedData(db *gorm.DB) error {
	// Тестовые данные для заявок
	applications := []models.Application{
		{
			UserID:        1,
			EventID:       101,
			EventName:     "Олимпиада по информатике",     // ВРЕМЕННО
			EventLocation: "Школа №122",                   // ВРЕМЕННО
			EventDate:     time.Now().Add(48 * time.Hour), // ВРЕМЕННО
			Status:        nil,
			SubmittedAt:   time.Now().Add(-48 * time.Hour), // 2 дня назад
			UpdatedAt:     time.Now(),
		},
		{
			UserID:        1,
			EventID:       102,
			EventName:     "Олимпиада по математике",       // ВРЕМЕННО
			EventLocation: "Школа №123",                    // ВРЕМЕННО
			EventDate:     time.Now().Add(72 * time.Hour),  // ВРЕМЕННО
			Status:        boolPtr(true),                   // Одобрено
			SubmittedAt:   time.Now().Add(-24 * time.Hour), // 1 день назад
			UpdatedAt:     time.Now(),
		},
		{
			UserID:        1,
			EventID:       103,
			EventName:     "Олимпиада по Русскому языку",   // ВРЕМЕННО
			EventLocation: "Школа №64",                     // ВРЕМЕННО
			EventDate:     time.Now().Add(24 * time.Hour),  // ВРЕМЕННО
			Status:        boolPtr(false),                  // Отклонено
			SubmittedAt:   time.Now().Add(-72 * time.Hour), // 3 дня назад
			UpdatedAt:     time.Now(),
		},
	}

	/*
		for _, application := range applications {
			if err := db.Create(&application).Error; err != nil {
				return err
			}
		} */

	// Создаем записи в базе данных
	for _, application := range applications {
		var existingApplication models.Application
		err := db.Where("user_id = ? AND event_id = ?", application.UserID, application.EventID).First(&existingApplication).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				// Record does not exist, create it
				if err := db.Create(&application).Error; err != nil {
					return err
				}
			} else {
				// Some other error occurred
				return err
			}
		} else {
			fmt.Printf("Skipping duplicate application: UserID=%d, EventID=%d\n", application.UserID, application.EventID)
		}
	}

	return nil
}

// Вспомогательная функция для указателя на bool
func boolPtr(b bool) *bool {
	return &b
}
