package rabbit_dto

type RabbitDTO struct {
	AppName string     `json:"appName"`
	Method  string     `json:"method"`
	Data    RabbitData `json:"data"`
}

type RabbitData struct {
	Table            string            `json:"table"`
	Attributes       map[string]string `json:"attributes"`
	SearchAttributes map[string]string `json:"searchAttributes"`
}
