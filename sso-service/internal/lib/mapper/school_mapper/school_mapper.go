package school_mapper

import (
	school_dto "main/internal/dto/school"
	"main/internal/models/school"

	"github.com/google/uuid"
)

func FromModelToDTO(school school.School) school_dto.SchoolResponseDTO {
	return school_dto.SchoolResponseDTO{
		ID:         school.ID.String(),
		Name:       school.Name,
		DistrictID: school.DistrictID.String(),
	}
}

func FromCreateDTOToModel(schoolDto school_dto.CreateSchoolRequestDTO) (school.School, error) {

	districtID, err := uuid.Parse(schoolDto.DistrictID)
	if err != nil {
		return school.School{}, nil
	}

	return school.School{
		Name:       schoolDto.Name,
		DistrictID: districtID,
	}, nil
}

func FromUpdateDTOToModel(schoolDto school_dto.UpdateSchoolRequestDTO, uid uuid.UUID) (school.School, error) {
	districtID, err := uuid.Parse(*schoolDto.DistrictID)
	if err != nil {
		return school.School{}, nil
	}

	return school.School{
		ID:         uid,
		Name:       *schoolDto.Name,
		DistrictID: districtID,
	}, nil
}
