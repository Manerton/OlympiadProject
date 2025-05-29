package user

import user_dto "main/internal/dto/user"

func ToDTO(user User) user_dto.UserResponseDTO {
	return user_dto.UserResponseDTO{
		ID:        user.ID.String(),
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Surname:   user.Surname,
		Sex:       user.Sex,
		// Role string or int?
		Role:        string(user.Role),
		PhoneNumber: user.PhoneNumber,
		BirthDate:   user.BirthDate.String(),
	}
}
