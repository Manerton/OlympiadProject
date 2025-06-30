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

func (r *SupportRequest) PrepareRequest(id uuid.UUID, service string) (bool, error) {
	const op = "lib.supportRequest.PrepareRequest"

	path := ""
	switch service {
	case UserService:
		path = fmt.Sprintf("%s/user/%d", r.cfg.JuryService, id)
	case EventService:
		path = fmt.Sprintf("%s/user/%d", r.cfg.EventService, id)
	}

	result, err := SendRequest("GET", path)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	if result.StatusCode != http.StatusOK {
		return false, fmt.Errorf("%s: %s", op, result.Error)
	}

	return true, nil
}

func SendRequest(method string, path string) (response.ApiResponse, error) {
	const op = "lib.supportRequest.SendRequest"
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
