package school_handler

import (
	"context"
	school_dto "main/internal/dto/school"
	"main/internal/lib/parser"
	"main/internal/lib/response"
	"net/http"

	"github.com/go-chi/render"
)

type SchoolService interface {
	GetAll(ctx context.Context, page, limit *int) ([]school_dto.SchoolResponeDTO, error)
}

type SchoolHandler struct {
	schoolService SchoolService
}

func New(schoolService SchoolService) *SchoolHandler {
	return &SchoolHandler{
		schoolService: schoolService,
	}
}

func (h *SchoolHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page, limit, err := parser.ParsePageLimit(pageStr, limitStr)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.ErrorResponse("failed parse page/limit"))
		return
	}

	schoolResponse, err := h.schoolService.GetAll(ctx, page, limit)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.ErrorResponse(err.Error()))
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, response.ApiResponse{
		Status:     response.SUCCESS,
		StatusCode: http.StatusOK,
		Data:       schoolResponse,
	})
}
