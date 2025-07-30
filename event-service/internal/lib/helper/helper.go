package helper

import (
	"encoding/json"
	"fmt"
	"main/internal/dto/rabbit_dto"
)

func PayloadDeleteConstructor(table string, id string) ([]byte, error) {
	const op = "helper.PayloadDeleteConstructor"

	rabbitDTO := rabbit_dto.RabbitDTO{
		AppName: "EventService",
		Method:  "delete",
		Data: rabbit_dto.RabbitData{
			Table:            table,
			SearchAttributes: map[string]any{"event_id": id},
		},
	}

	result, err := json.Marshal(rabbitDTO)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return result, nil
}
