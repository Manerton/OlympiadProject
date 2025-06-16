package school_mapper

import (
	school_dto "main/internal/dto/school"
	"main/internal/models/school"
)

func FromModelToDTO(school school.School) school_dto.SchoolResponeDTO {
	return school_dto.SchoolResponeDTO{
		ID:     school.ID.String(),
		Name:   school.Name,
		Region: school.Region,
	}
}

func FromCreateDTOToModel(schoolDto school_dto.CreateSchoolRequestDTO) school.School {
	return school.School{
		Name:   schoolDto.Name,
		Region: schoolDto.Region,
	}
}

func FromUpdateDTOToModel(schoolDto school_dto.UpdateSchoolRequestDTO) school.School {
	return school.School{
		Name:   schoolDto.Name,
		Region: schoolDto.Region,
	}
}
