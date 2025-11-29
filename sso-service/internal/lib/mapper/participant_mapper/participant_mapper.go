package participant_mapper

import (
	"fmt"
	register_dto "main/internal/dto/auth/register"
	participant_dto "main/internal/dto/participant"
	"main/internal/models/participant"
	"strconv"

	"github.com/google/uuid"
)

func ToDTO(participantModel participant.Participant) participant_dto.ParticipantResponseDTO {

	disability := strconv.Itoa(participantModel.Disability)
	citizenship := strconv.Itoa(participantModel.Citizenship)
	classNumber := strconv.Itoa(participantModel.ClassNumber)

	return participant_dto.ParticipantResponseDTO{
		ParticipantId: participantModel.ID.String(),
		UserId:        participantModel.UserId.String(),
		Disability:    disability,
		SchoolId:      participantModel.SchoolId.String(),

		Citizenship: citizenship,
		ClassNumber: classNumber,
	}
}

func FromRegisterToModel(registerDTO *register_dto.RegisterParticipantRequestDTO, userId uuid.UUID) participant.Participant {
	schoolId, err := uuid.Parse(registerDTO.SchoolId)
	if err != nil {

	}

	disability, err := strconv.Atoi(registerDTO.Disability)
	citizenship, err := strconv.Atoi(registerDTO.Citizenship)
	classNumber, err := strconv.Atoi(registerDTO.ClassNumber)

	return participant.Participant{
		UserId:      userId,
		Disability:  disability,
		SchoolId:    schoolId,
		Citizenship: citizenship,
		ClassNumber: classNumber,
	}
}

func FromUpdateToModel(updateDTO participant_dto.UpdateParticipantRequestDTO, uid uuid.UUID) (participant.Participant, error) {
	var model participant.Participant
	model.ID = uid

	if updateDTO.Disability != nil {
		val, err := strconv.Atoi(*updateDTO.Disability)
		if err != nil {
			return model, fmt.Errorf("invalid disability value: %w", err)
		}
		model.Disability = val
	}

	if updateDTO.Citizenship != nil {
		val, err := strconv.Atoi(*updateDTO.Citizenship)
		if err != nil {
			return model, fmt.Errorf("invalid citizenship value: %w", err)
		}
		model.Citizenship = val
	}

	if updateDTO.ClassNumber != nil {
		model.ClassNumber = *updateDTO.ClassNumber
	}

	if updateDTO.SchoolId != nil {
		sid, err := uuid.Parse(*updateDTO.SchoolId)
		if err != nil {
			return model, fmt.Errorf("invalid school_id value: %w", err)
		}
		model.SchoolId = sid
	}

	return model, nil
}
