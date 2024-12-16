package dtoConverter

import (
	"main/internal/dto/juryAssignmentsDto"
	"main/internal/models/juryassignments"
)

func ConvertDTOtoJuryAssignments(dto juryAssignmentsDto.JuryAssignmentsDTO) juryassignments.JuryAssignments {
	return juryassignments.JuryAssignments{
		ID:      dto.ID,
		EventID: dto.EventID,
		JuryID:  dto.JuryID,
	}
}

func ConvertManyJuryAssignmentsToDTO(models []juryassignments.JuryAssignments) []juryAssignmentsDto.JuryAssignmentsDTO {
	result := []juryAssignmentsDto.JuryAssignmentsDTO{}
	for _, model := range models {
		result = append(result, ConvertJuryAssignmentsToDTO(model))
	}
	return result
}

func ConvertJuryAssignmentsToDTO(model juryassignments.JuryAssignments) juryAssignmentsDto.JuryAssignmentsDTO {
	return juryAssignmentsDto.JuryAssignmentsDTO{
		ID:      model.ID,
		EventID: model.EventID,
		JuryID:  model.JuryID,
	}
}
