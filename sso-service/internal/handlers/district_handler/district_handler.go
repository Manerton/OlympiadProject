package district_handler

import (
	"context"
	district_dto "main/internal/dto/district"
)

type DistrictServiceInterface interface {
	GetAllByRegion(ctx context.Context, region string) ([]district_dto.DistrictDTOResponse, error)
	GetByID(ctx context.Context, id string) (district_dto.DistrictDTOResponse, error)
}

type DistrictHandler struct {
	service DistrictServiceInterface
}

func New(service DistrictServiceInterface) *DistrictHandler {
	return &DistrictHandler{
		service: service,
	}
}
