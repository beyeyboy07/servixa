package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"servixa/internal/middleware"
	"servixa/internal/response"
	"servixa/internal/service"
)

type VehicleHandler struct {
	vehicleService *service.VehicleService
}

func NewVehicleHandler(
	vehicleService *service.VehicleService,
) *VehicleHandler {
	return &VehicleHandler{
		vehicleService: vehicleService,
	}
}

type CreateVehicleRequest struct {
	PlateNumber string `json:"plate_number"`
	Brand       string `json:"brand"`
	Model       string `json:"model"`
	Year        int    `json:"year"`
}

type VehicleResponse struct {
	ID          uint   `json:"id"`
	PlateNumber string `json:"plate_number"`
	Brand       string `json:"brand"`
	Model       string `json:"model"`
	Year        int    `json:"year"`
	Status      string `json:"status"`
	IsActive    bool   `json:"is_active"`
}

type UpdateVehicleRequest struct {
	PlateNumber string `json:"plate_number"`
	Brand       string `json:"brand"`
	Model       string `json:"model"`
	Year        int    `json:"year"`
	Status      string `json:"status"`
}

func (h *VehicleHandler) Create(c *gin.Context) {
	_, err := middleware.GetUserID(c)
	if err != nil {
		response.Error(
			c,
			http.StatusUnauthorized,
			"unauthorized",
			"UNAUTHORIZED",
		)
		return
	}

	var request CreateVehicleRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		response.Error(
			c,
			http.StatusBadRequest,
			"invalid request body",
			"VALIDATION_ERROR",
		)
		return
	}

	vehicle, err := h.vehicleService.Create(
		request.PlateNumber,
		request.Brand,
		request.Model,
		request.Year,
	)

	if err != nil {
		switch err.Error() {
		case "plate number is required",
			"brand is required",
			"model is required",
			"year must be greater than 0":

			response.Error(
				c,
				http.StatusBadRequest,
				err.Error(),
				"VALIDATION_ERROR",
			)

		case "vehicle with plate number already exists":
			response.Error(
				c,
				http.StatusConflict,
				err.Error(),
				"VEHICLE_ALREADY_EXISTS",
			)

		default:
			response.Error(
				c,
				http.StatusInternalServerError,
				"failed to create vehicle",
				"INTERNAL_SERVER_ERROR",
			)
		}

		return
	}

	response.Success(
		c,
		http.StatusCreated,
		"vehicle created successfully",
		VehicleResponse{
			ID:          vehicle.ID,
			PlateNumber: vehicle.PlateNumber,
			Brand:       vehicle.Brand,
			Model:       vehicle.Model,
			Year:        vehicle.Year,
			Status:      vehicle.Status,
			IsActive:    vehicle.IsActive,
		},
	)
}

func (h *VehicleHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		response.Error(
			c,
			http.StatusBadRequest,
			"invalid vehicle id",
			"VALIDATION_ERROR",
		)
		return
	}

	vehicle, err := h.vehicleService.GetByID(uint(id))
	if err != nil {
		switch err.Error() {
		case "vehicle not found":
			response.Error(
				c,
				http.StatusNotFound,
				err.Error(),
				"VEHICLE_NOT_FOUND",
			)

		case "vehicle is inactive":
			response.Error(
				c,
				http.StatusForbidden,
				err.Error(),
				"VEHICLE_INACTIVE",
			)

		default:
			response.Error(
				c,
				http.StatusInternalServerError,
				"failed to get vehicle",
				"INTERNAL_SERVER_ERROR",
			)
		}

		return
	}

	response.Success(
		c,
		http.StatusOK,
		"vehicle retrieved successfully",
		VehicleResponse{
			ID:          vehicle.ID,
			PlateNumber: vehicle.PlateNumber,
			Brand:       vehicle.Brand,
			Model:       vehicle.Model,
			Year:        vehicle.Year,
			Status:      vehicle.Status,
			IsActive:    vehicle.IsActive,
		},
	)
}

func (h *VehicleHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		response.Error(
			c,
			http.StatusBadRequest,
			"invalid vehicle id",
			"VALIDATION_ERROR",
		)
		return
	}

	var request UpdateVehicleRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		response.Error(
			c,
			http.StatusBadRequest,
			"invalid request body",
			"VALIDATION_ERROR",
		)
		return
	}

	vehicle, err := h.vehicleService.Update(
		uint(id),
		request.PlateNumber,
		request.Brand,
		request.Model,
		request.Year,
		request.Status,
	)

	if err != nil {
		switch err.Error() {
		case "plate number is required",
			"brand is required",
			"model is required",
			"year must be greater than 0",
			"status is required":

			response.Error(
				c,
				http.StatusBadRequest,
				err.Error(),
				"VALIDATION_ERROR",
			)

		case "vehicle not found":
			response.Error(
				c,
				http.StatusNotFound,
				err.Error(),
				"VEHICLE_NOT_FOUND",
			)

		case "vehicle with plate number already exists":
			response.Error(
				c,
				http.StatusConflict,
				err.Error(),
				"VEHICLE_ALREADY_EXISTS",
			)

		default:
			response.Error(
				c,
				http.StatusInternalServerError,
				"failed to update vehicle",
				"INTERNAL_SERVER_ERROR",
			)
		}

		return
	}

	response.Success(
		c,
		http.StatusOK,
		"vehicle updated successfully",
		VehicleResponse{
			ID:          vehicle.ID,
			PlateNumber: vehicle.PlateNumber,
			Brand:       vehicle.Brand,
			Model:       vehicle.Model,
			Year:        vehicle.Year,
			Status:      vehicle.Status,
			IsActive:    vehicle.IsActive,
		},
	)
}

func (h *VehicleHandler) GetAvailable(c *gin.Context) {
	vehicles, err := h.vehicleService.GetAvailable()
	if err != nil {
		response.Error(
			c,
			http.StatusInternalServerError,
			"failed to retrieve vehicles",
			"INTERNAL_SERVER_ERROR",
		)
		return
	}

	result := make([]VehicleResponse, 0, len(vehicles))

	for _, vehicle := range vehicles {
		result = append(result, VehicleResponse{
			ID:          vehicle.ID,
			PlateNumber: vehicle.PlateNumber,
			Brand:       vehicle.Brand,
			Model:       vehicle.Model,
			Year:        vehicle.Year,
			Status:      vehicle.Status,
			IsActive:    vehicle.IsActive,
		})
	}

	response.Success(
		c,
		http.StatusOK,
		"available vehicles retrieved successfully",
		result,
	)
}

func (h *VehicleHandler) GetAll(c *gin.Context) {
	vehicles, err := h.vehicleService.GetAll()
	if err != nil {
		response.Error(
			c,
			http.StatusInternalServerError,
			"failed to retrieve vehicles",
			"INTERNAL_SERVER_ERROR",
		)
		return
	}

	result := make([]VehicleResponse, 0, len(vehicles))

	for _, vehicle := range vehicles {
		result = append(result, VehicleResponse{
			ID:          vehicle.ID,
			PlateNumber: vehicle.PlateNumber,
			Brand:       vehicle.Brand,
			Model:       vehicle.Model,
			Year:        vehicle.Year,
			Status:      vehicle.Status,
			IsActive:    vehicle.IsActive,
		})
	}

	response.Success(
		c,
		http.StatusOK,
		"vehicles retrieved successfully",
		result,
	)
}

func (h *VehicleHandler) Deactivate(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		response.Error(
			c,
			http.StatusBadRequest,
			"invalid vehicle id",
			"VALIDATION_ERROR",
		)
		return
	}

	vehicle, err := h.vehicleService.Deactivate(uint(id))
	if err != nil {
		switch err.Error() {
		case "vehicle not found":
			response.Error(
				c,
				http.StatusNotFound,
				err.Error(),
				"VEHICLE_NOT_FOUND",
			)
		case "vehicle is already inactive":
			response.Error(
				c,
				http.StatusConflict,
				err.Error(),
				"VEHICLE_ALREADY_INACTIVE",
			)
		case "vehicle is currently in use":
			response.Error(
				c,
				http.StatusConflict,
				err.Error(),
				"VEHICLE_IN_USE",
			)
		default:
			response.Error(
				c,
				http.StatusInternalServerError,
				"failed to deactivate vehicle",
				"INTERNAL_SERVER_ERROR",
			)
		}
		return
	}

	response.Success(
		c,
		http.StatusOK,
		"vehicle deactivated successfully",
		VehicleResponse{
			ID:          vehicle.ID,
			PlateNumber: vehicle.PlateNumber,
			Brand:       vehicle.Brand,
			Model:       vehicle.Model,
			Year:        vehicle.Year,
			Status:      vehicle.Status,
			IsActive:    vehicle.IsActive,
		},
	)
}

func (h *VehicleHandler) Activate(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		response.Error(
			c,
			http.StatusBadRequest,
			"invalid vehicle id",
			"VALIDATION_ERROR",
		)
		return
	}

	vehicle, err := h.vehicleService.Activate(uint(id))
	if err != nil {
		switch err.Error() {
		case "vehicle not found":
			response.Error(
				c,
				http.StatusNotFound,
				err.Error(),
				"VEHICLE_NOT_FOUND",
			)
		case "vehicle is already active":
			response.Error(
				c,
				http.StatusConflict,
				err.Error(),
				"VEHICLE_ALREADY_ACTIVE",
			)
		default:
			response.Error(
				c,
				http.StatusInternalServerError,
				"failed to activate vehicle",
				"INTERNAL_SERVER_ERROR",
			)
		}
		return
	}

	response.Success(
		c,
		http.StatusOK,
		"vehicle activated successfully",
		VehicleResponse{
			ID:          vehicle.ID,
			PlateNumber: vehicle.PlateNumber,
			Brand:       vehicle.Brand,
			Model:       vehicle.Model,
			Year:        vehicle.Year,
			Status:      vehicle.Status,
			IsActive:    vehicle.IsActive,
		},
	)
}
