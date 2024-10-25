package subject_service

import (
	"fmt"
	"log/slog"
	"main/internal/dto/subject_dto"
	"main/internal/models/subject"
	"main/internal/repositories/subject_repository"

	"gorm.io/gorm"
)

type SubjectService struct {
	log        *slog.Logger
	db         *gorm.DB
	repository *subject_repository.SubjectRepository
}

func NewSubjectService(log *slog.Logger, db *gorm.DB, sr *subject_repository.SubjectRepository) *SubjectService {
	return &SubjectService{
		log:        log,
		db:         db,
		repository: sr,
	}
}

func (s *SubjectService) GetSubjectByID(id uint) (subject_dto.SubjectDTO, error) {
	const op = "service.subject_service.GetEventByID"
	// s.log.With(slog.String("op", op))
	subject, err := s.repository.GetSubjectByID(s.db, id)
	if err != nil {
		// s.log.Error("failed to get subject", liblogger.Err(err))
		return subject_dto.SubjectDTO{}, fmt.Errorf("%s: %w", op, err)
	}
	return ConvertSubjectToDTO(subject), nil
}

func (s *SubjectService) GetAllSubjects() ([]subject_dto.SubjectDTO, error) {
	const op = "service.subject_service.GetAllSubjects"
	subjects, err := s.repository.GetAllSubjects(s.db)
	if err != nil {
		return nil, fmt.Errorf("%s, %w", op, err)
	}
	return ConverManySubjectsToDTO(subjects), nil
}

func (s *SubjectService) CreateSubject(subject_dto subject_dto.SubjectDTO) (uint, error) {
	const op = "service.subject_service.CreateSubject"
	subject := ConvertDTOtoSubject(subject_dto)

	id, err := s.repository.CreateSubject(s.db, subject)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}

func (s *SubjectService) UpdateSubject(subject_dto subject_dto.SubjectDTO) (uint, error) {
	const op = "service.subject_service.UpdateSubject"
	subject := ConvertDTOtoSubject(subject_dto)

	id, err := s.repository.UpdateSubject(s.db, subject)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}

func (s *SubjectService) DeleteSubject(id uint) error {
	const op = "service.subject_service.UpdateSubject"
	err := s.repository.DeleteSubject(s.db, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func ConvertDTOtoSubject(subject_dto subject_dto.SubjectDTO) subject.Subject {
	return subject.Subject{
		ID:   subject_dto.ID,
		Name: subject_dto.Name,
	}
}

func ConvertSubjectToDTO(subject subject.Subject) subject_dto.SubjectDTO {
	return subject_dto.SubjectDTO{
		ID:   subject.ID,
		Name: subject.Name,
	}
}

func ConverManySubjectsToDTO(subjects []subject.Subject) []subject_dto.SubjectDTO {
	var subjects_dto []subject_dto.SubjectDTO
	for _, subject := range subjects {
		subjects_dto = append(subjects_dto, ConvertSubjectToDTO(subject))
	}
	return subjects_dto
}
