package participant_mapper

import (
	register_dto "main/internal/dto/auth/register"
	participant_dto "main/internal/dto/participant"
	"main/internal/models/participant"

	"github.com/google/uuid"
)

func ToDTO(participantModel participant.Participant) participant_dto.ParticipantResponseDTO {
	return participant_dto.ParticipantResponseDTO{
		UserId:      participantModel.UserId.String(),
		Disability:  participantModel.Disability,
		SchoolId:    participantModel.SchoolId.String(),
		City:        participantModel.City,
		Reason:      participantModel.Reason,
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
		City:        registerDTO.City,
		Citizenship: registerDTO.Citizenship,
		Reason:      registerDTO.Reason,
		ClassNumber: registerDTO.ClassNumber,
	}
}

func FromUpdateToModel(updateDTO participant_dto.UpdateParticipantRequestDTO, uid uuid.UUID) participant.Participant {
	return participant.Participant{
		ID:          uid,
		Disability:  *updateDTO.Disability,
		City:        *updateDTO.City,
		Reason:      *updateDTO.Reason,
		Citizenship: *updateDTO.Citizenship,
		ClassNumber: *updateDTO.ClassNumber,
	}
}
