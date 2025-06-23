package school_handler

import (
	"context"
	school_dto "main/internal/dto/school"
	"main/internal/lib/parser"
	"main/internal/lib/response"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type SchoolService interface {
	GetAll(ctx context.Context, page, limit *int) ([]school_dto.SchoolResponeDTO, error)
	GetById(ctx context.Context, id string) (school_dto.SchoolResponeDTO, error)

	Update(ctx context.Context, id string, schoolDTO school_dto.UpdateSchoolRequestDTO) error
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

	render.JSON(w, r, response.ApiResponse{
		Status:     response.SUCCESS,
		StatusCode: http.StatusOK,
		Data:       schoolResponse,
	})
}

func (h *SchoolHandler) GetById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	schoolResponse, err := h.schoolService.GetById(ctx, id)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.ErrorResponse(err.Error()))
		return
	}

	render.JSON(w, r, response.ApiResponse{
		Status:     response.SUCCESS,
		StatusCode: http.StatusOK,
		Data:       schoolResponse,
	})
}

func (h *SchoolHandler) Create(w http.ResponseWriter, r *http.Request) {
	// TODO!! Create?
}

func (h *SchoolHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	schoolDTO := school_dto.UpdateSchoolRequestDTO{}
	err := render.DecodeJSON(r.Body, &schoolDTO)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.ErrorResponse("failed decode json"))
		return
	}

	err = h.schoolService.Update(ctx, id, schoolDTO)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.ErrorResponse(err.Error()))
		return
	}

	render.JSON(w, r, response.SuccessResponse("success update"))
}
