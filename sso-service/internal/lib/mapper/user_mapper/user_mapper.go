package user_mapper

import (
	register_dto "main/internal/dto/auth/register"
	user_dto "main/internal/dto/user"
	"main/internal/models/user"
)

func ToDTO(user user.User) user_dto.UserResponseDTO {
	return user_dto.UserResponseDTO{
		ID:          user.ID.String(),
		Email:       user.Email,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Surname:     user.Surname,
		Gender:      user.Gender,
		Role:        string(user.Role),
		PhoneNumber: user.PhoneNumber,
		BirthDate:   user.BirthDate.String(),
		Activated:   user.Activated,
	}
}

func FromRegisterToModel(regiterDTO *register_dto.RegisterParticipantRequestDTO) user.User {
	return user.User{
		Email:       regiterDTO.Email,
		FirstName:   regiterDTO.FirstName,
		LastName:    regiterDTO.LastName,
		Surname:     regiterDTO.Surname,
		PhoneNumber: regiterDTO.PhoneNumber,
		BirthDate:   regiterDTO.BirthDate,
		Gender:      regiterDTO.Gender,
		Role:        user.PARTICIPANT,
		Activated:   false,
	}
}

func FromRegisterUserToModel(regiterDTO *register_dto.RegusterUserRequestDTO) user.User {
	return user.User{
		Email:       regiterDTO.Email,
		FirstName:   regiterDTO.FirstName,
		LastName:    regiterDTO.LastName,
		Surname:     regiterDTO.Surname,
		PhoneNumber: regiterDTO.PhoneNumber,
		Gender:      regiterDTO.Gender,
		Role:        user.RoleType(regiterDTO.Role),
		BirthDate:   regiterDTO.BirthDate,
		Activated:   true,
	}
}

func FromUpdateToModel(updateDTO *user_dto.UpdateUserRequestDTO) user.User {
	return user.User{
		Email:       *updateDTO.Email,
		FirstName:   *updateDTO.FirstName,
		LastName:    *updateDTO.LastName,
		Surname:     *updateDTO.Surname,
		PhoneNumber: *updateDTO.PhoneNumber,
	}
}
