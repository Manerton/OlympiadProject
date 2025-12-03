package links_service

import (
	"context"
	"log/slog"
	link_dto "main/internal/dto/link"
	"main/internal/lib/errs"
	"main/internal/lib/liblogger"
	"main/internal/models/district"
	"main/internal/models/school"
	"main/internal/storage/orm"
	"strconv"

	"github.com/google/uuid"
)

type SchoolRepositoryInterface interface {
	GetAllByDistrict(ctx context.Context, orm orm.ORM, districtId uuid.UUID) ([]school.School, error)
}

type DistrictRepositoryInterface interface {
	GetAllByRegion(ctx context.Context, orm orm.ORM, region int) ([]district.District, error)
}

type LinksService struct {
	log                *slog.Logger
	db                 orm.ORM
	schoolRepository   SchoolRepositoryInterface
	districtRepository DistrictRepositoryInterface
}

func New(log *slog.Logger, orm orm.ORM, schoolRepository SchoolRepositoryInterface, districtRepository DistrictRepositoryInterface) *LinksService {
	return &LinksService{
		log:                log,
		db:                 orm,
		schoolRepository:   schoolRepository,
		districtRepository: districtRepository,
	}
}

func (s *LinksService) GetLinks(ctx context.Context, region string) ([]link_dto.LinkDTO, error) {
	const op = "services.LinksService.GetLinks"

	log := s.log.With(
		slog.String("op", op),
	)

	regionInt, err := strconv.Atoi(region)
	if err != nil {
		log.Error("invalid region parameter", slog.String("region", region), liblogger.Err(err))
		return nil, errs.ErrBadRequest.Wrap("invalid region parameter")
	}

	allDistricts, err := s.districtRepository.GetAllByRegion(ctx, s.db, regionInt)
	if err != nil {
		log.Error("failed get districts by region", slog.String("region", region), liblogger.Err(err))
		return nil, errs.ErrInternalError.Wrap("failed get districts by region")
	}

	_ = allDistricts

	return nil, nil
}
