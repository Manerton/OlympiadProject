package dtoConverter

import (
	"main/internal/dto/juryAssignmentsDto"
	"main/internal/models/juryAssignments"
)

func ConvertDTOtoJuryAssignments(dto juryAssignmentsDto.JuryAssignmentsDTO) juryAssignments.JuryAssignments {
	return juryAssignments.JuryAssignments{
		ID:      dto.ID,
		EventID: dto.EventID,
		JuryID:  dto.JuryID,
	}
}

func ConvertManyJuryAssignmentsToDTO(models []juryAssignments.JuryAssignments) []juryAssignmentsDto.JuryAssignmentsDTO {
	result := []juryAssignmentsDto.JuryAssignmentsDTO{}
	for _, model := range models {
		result = append(result, ConvertJuryAssignmentsToDTO(model))
	}
	return result
}

func ConvertJuryAssignmentsToDTO(model juryAssignments.JuryAssignments) juryAssignmentsDto.JuryAssignmentsDTO {
	return juryAssignmentsDto.JuryAssignmentsDTO{
		ID:      model.ID,
		EventID: model.EventID,
		JuryID:  model.JuryID,
	}
}
