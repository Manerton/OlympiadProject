package subject_handler

import (
	"log/slog"
	"main/internal/lib/response"
	"main/internal/models/subject"
	"net/http"

	"github.com/go-chi/render"
)

type SubjectHandler struct {
	// service        *subject_service.SubjectService
	subjectStorage *subject.SubjectStorage
	log            *slog.Logger
}

func NewSubjectHandler(storage *subject.SubjectStorage, log *slog.Logger) *SubjectHandler {
	return &SubjectHandler{subjectStorage: storage, log: log}
}

func (h *SubjectHandler) GetAllSubjects(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.subject_handler.GetAllSubjects"
	log := h.log.With(
		slog.String("op", op),
	)

	subjects := h.subjectStorage.GetAllSubjects()
	// if err != nil {
	// 	log.Error("failed to get subject", liblogger.Err(err))
	// 	render.JSON(w, r, response.Error("failed to get subject"))
	// 	return
	// }
	log.Info("subject getted")

	render.JSON(w, r, response.ApiResponse{
		Status: response.StatusOK,
		Data:   subjects,
	})
}

// func (h *SubjectHandler) GetSubjectByID(w http.ResponseWriter, r *http.Request) {
// 	const op = "handlers.subject_handler.GetSubjectByID"
// 	log := h.log.With(
// 		slog.String("op", op),
// 	)

// 	received_id := chi.URLParam(r, "id")
// 	searched_id, err := strconv.ParseUint(received_id, 10, 32)
// 	if err != nil {
// 		log.Error("failed to parse id to uint", slog.String("received id", received_id), liblogger.Err(err))
// 	}
// 	log.Info("subject_id on request body decoded", slog.Any("subject_id", searched_id))

// 	subject_dto, err := h.service.GetSubjectByID(uint(searched_id))
// 	if err != nil {
// 		log.Error("failed to get subject", liblogger.Err(err))
// 		render.JSON(w, r, response.Error("failed to get subject"))
// 		return
// 	}
// 	log.Info("subject getted", slog.Any("subject", subject_dto))

// 	render.JSON(w, r, response.ApiResponse{
// 		Status: response.StatusOK,
// 		Data:   subject_dto,
// 	})
// }

// func (h *SubjectHandler) CreateSubject(w http.ResponseWriter, r *http.Request) {
// 	const op = "handlers.subject_handler.CreateSubject"
// 	log := h.log.With(
// 		slog.String("op", op),
// 	)

// 	var subject_dto subject_dto.SubjectDTO
// 	err := render.DecodeJSON(r.Body, &subject_dto)
// 	if errors.Is(err, io.EOF) {
// 		log.Error("render body is empty")
// 		render.JSON(w, r, response.Error("empty request"))
// 		return
// 	}
// 	if err != nil {
// 		log.Error("failed to decode request body", liblogger.Err(err))
// 		render.JSON(w, r, response.Error("failde to decode request"))
// 		return
// 	}
// 	log.Info("subject on request body decoded", slog.Any("subject", subject_dto))
// 	id, err := h.service.CreateSubject(subject_dto)
// 	if err != nil {
// 		log.Error("failed to add subject", liblogger.Err(err))
// 		render.JSON(w, r, response.Error("failed to add subject"))
// 		return
// 	}
// 	log.Info("subject added", slog.Any("id", id))

// 	render.JSON(w, r, response.Success(fmt.Sprintf("id = %d", id)))
// }

// func (h *SubjectHandler) UpdateSubject(w http.ResponseWriter, r *http.Request) {
// 	const op = "handlers.subject_handler.UpdateSubject"
// 	log := h.log.With(
// 		slog.String("op", op),
// 	)

// 	var subject_dto subject_dto.SubjectDTO
// 	err := render.DecodeJSON(r.Body, &subject_dto)
// 	if errors.Is(err, io.EOF) {
// 		log.Error("render body is empty")
// 		render.JSON(w, r, response.Error("empty request"))
// 		return
// 	}
// 	if err != nil {
// 		log.Error("failed to decode request body", liblogger.Err(err))
// 		render.JSON(w, r, response.Error("failde to decode request"))
// 		return
// 	}
// 	log.Info("subject on request body decoded", slog.Any("subject", subject_dto))
// 	id, err := h.service.UpdateSubject(subject_dto)
// 	if err != nil {
// 		log.Error("failed to update subject", liblogger.Err(err))
// 		render.JSON(w, r, response.Error("failed to update subject"))
// 		return
// 	}
// 	log.Info("subject updated", slog.Any("id", id))

// 	render.JSON(w, r, response.Success(fmt.Sprintf("id = %d", id)))
// }

// func (h *SubjectHandler) DeleteSubject(w http.ResponseWriter, r *http.Request) {
// 	const op = "handlers.subject_handler.DeleteSubject"
// 	log := h.log.With(
// 		slog.String("op", op),
// 	)

// 	received_id := chi.URLParam(r, "id")
// 	deleted_id, err := strconv.ParseUint(received_id, 10, 32)
// 	if err != nil {
// 		log.Error("failed to parse id to uint", slog.String("received id", received_id), liblogger.Err(err))
// 		render.JSON(w, r, response.Error("failed to parse id"))
// 	}
// 	log.Info("subject_id on request body decoded", slog.Any("subject_id", deleted_id))

// 	err = h.service.DeleteSubject(uint(deleted_id))
// 	if err != nil {
// 		log.Error("failed to delete subject", liblogger.Err(err))
// 		render.JSON(w, r, response.Error("failed to delete subject"))
// 		return
// 	}
// 	log.Info("subject deleted", slog.Any("id", deleted_id))

// 	render.JSON(w, r, response.Success("Object deleted"))
// }
