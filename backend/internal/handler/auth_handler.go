package handler

import (
	"net/http"

	"servixa/internal/model"
	"servixa/internal/response"
	"servixa/internal/service"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(
	authService *service.AuthService,
) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserResponse struct {
	ID       uint           `json:"id"`
	Name     string         `json:"name"`
	Email    string         `json:"email"`
	Role     model.UserRole `json:"role"`
	IsActive bool           `json:"is_active"`
}
type RegisterResponse struct {
	ID       uint           `json:"id"`
	Name     string         `json:"name"`
	Email    string         `json:"email"`
	Role     model.UserRole `json:"role"`
	IsActive bool           `json:"is_active"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	User         UserResponse `json:"user"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type RefreshTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type LogoutRequest struct {
	RefreshToken string `json:"refresh_token"`
}

func (h *AuthHandler) Register(c *gin.Context) {
	var request RegisterRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		response.Error(
			c,
			http.StatusBadRequest,
			"invalid request body",
			"INVALID_REQUEST",
		)
		return
	}

	user, err := h.authService.Register(
		request.Name,
		request.Email,
		request.Password,
	)

	if err != nil {
		switch err.Error() {
		case "name is required":
			response.Error(
				c,
				http.StatusBadRequest,
				err.Error(),
				"VALIDATION_ERROR",
			)
		case "email is required":
			response.Error(
				c,
				http.StatusBadRequest,
				err.Error(),
				"VALIDATION_ERROR",
			)
		case "password is required":
			response.Error(
				c,
				http.StatusBadRequest,
				err.Error(),
				"VALIDATION_ERROR",
			)
		case "email already registered":
			response.Error(
				c,
				http.StatusConflict,
				err.Error(),
				"EMAIL_ALREADY_REGISTERED",
			)
		default:
			response.Error(
				c,
				http.StatusInternalServerError,
				"failed to register user",
				"INTERNAL_ERROR",
			)
		}

		return
	}

	response.Success(
		c,
		http.StatusCreated,
		"user registered successfully",
		RegisterResponse{
			ID:       user.ID,
			Name:     user.Name,
			Email:    user.Email,
			Role:     user.Role,
			IsActive: user.IsActive,
		},
	)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var request LoginRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		response.Error(
			c,
			http.StatusBadRequest,
			"invalid request body",
			"INVALID_REQUEST",
		)
		return
	}

	result, err := h.authService.Login(
		request.Email,
		request.Password,
	)

	if err != nil {
		switch err.Error() {
		case "email is required":
			response.Error(
				c,
				http.StatusBadRequest,
				err.Error(),
				"VALIDATION_ERROR",
			)

		case "password is required":
			response.Error(
				c,
				http.StatusBadRequest,
				err.Error(),
				"VALIDATION_ERROR",
			)

		case "invalid email or password":
			response.Error(
				c,
				http.StatusUnauthorized,
				err.Error(),
				"INVALID_CREDENTIALS",
			)

		case "user is inactive":
			response.Error(
				c,
				http.StatusForbidden,
				err.Error(),
				"USER_INACTIVE",
			)

		default:
			response.Error(
				c,
				http.StatusInternalServerError,
				"failed to login",
				"INTERNAL_ERROR",
			)
		}

		return
	}

	response.Success(
		c,
		http.StatusOK,
		"login successful",
		LoginResponse{
			AccessToken:  result.AccessToken,
			RefreshToken: result.RefreshToken,
			User: UserResponse{
				ID:       result.User.ID,
				Name:     result.User.Name,
				Email:    result.User.Email,
				Role:     result.User.Role,
				IsActive: result.User.IsActive,
			},
		},
	)
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	var request RefreshTokenRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		response.Error(
			c,
			http.StatusBadRequest,
			"invalid request body",
			"INVALID_REQUEST",
		)
		return
	}

	accessToken, refreshToken, err := h.authService.RefreshAccessToken(
		request.RefreshToken,
	)

	if err != nil {
		switch err.Error() {
		case "refresh token is required",
			"invalid refresh token",
			"refresh token has been revoked",
			"refresh token has expired":
			response.Error(
				c,
				http.StatusUnauthorized,
				err.Error(),
				"INVALID_REFRESH_TOKEN",
			)

		case "user is inactive":
			response.Error(
				c,
				http.StatusForbidden,
				err.Error(),
				"USER_INACTIVE",
			)

		default:
			response.Error(
				c,
				http.StatusInternalServerError,
				"failed to refresh access token",
				"INTERNAL_ERROR",
			)
		}

		return
	}

	response.Success(
		c,
		http.StatusOK,
		"access token refreshed successfully",
		RefreshTokenResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	var request LogoutRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		response.Error(
			c,
			http.StatusBadRequest,
			"invalid request body",
			"INVALID_REQUEST",
		)
		return
	}

	if err := h.authService.Logout(request.RefreshToken); err != nil {
		switch err.Error() {
		case "refresh token is required",
			"invalid refresh token",
			"refresh token has been revoked",
			"refresh token has expired":
			response.Error(
				c,
				http.StatusUnauthorized,
				err.Error(),
				"INVALID_REFRESH_TOKEN",
			)

		default:
			response.Error(
				c,
				http.StatusInternalServerError,
				"failed to logout",
				"INTERNAL_ERROR",
			)
		}

		return
	}

	response.Success(
		c,
		http.StatusOK,
		"logout successful",
		nil,
	)
}
