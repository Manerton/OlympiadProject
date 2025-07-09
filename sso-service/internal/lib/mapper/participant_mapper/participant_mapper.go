package participant_mapper

import (
	register_dto "main/internal/dto/auth/register"
	participant_dto "main/internal/dto/participant"
	"main/internal/models/participant"

	"github.com/google/uuid"
)

func ToDTO(participantModel participant.Participant) participant_dto.ParticipantResponseDTO {
	return participant_dto.ParticipantResponseDTO{
		ID:         participantModel.ID.String(),
		UserId:     participantModel.UserId.String(),
		Disability: participantModel.Disability,
		SchoolId:   participantModel.SchoolId.String(),

		Citizenship: participantModel.Citizenship,
		ClassNumber: participantModel.ClassNumber,
	}
}

func FromRegisterToModel(registerDTO *register_dto.RegisterParticipantRequestDTO, userId uuid.UUID) participant.Participant {
	schoolId, err := uuid.Parse(registerDTO.SchoolId)
	if err != nil {

	}

	return participant.Participant{
		UserId:      userId,
		Disability:  registerDTO.Disability,
		SchoolId:    schoolId,
		Citizenship: registerDTO.Citizenship,
		ClassNumber: registerDTO.ClassNumber,
	}
}

func FromUpdateToModel(updateDTO participant_dto.UpdateParticipantRequestDTO, uid uuid.UUID) participant.Participant {
	return participant.Participant{
		ID:          uid,
		Disability:  *updateDTO.Disability,
		Citizenship: *updateDTO.Citizenship,
		ClassNumber: *updateDTO.ClassNumber,
	}
}
