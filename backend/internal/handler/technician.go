package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"servixa/internal/middleware"
	"servixa/internal/response"
	"servixa/internal/service"
)

type TechnicianHandler struct {
	technicianService *service.TechnicianService
}

func NewTechnicianHandler(
	technicianService *service.TechnicianService,
) *TechnicianHandler {
	return &TechnicianHandler{
		technicianService: technicianService,
	}
}

type CreateTechnicianRequest struct {
	Name           string `json:"name"`
	Email          string `json:"email"`
	Phone          string `json:"phone"`
	Specialization string `json:"specialization"`
}

type TechnicianResponse struct {
	ID             uint   `json:"id"`
	UserID         uint   `json:"user_id"`
	Name           string `json:"name"`
	Email          string `json:"email"`
	Phone          string `json:"phone"`
	Specialization string `json:"specialization"`
	IsAvailable    bool   `json:"is_available"`
	IsActive       bool   `json:"is_active"`
}

func (h *TechnicianHandler) Create(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		response.Error(
			c,
			http.StatusUnauthorized,
			"unauthorized",
			"UNAUTHORIZED",
		)
		return
	}

	var request CreateTechnicianRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		response.Error(
			c,
			http.StatusBadRequest,
			"invalid request body",
			"VALIDATION_ERROR",
		)
		return
	}

	technician, err := h.technicianService.Create(
		userID,
		request.Name,
		request.Email,
		request.Phone,
		request.Specialization,
	)

	if err != nil {
		switch err.Error() {
		case "name is required",
			"email is required",
			"phone is required",
			"specialization is required":
			response.Error(
				c,
				http.StatusBadRequest,
				err.Error(),
				"VALIDATION_ERROR",
			)

		case "user not found":
			response.Error(
				c,
				http.StatusNotFound,
				err.Error(),
				"USER_NOT_FOUND",
			)

		case "only technician users can create technician profile":
			response.Error(
				c,
				http.StatusForbidden,
				err.Error(),
				"FORBIDDEN",
			)

		case "user is inactive":
			response.Error(
				c,
				http.StatusForbidden,
				err.Error(),
				"USER_INACTIVE",
			)

		case "technician profile already exists":
			response.Error(
				c,
				http.StatusConflict,
				err.Error(),
				"TECHNICIAN_ALREADY_EXISTS",
			)

		default:
			response.Error(
				c,
				http.StatusInternalServerError,
				"failed to create technician profile",
				"INTERNAL_SERVER_ERROR",
			)
		}

		return
	}

	response.Success(
		c,
		http.StatusCreated,
		"technician profile created successfully",
		TechnicianResponse{
			ID:             technician.ID,
			UserID:         technician.UserID,
			Name:           technician.Name,
			Email:          technician.Email,
			Phone:          technician.Phone,
			Specialization: technician.Specialization,
			IsAvailable:    technician.IsAvailable,
			IsActive:       technician.IsActive,
		},
	)
}

func (h *TechnicianHandler) GetMe(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		response.Error(
			c,
			http.StatusUnauthorized,
			"unauthorized",
			"UNAUTHORIZED",
		)
		return
	}

	technician, err := h.technicianService.GetByUserID(userID)
	if err != nil {
		switch err.Error() {
		case "technician profile not found":
			response.Error(
				c,
				http.StatusNotFound,
				err.Error(),
				"TECHNICIAN_NOT_FOUND",
			)

		case "technician profile is inactive":
			response.Error(
				c,
				http.StatusForbidden,
				err.Error(),
				"TECHNICIAN_INACTIVE",
			)

		default:
			response.Error(
				c,
				http.StatusInternalServerError,
				"failed to get technician profile",
				"INTERNAL_SERVER_ERROR",
			)
		}

		return
	}

	response.Success(
		c,
		http.StatusOK,
		"technician profile retrieved successfully",
		TechnicianResponse{
			ID:             technician.ID,
			UserID:         technician.UserID,
			Name:           technician.Name,
			Email:          technician.Email,
			Phone:          technician.Phone,
			Specialization: technician.Specialization,
			IsAvailable:    technician.IsAvailable,
			IsActive:       technician.IsActive,
		},
	)
}

func (h *TechnicianHandler) GetAvailable(c *gin.Context) {
	technicians, err := h.technicianService.GetAvailable()
	if err != nil {
		response.Error(
			c,
			http.StatusInternalServerError,
			"failed to retrieve technicians",
			"INTERNAL_SERVER_ERROR",
		)
		return
	}

	result := make([]TechnicianResponse, 0, len(technicians))

	for _, technician := range technicians {
		result = append(result, TechnicianResponse{
			ID:             technician.ID,
			UserID:         technician.UserID,
			Name:           technician.Name,
			Email:          technician.Email,
			Phone:          technician.Phone,
			Specialization: technician.Specialization,
			IsAvailable:    technician.IsAvailable,
			IsActive:       technician.IsActive,
		})
	}

	response.Success(
		c,
		http.StatusOK,
		"available technicians retrieved successfully",
		result,
	)
}
