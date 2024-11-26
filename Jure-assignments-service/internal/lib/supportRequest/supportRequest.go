package supportRequest

import (
	"fmt"
	"net/http"

	"github.com/go-chi/render"
)

type Response struct {
	Status string
	Data   interface{}
}

func SupportRequest(method string, path string) (Response, error) {
	const op = "lib.supportRequest.SupportRequest"
	client := &http.Client{}
	request, err := http.NewRequest(method, path, nil)
	if err != nil {
		return Response{}, fmt.Errorf("%s: %w", op, err)
	}
	response, err := client.Do(request)
	if err != nil {
		return Response{}, fmt.Errorf("%s: %w", op, err)
	}
	defer response.Body.Close()

	decodeResponse := Response{}
	err = render.DecodeJSON(response.Body, &decodeResponse)
	if err != nil {
		return Response{}, fmt.Errorf("%s: %w", op, err)
	}
	return decodeResponse, nil
}
