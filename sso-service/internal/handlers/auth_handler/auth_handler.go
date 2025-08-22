package auth_handler

import (
	"context"
	login_dto "main/internal/dto/auth/login"
	register_dto "main/internal/dto/auth/register"
	"main/internal/lib/errs"
	"main/internal/lib/response"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type AuthService interface {
	Login(ctx context.Context, loginRequest *login_dto.LoginRequestDTO) (*login_dto.AuthResultDTO, error)
	Logout(ctx context.Context, tokeId string) error

	RegisterParticipant(ctx context.Context, registerRequest *register_dto.RegisterParticipantRequestDTO) error
	RegisterUser(ctx context.Context, userRequest *register_dto.RegisterUserRequestDTO) error

	Refresh(ctx context.Context, refershToken string) (*login_dto.AuthResultDTO, error)
	RevokeAllUserTokens(ctx context.Context, userId string) error
	RevokeToken(ctx context.Context, id string) error
}

type AuthHandler struct {
	authService AuthService
}

func New(authService AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// @Summary Login
// @Description Авторизация пользователя
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body login_dto.LoginRequestDTO true "Данные для входа"
// @Success 200 {object} response.ApiResponse
// @Failure 401 {object} response.ApiResponse
// @Router /api/users/login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var loginRequest login_dto.LoginRequestDTO

	err := render.DecodeJSON(r.Body, &loginRequest)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.ErrorApiResponse(errs.ErrBadRequest.Wrap("failed decode body")))
		return
	}

	authResult, err := h.authService.Login(ctx, &loginRequest)
	if err != nil {
		if apiErr, ok := errs.IsApiError(err); ok {
			render.Status(r, apiErr.HttpCode)
			render.JSON(w, r, response.ErrorApiResponse(apiErr))
			return
		}

		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, response.ErrorApiResponse(errs.ErrInternalError))
		return
	}

	// Set the JWT token in a cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    authResult.RefreshToken,
		HttpOnly: true,
		Secure:   false,
		// Domain:   "172.16.1.39",
		SameSite: http.SameSiteLaxMode,
		// Path:     "api/",
		Expires: time.Now().Add(time.Duration(authResult.ExpiresInRefresh) * time.Second), // Match the token expiration
	})

	render.Status(r, http.StatusOK)
	render.JSON(w, r, response.ApiResponse{
		Status:     "SUCCESS",
		StatusCode: http.StatusOK,
		Data: login_dto.LoginResponseDTO{
			AccessToken: authResult.AccessToken,
			ExpiresIn:   authResult.ExpiresInAccess,
		},
	})
}

// @Summery Logout
// @Description Выход из учётной записи
// @Tags auth
// @Produce json
// @Success 200 {object} response.ApiResponse
// @Failure 400 {object} response.ApiResponse
// @Router /api/users/logout [post]
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.ErrorApiResponse(errs.ErrTokenNotFound.Wrap("refresh token not found")))
		return
	}

	cookieStr := cookie.Value
	err = h.authService.Logout(ctx, cookieStr)
	if err != nil {
		if apiErr, ok := errs.IsApiError(err); ok {
			render.Status(r, apiErr.HttpCode)
			render.JSON(w, r, response.ErrorApiResponse(apiErr))
			return
		}

		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, response.ErrorApiResponse(errs.ErrInternalError))
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		HttpOnly: true,
		Secure:   false,
		// Domain:   "172.16.1.39",
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Unix(0, 0), // дата в прошлом
		MaxAge:   -1,              // сразу удалить
	})

	render.JSON(w, r, response.SuccessResponse("success logout"))
}

// @Summery Register
// @Description Регистрация пользователя как ученика
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body register_dto.RegisterParticipantRequestDTO true "Данные для регистрации"
// @Success 200 {object} response.ApiResponse
// @Failure 400 {object} response.ApiResponse
// @Router /api/users/register [post]
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var registerRequest register_dto.RegisterParticipantRequestDTO

	err := render.DecodeJSON(r.Body, &registerRequest)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.ErrorApiResponse(errs.ErrBadRequest.Wrap("failed decode body")))
		return
	}

	err = h.authService.RegisterParticipant(ctx, &registerRequest)
	if err != nil {
		if apiErr, ok := errs.IsApiError(err); ok {
			render.Status(r, apiErr.HttpCode)
			render.JSON(w, r, response.ErrorApiResponse(apiErr))
			return
		}

		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, response.ErrorApiResponse(errs.ErrInternalError))
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, response.SuccessResponse("Register success"))
}

// @Summery AdminRegister
// @Security BearerAuth
// @Description Регистрация пользователя для админа
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body register_dto.RegisterUserRequestDTO true "Данные для регистрации от панели админа"
// @Success 200 {object} response.ApiResponse
// @Failure 400 {object} response.ApiResponse
// @Router /api/byadmin/register [post]
func (h *AuthHandler) AdminRegister(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	registerRequest := &register_dto.RegisterUserRequestDTO{}

	err := render.DecodeJSON(r.Body, registerRequest)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.ErrorApiResponse(errs.ErrBadRequest.Wrap("failed decode body")))
		return
	}

	err = h.authService.RegisterUser(ctx, registerRequest)
	if err != nil {
		if apiErr, ok := errs.IsApiError(err); ok {
			render.Status(r, apiErr.HttpCode)
			render.JSON(w, r, response.ErrorApiResponse(apiErr))
			return
		}

		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, response.ErrorApiResponse(errs.ErrInternalError))
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, response.SuccessResponse("register success"))
}

// @Summary Refresh
// @Description Обновление access-токена по refresh-токену из cookie
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponse{data=login_dto.LoginResponseDTO}
// @Failure 400 {object} response.ApiResponse
// @Router /api/users/refresh [post]
func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	refreshToken, err := r.Cookie("refresh_token")
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.ErrorApiResponse(errs.ErrTokenNotFound.Wrap("refresh token not found")))
		return
	}

	loginDTO, err := h.authService.Refresh(ctx, refreshToken.Value)
	if err != nil {
		if apiErr, ok := errs.IsApiError(err); ok {
			render.Status(r, apiErr.HttpCode)
			render.JSON(w, r, response.ErrorApiResponse(apiErr))
			return
		}

		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, response.ErrorApiResponse(errs.ErrInternalError))
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    loginDTO.RefreshToken,
		HttpOnly: true,
		Secure:   false,
		// Domain:   "172.16.1.39",
		SameSite: http.SameSiteLaxMode,
		// Path:     "api/users/refresh",
		Expires: time.Now().Add(time.Duration(loginDTO.ExpiresInRefresh) * time.Second), // Match the token expiration
	})

	render.Status(r, http.StatusOK)
	render.JSON(w, r, response.ApiResponse{
		Status:     "SUCCESS",
		StatusCode: http.StatusOK,
		Data: login_dto.LoginResponseDTO{
			AccessToken: loginDTO.AccessToken,
			ExpiresIn:   loginDTO.ExpiresInAccess,
		},
	})
}

// @Summery Revoke token
// @Security BearerAuth
// @Description Блокировка refresh токена
// @Tags auth
// @Produce json
// @Param id path string true "id refresh токена"
// @Success 200 {object} response.ApiResponse
// @Failure 400 {object} response.ApiResponse
// @Router /api/users/revoke/{id} [post]
func (h *AuthHandler) RevokeToken(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	err := h.authService.RevokeToken(ctx, id)
	if err != nil {
		if apiErr, ok := errs.IsApiError(err); ok {
			render.Status(r, apiErr.HttpCode)
			render.JSON(w, r, response.ErrorApiResponse(apiErr))
			return
		}

		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, response.ErrorApiResponse(errs.ErrInternalError))
		return
	}

	render.JSON(w, r, response.SuccessResponse("success token revoked"))
}

// @Summery Revoke all token by user
// @Security BearerAuth
// @Description Блокировка всех refresh токенов пользователя
// @Tags auth
// @Produce json
// @Param id path string true "id пользователя"
// @Success 200 {object} response.ApiResponse
// @Failure 400 {object} response.ApiResponse
// @Router /api/users/revoke-all/{id} [post]
func (h *AuthHandler) RevokeAllUserTokens(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID := chi.URLParam(r, "id")

	err := h.authService.RevokeAllUserTokens(ctx, userID)
	if err != nil {
		if apiErr, ok := errs.IsApiError(err); ok {
			render.Status(r, apiErr.HttpCode)
			render.JSON(w, r, response.ErrorApiResponse(apiErr))
			return
		}

		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, response.ErrorApiResponse(errs.ErrInternalError))
		return
	}

	render.JSON(w, r, response.SuccessResponse("success tokens revoked"))
}
