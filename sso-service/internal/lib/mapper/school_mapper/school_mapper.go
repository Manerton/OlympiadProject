package school_mapper

import (
	school_dto "main/internal/dto/school"
	"main/internal/models/school"
	"strconv"

	"github.com/google/uuid"
)

func FromModelToDTO(school school.School) school_dto.SchoolResponeDTO {
	return school_dto.SchoolResponeDTO{
		ID:     school.ID.String(),
		Name:   school.Name,
		Region: school.Region,
	}
}

func FromCreateDTOToModel(schoolDto school_dto.CreateSchoolRequestDTO) (school.School, error) {

	regionNum, err := strconv.Atoi(schoolDto.Region)
	if err != nil {
		return school.School{}, nil
	}

	return school.School{
		Name:   schoolDto.Name,
		Region: regionNum,
	}, nil
}

func FromUpdateDTOToModel(schoolDto school_dto.UpdateSchoolRequestDTO, uid uuid.UUID) (school.School, error) {
	if schoolDto.Region == nil {
		return school.School{}, nil
	}

	regionNum, err := strconv.Atoi(*schoolDto.Region)
	if err != nil {
		return school.School{}, nil
	}

	return school.School{
		ID:     uid,
		Name:   *schoolDto.Name,
		Region: regionNum,
	}, nil
}
