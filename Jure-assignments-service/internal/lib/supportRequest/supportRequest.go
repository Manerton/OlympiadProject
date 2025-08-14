package supportRequest

import (
	"fmt"
	"main/internal/config"
	"main/internal/lib/response"
	"net/http"

	"github.com/go-chi/render"
	"github.com/google/uuid"
)

const (
	UserService  = "userService"
	EventService = "evetnService"
)

type SupportRequest struct {
	cfg *config.AdditionalAddressesConfig
}

func New(additionalAddress *config.AdditionalAddressesConfig) *SupportRequest {
	return &SupportRequest{
		cfg: additionalAddress,
	}
}

func (r *SupportRequest) PrepareRequest(id uuid.UUID, service string) (bool, error) {
	const op = "lib.supportRequest.PrepareRequest"

	path := ""
	switch service {
	case UserService:
		path = fmt.Sprintf("%s/api/users/%s", r.cfg.UserService, id.String())
	case EventService:
		path = fmt.Sprintf("%s/api/events/%s", r.cfg.EventService, id.String())
	}

	result, err := sendRequest("GET", path)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	if result.StatusCode != http.StatusOK {
		return false, fmt.Errorf("%s: %s", op, result.ErrorCode)
	}

	return true, nil
}

func sendRequest(method string, path string) (response.ApiResponse, error) {
	const op = "lib.supportRequest.sendRequest"
	client := &http.Client{}
	request, err := http.NewRequest(method, path, nil)
	if err != nil {
		return response.ApiResponse{}, fmt.Errorf("%s: %w", op, err)
	}
	nowResponse, err := client.Do(request)
	if err != nil {
		return response.ApiResponse{}, fmt.Errorf("%s: %w", op, err)
	}
	defer nowResponse.Body.Close()

	decodeResponse := response.ApiResponse{}
	err = render.DecodeJSON(nowResponse.Body, &decodeResponse)
	if err != nil {
		return response.ApiResponse{}, fmt.Errorf("%s: %w", op, err)
	}
	return decodeResponse, nil
}
