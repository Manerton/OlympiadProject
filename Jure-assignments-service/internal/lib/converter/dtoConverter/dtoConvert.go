package dtoConverter

import (
	"main/internal/dto/juryAssignmentsDto"
	"main/internal/models/jury_assignments"
)

func ConvertDTOtoJuryAssignments(dto juryAssignmentsDto.JuryAssignmentsDTO) jury_assignments.JuryAssignments {
	return jury_assignments.JuryAssignments{
		ID:      dto.ID,
		EventID: dto.EventID,
		JuryID:  dto.JuryID,
	}
}

func ConvertManyJuryAssignmentsToDTO(models []jury_assignments.JuryAssignments) []juryAssignmentsDto.JuryAssignmentsDTO {
	result := []juryAssignmentsDto.JuryAssignmentsDTO{}
	for _, model := range models {
		result = append(result, ConvertJuryAssignmentsToDTO(model))
	}
	return result
}

func ConvertJuryAssignmentsToDTO(model jury_assignments.JuryAssignments) juryAssignmentsDto.JuryAssignmentsDTO {
	return juryAssignmentsDto.JuryAssignmentsDTO{
		ID:      model.ID,
		EventID: model.EventID,
		JuryID:  model.JuryID,
	}
}
