package ApplicationRepository

import (
	"context"
	"fmt"
	models "main/internal/models"

	"main/internal/storage/orm"

	"github.com/google/uuid"
)

// Структура репозитория
type ApplicationRepository struct{}

// -1. Получить заявки по фильтру
func (r *ApplicationRepository) GetAllByFilter(ctx context.Context, orm orm.ORM, filter models.Application, offset, limit *int, order *string) ([]models.Application, error) {
	const op = "repositories.UserRepository.GetByFilter"

	AppResult := []models.Application{}

	if err := orm.Find(ctx, models.Application{}, nil, offset, limit, order, &AppResult, filter); err != nil {
		return []models.Application{}, fmt.Errorf("%s: %w", op, err)
	}

	return AppResult, nil
}

// 0. Получить количество всех заявок
func (r *ApplicationRepository) GetCount(ctx context.Context, orm orm.ORM) (int64, error) {
	const op = "repositories.UserRepository.GetCount"

	var countResult int64 = 0
	err := orm.Count(ctx, models.Application{}, &countResult, nil)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return countResult, nil
}

// 1. Создание новой заявки
func (r *ApplicationRepository) Create(ctx context.Context, orm orm.ORM, application models.Application) (uuid.UUID, error) {
	const op = "repositories.ApplicationRepository.CreateApplication"
	err := orm.Create(ctx, &application)
	if err != nil {
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}
	return application.ID, nil
}

// 2. Получение заявки по ID
func (r *ApplicationRepository) GetByID(ctx context.Context, orm orm.ORM, id uuid.UUID) (models.Application, error) {
	const op = "repositories.application_repository.GetApplicationByID"

	application := models.Application{}
	err := orm.First(ctx, models.Application{}, nil, &application, models.Application{ID: id})
	if err != nil {
		return models.Application{}, fmt.Errorf("%s: %w", op, err)
	}
	return application, nil
}

// 3. Получение всех заявок
func (r *ApplicationRepository) GetAllApplications(ctx context.Context, orm orm.ORM, offset *int, limit *int) ([]models.Application, error) {
	const op = "repositories.application_repository.GetAllApplications"

	applications := []models.Application{}
	err := orm.Find(
		ctx,
		models.Application{}, // Модель
		nil,                  // Поля (nil - выбираем все)
		offset,               // Offset (nil - без offset)
		limit,                // Limit (nil - без лимита)
		nil,                  // Order (nil - без сортировки)
		&applications,        // Куда записать результат
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return applications, nil
}

// 4. Обновление статуса заявки
func (r *ApplicationRepository) UpdateApplication(ctx context.Context, orm orm.ORM, application models.Application) error {
	const op = "repositories.application_repository.UpdateApplicationStatus"
	err := orm.Updates(ctx, nil, &application)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// 5. Удаление заявки по ID
func (r *ApplicationRepository) DeleteApplicationByID(ctx context.Context, orm orm.ORM, id uuid.UUID) error {
	const op = "repositories.application_repository.DeleteApplicationByID"

	application := models.Application{}
	err := orm.Delete(ctx, models.Application{}, nil, &application, models.Application{ID: id})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

// Удаление по условию
func (r *ApplicationRepository) DeleteByFilter(ctx context.Context, orm orm.ORM, model models.Application) error {
	const op = "repositories.application_repository.DeleteByFilter"

	err := orm.Delete(ctx, models.Application{}, nil, nil, model)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

// 6. Получение всех заявок по ID пользователя
func (r *ApplicationRepository) GetApplicationsByUserID(ctx context.Context, orm orm.ORM, userID uuid.UUID, offset *int, limit *int) ([]models.Application, error) {
	const op = "repositories.application_repository.GetApplicationsByUserID"

	var applications []models.Application

	// Условие для поиска (WHERE user_id = ?)
	condition := models.Application{UserID: userID}

	err := orm.Find(
		ctx,
		models.Application{}, // Модель
		nil,                  // Поля (nil - выбираем все)
		offset,               // Offset (nil - без offset)
		limit,                // Limit (nil - без лимита)
		nil,                  // Order (nil - без сортировки)
		&applications,        // Куда записать результат
		condition,            // Условия WHERE
	)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return applications, nil
}

// 7. Получение всех заявок по ID события
func (r *ApplicationRepository) GetApplicationsByEventID(ctx context.Context, orm orm.ORM, eventID uuid.UUID, offset *int, limit *int) ([]models.Application, error) {
	const op = "repositories.application_repository.GetApplicationsByEventID"

	var applications []models.Application

	// Условие для поиска (WHERE user_id = ?)
	condition := models.Application{EventID: eventID}

	err := orm.Find(
		ctx,
		models.Application{}, // Модель
		nil,                  // Поля (nil - выбираем все)
		offset,               // Offset (nil - без offset)
		limit,                // Limit (nil - без лимита)
		nil,                  // Order (nil - без сортировки)
		&applications,        // Куда записать результат
		condition,            // Условия WHERE
	)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return applications, nil
}
