package ApplicationService

import (
	ApplicationDto "OlimpiadPortal/ApplicationService/internal/dto"
	"OlimpiadPortal/ApplicationService/internal/models"
	application_repository "OlimpiadPortal/ApplicationService/internal/repositories"
	"fmt"

	"gorm.io/gorm"
)

type ApplicationService struct {
	db         *gorm.DB
	repository *application_repository.ApplicationRepository
}

func NewApplicationService(db *gorm.DB, repo *application_repository.ApplicationRepository) *ApplicationService {
	return &ApplicationService{
		db:         db,
		repository: repo,
	}
}

// Получение всех заявок
func (s *ApplicationService) GetAllApplications() ([]ApplicationDto.ApplicationResponseDTO, error) {
	const op = "services.application_service.GetAllApplications"
	applications, err := s.repository.GetAllApplications(s.db)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return ConvertManyApplicationsToDTO(applications), nil
}

// Получение заявки по ID
func (s *ApplicationService) GetApplicationByID(id uint) (ApplicationDto.ApplicationResponseDTO, error) {
	const op = "services.application_service.GetApplicationByID"
	application, err := s.repository.GetApplicationByID(s.db, id)
	if err != nil {
		return ApplicationDto.ApplicationResponseDTO{}, fmt.Errorf("%s: %w", op, err)
	}
	return ConvertApplicationToDTO(application), nil
}

// Создание новой заявки
func (s *ApplicationService) CreateApplication(applicationDTO ApplicationDto.CreateApplicationDTO) (uint, error) {
	const op = "services.application_service.CreateApplication"
	application := ConvertDTOtoApplication(applicationDTO)
	if err := s.repository.CreateApplication(s.db, &application); err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return application.ApplicationID, nil
}

// Обновление статуса заявки
func (s *ApplicationService) UpdateApplicationStatus(id uint, statusDTO ApplicationDto.UpdateApplicationStatusDTO) error {
	const op = "services.application_service.UpdateApplicationStatus"
	if err := s.repository.UpdateApplicationStatus(s.db, id, statusDTO.Status); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

// Удаление заявки по ID
func (s *ApplicationService) DeleteApplication(id uint) error {
	const op = "services.application_service.DeleteApplication"
	if err := s.repository.DeleteApplicationByID(s.db, id); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

// Функции для преобразования между DTO и моделью

func ConvertDTOtoApplication(dto ApplicationDto.CreateApplicationDTO) models.Application {
	return models.Application{
		UserID:  dto.UserID,
		EventID: dto.EventID,
	}
}

func ConvertApplicationToDTO(application models.Application) ApplicationDto.ApplicationResponseDTO {
	return ApplicationDto.ApplicationResponseDTO{
		ApplicationID: application.ApplicationID,
		UserID:        application.UserID,
		EventID:       application.EventID,
		Status:        application.Status,
		SubmittedAt:   application.SubmittedAt,
		UpdatedAt:     application.UpdatedAt,
	}
}

func ConvertManyApplicationsToDTO(applications []models.Application) []ApplicationDto.ApplicationResponseDTO {
	var applicationsDTO []ApplicationDto.ApplicationResponseDTO
	for _, application := range applications {
		applicationsDTO = append(applicationsDTO, ConvertApplicationToDTO(application))
	}
	return applicationsDTO
}
