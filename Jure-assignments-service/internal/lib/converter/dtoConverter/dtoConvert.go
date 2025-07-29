package dtoConverter

import (
	"fmt"
	"main/internal/dto/juryAssignmentsDto"
	"main/internal/models/jury_assignments"

	"github.com/google/uuid"
)

func ConvertDTOtoJuryAssignments(dto juryAssignmentsDto.JuryAssignmentsResponseDTO) jury_assignments.JuryAssignments {
	return jury_assignments.JuryAssignments{
		ID:      dto.ID,
		EventID: dto.EventID,
		UserID:  dto.UserID,
	}
}

func FromCreateDTOtoModel(dto juryAssignmentsDto.CreateJuryAssignmentsDTO) (jury_assignments.JuryAssignments, error) {

	juryUuid, err := uuid.Parse(dto.UserID)
	if err != nil {
		return jury_assignments.JuryAssignments{}, fmt.Errorf("%s", err)
	}

	eventUuid, err := uuid.Parse(dto.EventID)
	if err != nil {
		return jury_assignments.JuryAssignments{}, fmt.Errorf("%s", err)
	}

	return jury_assignments.JuryAssignments{
		UserID:  juryUuid,
		EventID: eventUuid,
	}, nil
}

func FromUpdateDTOtoModel(dto juryAssignmentsDto.UpdateJuryAssignmentsDTO, id uuid.UUID) (jury_assignments.JuryAssignments, error) {

	var juryUuid uuid.UUID
	var err error
	if dto.UserID != nil {
		juryUuid, err = uuid.Parse(*dto.UserID)
		if err != nil {
			return jury_assignments.JuryAssignments{}, fmt.Errorf("%s", err)
		}
	}

	var eventUuid uuid.UUID
	if dto.EventID != nil {
		eventUuid, err = uuid.Parse(*dto.EventID)
		if err != nil {
			return jury_assignments.JuryAssignments{}, fmt.Errorf("%s", err)
		}
	}

	return jury_assignments.JuryAssignments{
		ID:      id,
		UserID:  juryUuid,
		EventID: eventUuid,
	}, nil
}

func ConvertManyJuryAssignmentsToDTO(models []jury_assignments.JuryAssignments) []juryAssignmentsDto.JuryAssignmentsResponseDTO {
	result := []juryAssignmentsDto.JuryAssignmentsResponseDTO{}
	for _, model := range models {
		result = append(result, ConvertJuryAssignmentsToDTO(model))
	}
	return result
}

func ConvertJuryAssignmentsToDTO(model jury_assignments.JuryAssignments) juryAssignmentsDto.JuryAssignmentsResponseDTO {
	return juryAssignmentsDto.JuryAssignmentsResponseDTO{
		ID:      model.ID,
		EventID: model.EventID,
		UserID:  model.UserID,
	}
}
