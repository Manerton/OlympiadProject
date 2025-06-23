package school_service

import (
	"context"
	"fmt"
	"log/slog"
	school_dto "main/internal/dto/school"
	"main/internal/lib/liblogger"
	"main/internal/lib/mapper/school_mapper"
	"main/internal/models/school"
	"main/internal/storage/orm"

	"github.com/google/uuid"
)

type SchoolRepository interface {
	GetById(ctx context.Context, orm orm.ORM, id uuid.UUID) (school.School, error)
	GetAll(ctx context.Context, orm orm.ORM, offset, limit *int) ([]school.School, error)

	Create(ctx context.Context, orm orm.ORM, school school.School) (uuid.UUID, error)
	Update(ctx context.Context, orm orm.ORM, school school.School) error
	Delete(ctx context.Context, orm orm.ORM, id uuid.UUID) error
}

type SchoolService struct {
	db               orm.ORM
	log              *slog.Logger
	schoolRepository SchoolRepository
}

func New(log *slog.Logger, orm orm.ORM, schoolRepository SchoolRepository) *SchoolService {
	return &SchoolService{
		log:              log,
		db:               orm,
		schoolRepository: schoolRepository,
	}
}

func (s *SchoolService) GetById(ctx context.Context, id string) (school_dto.SchoolResponeDTO, error) {
	const op = "services.school_service.GetById"
	const errMsg = "failed get school"

	log := s.log.With(
		slog.String("op", op),
	)

	uid, err := uuid.Parse(id)
	if err != nil {
		log.Error("failed parse id", liblogger.Err(err))
		return school_dto.SchoolResponeDTO{}, fmt.Errorf("%s", errMsg)
	}

	schoolModel, err := s.schoolRepository.GetById(ctx, s.db, uid)
	if err != nil {
		log.Error("failed get school", liblogger.Err(err))
		return school_dto.SchoolResponeDTO{}, fmt.Errorf("%s", errMsg)
	}

	return school_mapper.FromModelToDTO(schoolModel), nil
}

func (s *SchoolService) GetAll(ctx context.Context, page, limit *int) ([]school_dto.SchoolResponeDTO, error) {
	const op = "services.school_service.GetAll"
	const errMsg = "failed get all schools"

	log := s.log.With(
		slog.String("op", op),
	)

	offset := new(int)
	if page != nil && limit != nil {
		*offset = (*page - 1) * (*limit)
	}

	schoolsResult, err := s.schoolRepository.GetAll(ctx, s.db, offset, limit)
	if err != nil {
		log.Error("failed get all schools", liblogger.Err(err))
		return nil, fmt.Errorf("%s", errMsg)
	}

	schoolDTO := make([]school_dto.SchoolResponeDTO, 0, len(schoolsResult))
	for _, schoolRes := range schoolsResult {
		schoolDTO = append(schoolDTO, school_mapper.FromModelToDTO(schoolRes))
	}
	return schoolDTO, nil

}

func (s *SchoolService) Create(ctx context.Context, schoolDTO school_dto.CreateSchoolRequestDTO) error {
	const op = "services.school_service.Create"
	const errMsg = "failed create school"

	log := s.log.With(
		slog.String("op", op),
	)

	schoolModel := school_mapper.FromCreateDTOToModel(schoolDTO)
	_, err := s.schoolRepository.Create(ctx, s.db, schoolModel)
	if err != nil {
		log.Error("failed create school", liblogger.Err(err))
		return fmt.Errorf("%s", errMsg)
	}

	return nil
}

func (s *SchoolService) Update(ctx context.Context, id string, schoolDTO school_dto.UpdateSchoolRequestDTO) error {
	const op = "services.school_service.Update"
	const errMsg = "failed update school"

	log := s.log.With(
		slog.String("op", op),
	)

	uid, err := uuid.Parse(id)
	if err != nil {
		log.Error("failed parse id", liblogger.Err(err))
		return fmt.Errorf("failed parse id")
	}

	schoolModel := school_mapper.FromUpdateDTOToModel(schoolDTO, uid)
	err = s.schoolRepository.Update(ctx, s.db, schoolModel)
	if err != nil {
		log.Error("failed update school", liblogger.Err(err))
		return fmt.Errorf("%s", errMsg)
	}

	return nil
}

func (s *SchoolService) Delete(ctx context.Context, id string) error {
	const op = "services.school_service.Delete"
	const errMsg = "failed delete school"

	log := s.log.With(
		slog.String("op", op),
	)

	uid, err := uuid.Parse(id)
	if err != nil {
		log.Error("failed parse id", liblogger.Err(err))
		return fmt.Errorf("%s", errMsg)
	}

	err = s.schoolRepository.Delete(ctx, s.db, uid)
	if err != nil {
		log.Error("failed delete school", liblogger.Err(err))
		return fmt.Errorf("%s", errMsg)
	}

	return nil
}
