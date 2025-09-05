package postgresql

import (
	"fmt"
	"strings"
	"time"

	"main/internal/models"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
	// Предопределенные UUID пользователей
	userIDs := []string{
		"a08957c4-7911-4289-bfed-42f0466d3e06",
		"af3aa7a9-f824-44e2-a3f6-275f94b04104",
		"e49898de-1cd8-462b-882f-d550904dda9e",
		"d3b41646-4913-467b-860b-7290e52ef477",
		"b90bc433-6210-4836-bafc-837c88a794ea",
		"1871cbfc-c50f-43aa-86a1-619ffb569b25",
		"d5ec1598-6468-4cf4-ae07-64876019b93e",
		"3d297cc8-91b2-4b88-83e2-dc08d6106750",
		"ebe3ad07-69f1-41bb-b080-fc49a4dc16fa",
		"6aec8f4c-9330-4afc-801b-4d4eaa6cf75d",
		"d2c38359-5d05-4789-945d-3af037e27c99",
		"106066d9-c4eb-4aea-bcb5-0f395b97f5eb",
		"81ab90de-f856-4cb9-912b-0804072c3611",
	}

	// Предопределенные UUID олимпиад
	eventIDs := []string{
		"a8edf681-fb9a-458a-b6f0-91035a67dbba",
		"46e160c4-7fe4-4f14-9193-e06ec8b41c34",
		"46e160c4-7fe4-4f14-9193-e06ec8b41c34",
	}

	// Тестовые данные для заявок
	var applications []models.Application

	// Для каждой олимпиады создаем 4 заявки
	for i, eventID := range eventIDs {
		// Получаем 4 уникальных пользователя для этой олимпиады
		startIdx := i * 4
		endIdx := startIdx + 4
		eventUsers := userIDs[startIdx:endIdx]

		// Создаем 4 заявки для текущей олимпиады
		for j, userID := range eventUsers {
			userUUID, _ := uuid.Parse(userID)
			eventUUID, _ := uuid.Parse(eventID)

			// Уникальные параметры для каждой заявки
			status := j + 1 // 1, 2, 3
			if status > 3 {
				status = 3
			}
			reason := j%2 + 1 // чередуем 1 и 2
			code := fmt.Sprintf("%02d_%03d_%02d", 9+j, 100+j*5, 20+i+j)

			applications = append(applications, models.Application{
				UserID:      userUUID,
				EventID:     eventUUID,
				Status:      status,
				Reason:      reason,
				Code:        code,
				SubmittedAt: time.Now().Add(-time.Duration(24*(i*4+j+1)) * time.Hour),
				UpdatedAt:   time.Now(),
			})
		}
	}

	// Создаем заявки в базе данных
	if err := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&applications).Error; err != nil {
		return fmt.Errorf("failed to seed applications: %w", err)
	}

	return nil
}

// Вспомогательная функция для указателя на bool
func boolPtr(b bool) *bool {
	return &b
}
