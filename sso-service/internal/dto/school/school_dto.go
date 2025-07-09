package school_dto

type CreateSchoolRequestDTO struct {
	Name   string `json:"name"`
	Region int    `json:"region"`
}

type UpdateSchoolRequestDTO struct {
	Name   *string `json:"name"`
	Region *int    `json:"region"`
}

type SchoolResponeDTO struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Region int    `json:"region"`
}
