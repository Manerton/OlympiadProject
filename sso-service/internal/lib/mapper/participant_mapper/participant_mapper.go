package paricipant_mapper

import (
	register_dto "main/internal/dto/auth/register"
	"main/internal/models/participant"

	"github.com/google/uuid"
)

func FromRegisterToModel(registerDTO *register_dto.RegisterParticipantRequestDTO, userId uuid.UUID) participant.Participant {
	return participant.Participant{
		UserId:      userId,
		OVZ:         registerDTO.OVZ,
		SchoolName:  registerDTO.SchoolName,
		City:        registerDTO.City,
		Citizenship: registerDTO.Citizenship,
		Reason:      registerDTO.Reason,
		ClassNumber: registerDTO.ClassNumber,
	}
}
