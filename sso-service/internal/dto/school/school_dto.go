package school_dto

type CreateSchoolRequestDTO struct {
	Name   string `json:"name"`
	Region string `json:"region"`
}

type UpdateSchoolRequestDTO struct {
	Name   *string `json:"name"`
	Region *string `json:"region"`
}

type SchoolResponseDTO struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Region int    `json:"region"`
}
