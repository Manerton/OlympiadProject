package request

import "main/internal/dto/event_dto"

type DetailRequest struct {
	Fields *[]string
	Limit  *int
	Offset *int
	event_dto.EventDTO
}
