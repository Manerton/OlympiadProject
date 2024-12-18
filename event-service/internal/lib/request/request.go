package request

import "main/internal/dto/event_dto"

type DetailRequest struct {
	Fields *[]string `json:"fields"`
	Limit  *int      `json:"limit"`
	Offset *int      `json:"offset"`
	event_dto.EventDTO
}
