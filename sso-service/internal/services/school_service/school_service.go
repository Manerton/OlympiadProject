package school_service

import (
	"context"
	"errors"
	"log/slog"
	school_dto "main/internal/dto/school"
	"main/internal/lib/errs"
	"main/internal/lib/liblogger"
	"main/internal/lib/mapper/school_mapper"
	"main/internal/models/school"
	"main/internal/storage/orm"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SchoolRepository interface {
	GetById(ctx context.Context, orm orm.ORM, id uuid.UUID) (school.School, error)
	GetAll(ctx context.Context, orm orm.ORM, offset, limit *int) ([]school.School, error)

	GetCount(ctx context.Context, orm orm.ORM) (int64, error)

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

func (s *SchoolService) GetCount(ctx context.Context) (int64, error) {
	const op = "services.school_service.GetCount"

	log := s.log.With(
		slog.String("op", op),
	)

	schoolCount, err := s.schoolRepository.GetCount(ctx, s.db)
	if err != nil {
		log.Error("failed get count schools", liblogger.Err(err))
		return 0, errs.ErrInternalError.Wrap("failed get cound schools")
	}

	return schoolCount, nil
}

func (s *SchoolService) GetById(ctx context.Context, id string) (school_dto.SchoolResponseDTO, error) {
	const op = "services.school_service.GetById"

	log := s.log.With(
		slog.String("op", op),
	)

	uid, err := uuid.Parse(id)
	if err != nil {
		log.Error("failed parse id", slog.String("id", id), liblogger.Err(err))
		return school_dto.SchoolResponseDTO{}, errs.ErrBadRequest.Wrap("failed parse uuid")
	}

	schoolModel, err := s.schoolRepository.GetById(ctx, s.db, uid)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		log.Warn("school not found", liblogger.Err(err))
		return school_dto.SchoolResponseDTO{}, errs.ErrSchoolNotFound
	}

	if err != nil {
		log.Error("failed get school", liblogger.Err(err))
		return school_dto.SchoolResponseDTO{}, errs.ErrInternalError.Wrap("failed find school")
	}

	return school_mapper.FromModelToDTO(schoolModel), nil
}

func (s *SchoolService) GetAll(ctx context.Context, page, limit *int) ([]school_dto.SchoolResponseDTO, error) {
	const op = "services.school_service.GetAll"

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
		return nil, errs.ErrInternalError.Wrap("failed get all schools")
	}

	schoolDTO := make([]school_dto.SchoolResponseDTO, 0, len(schoolsResult))
	for _, schoolRes := range schoolsResult {
		schoolDTO = append(schoolDTO, school_mapper.FromModelToDTO(schoolRes))
	}
	return schoolDTO, nil

}

func (s *SchoolService) Create(ctx context.Context, schoolDTO school_dto.CreateSchoolRequestDTO) (uuid.UUID, error) {
	const op = "services.school_service.Create"

	log := s.log.With(
		slog.String("op", op),
	)

	schoolModel, err := school_mapper.FromCreateDTOToModel(schoolDTO)
	if err != nil {
		log.Error("failed convert dto to model", slog.Any("dto", schoolDTO), liblogger.Err(err))
		return uuid.Nil, errs.ErrBadRequest.Wrap("failed convert dto to model")
	}
	uid, err := s.schoolRepository.Create(ctx, s.db, schoolModel)
	if err != nil {
		log.Error("failed create school", liblogger.Err(err))
		return uuid.Nil, errs.ErrInternalError.Wrap("failed create school")
	}

	return uid, nil
}

func (s *SchoolService) Update(ctx context.Context, id string, schoolDTO school_dto.UpdateSchoolRequestDTO) error {
	const op = "services.school_service.Update"

	log := s.log.With(
		slog.String("op", op),
	)

	uid, err := uuid.Parse(id)
	if err != nil {
		log.Error("failed parse id", slog.String("id", id), liblogger.Err(err))
		return errs.ErrBadRequest.Wrap("failed parse uuid")
	}

	schoolModel, err := school_mapper.FromUpdateDTOToModel(schoolDTO, uid)
	if err != nil {
		log.Error("failed convert dto to model", slog.Any("dto", schoolDTO), liblogger.Err(err))
		return errs.ErrBadRequest.Wrap("failed convert dto to model")
	}

	err = s.schoolRepository.Update(ctx, s.db, schoolModel)
	if err != nil {
		log.Error("failed update school", liblogger.Err(err))
		return errs.ErrInternalError.Wrap("failed update school")
	}

	return nil
}

func (s *SchoolService) Delete(ctx context.Context, id string) error {
	const op = "services.school_service.Delete"

	log := s.log.With(
		slog.String("op", op),
	)

	uid, err := uuid.Parse(id)
	if err != nil {
		log.Error("failed parse id", slog.String("id", id), liblogger.Err(err))
		return errs.ErrBadRequest.Wrap("failed parse uuid")
	}

	err = s.schoolRepository.Delete(ctx, s.db, uid)
	if err != nil {
		log.Error("failed delete school", liblogger.Err(err))
		return errs.ErrInternalError.Wrap("failed delete school")
	}

	return nil
}
