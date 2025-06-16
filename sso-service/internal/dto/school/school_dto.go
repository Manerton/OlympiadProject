package school_dto

type CreateSchoolRequestDTO struct {
	Name   string `json:"name"`
	Region string `json:"region"`
}

type UpdateSchoolRequestDTO struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Region string `json:"region"`
}

type SchoolResponeDTO struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Region string `json:"region"`
}
