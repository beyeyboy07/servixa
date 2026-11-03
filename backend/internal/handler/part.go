package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"servixa/internal/middleware"
	"servixa/internal/response"
	"servixa/internal/service"
)

type PartHandler struct {
	partService *service.PartService
}

func NewPartHandler(
	partService *service.PartService,
) *PartHandler {
	return &PartHandler{
		partService: partService,
	}
}

type CreatePartRequest struct {
	Code        string  `json:"code"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Unit        string  `json:"unit"`
	Stock       int     `json:"stock"`
	MinStock    int     `json:"min_stock"`
	UnitPrice   float64 `json:"unit_price"`
}

type PartResponse struct {
	ID          uint    `json:"id"`
	Code        string  `json:"code"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Unit        string  `json:"unit"`
	Stock       int     `json:"stock"`
	MinStock    int     `json:"min_stock"`
	UnitPrice   float64 `json:"unit_price"`
	IsActive    bool    `json:"is_active"`
}

func (h *PartHandler) Create(c *gin.Context) {
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

	var request CreatePartRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		response.Error(
			c,
			http.StatusBadRequest,
			"invalid request body",
			"VALIDATION_ERROR",
		)
		return
	}

	part, err := h.partService.Create(
		request.Code,
		request.Name,
		request.Description,
		request.Unit,
		request.Stock,
		request.MinStock,
		request.UnitPrice,
	)

	if err != nil {
		switch err.Error() {
		case "code is required",
			"name is required",
			"unit is required",
			"stock cannot be negative",
			"minimum stock cannot be negative",
			"unit price cannot be negative":

			response.Error(
				c,
				http.StatusBadRequest,
				err.Error(),
				"VALIDATION_ERROR",
			)

		case "part with code already exists":
			response.Error(
				c,
				http.StatusConflict,
				err.Error(),
				"PART_ALREADY_EXISTS",
			)

		default:
			response.Error(
				c,
				http.StatusInternalServerError,
				"failed to create part",
				"INTERNAL_SERVER_ERROR",
			)
		}

		return
	}

	response.Success(
		c,
		http.StatusCreated,
		"part created successfully",
		PartResponse{
			ID:          part.ID,
			Code:        part.Code,
			Name:        part.Name,
			Description: part.Description,
			Unit:        part.Unit,
			Stock:       part.Stock,
			MinStock:    part.MinStock,
			UnitPrice:   part.UnitPrice,
			IsActive:    part.IsActive,
		},
	)
}

func (h *PartHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		response.Error(
			c,
			http.StatusBadRequest,
			"invalid part id",
			"VALIDATION_ERROR",
		)
		return
	}

	part, err := h.partService.GetByID(uint(id))
	if err != nil {
		switch err.Error() {
		case "part not found":
			response.Error(
				c,
				http.StatusNotFound,
				err.Error(),
				"PART_NOT_FOUND",
			)

		case "part is inactive":
			response.Error(
				c,
				http.StatusForbidden,
				err.Error(),
				"PART_INACTIVE",
			)

		default:
			response.Error(
				c,
				http.StatusInternalServerError,
				"failed to get part",
				"INTERNAL_SERVER_ERROR",
			)
		}

		return
	}

	response.Success(
		c,
		http.StatusOK,
		"part retrieved successfully",
		PartResponse{
			ID:          part.ID,
			Code:        part.Code,
			Name:        part.Name,
			Description: part.Description,
			Unit:        part.Unit,
			Stock:       part.Stock,
			MinStock:    part.MinStock,
			UnitPrice:   part.UnitPrice,
			IsActive:    part.IsActive,
		},
	)
}

func (h *PartHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		response.Error(
			c,
			http.StatusBadRequest,
			"invalid part id",
			"VALIDATION_ERROR",
		)
		return
	}

	var request struct {
		Code        string  `json:"code"`
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Unit        string  `json:"unit"`
		Stock       int     `json:"stock"`
		MinStock    int     `json:"min_stock"`
		UnitPrice   float64 `json:"unit_price"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		response.Error(
			c,
			http.StatusBadRequest,
			"invalid request body",
			"VALIDATION_ERROR",
		)
		return
	}

	part, err := h.partService.Update(
		uint(id),
		request.Code,
		request.Name,
		request.Description,
		request.Unit,
		request.Stock,
		request.MinStock,
		request.UnitPrice,
	)

	if err != nil {
		switch err.Error() {
		case "part not found":
			response.Error(
				c,
				http.StatusNotFound,
				err.Error(),
				"PART_NOT_FOUND",
			)

		case "part is inactive":
			response.Error(
				c,
				http.StatusForbidden,
				err.Error(),
				"PART_INACTIVE",
			)

		case "code is required",
			"name is required",
			"unit is required",
			"stock cannot be negative",
			"minimum stock cannot be negative",
			"unit price cannot be negative":

			response.Error(
				c,
				http.StatusBadRequest,
				err.Error(),
				"VALIDATION_ERROR",
			)

		case "part with code already exists":
			response.Error(
				c,
				http.StatusConflict,
				err.Error(),
				"PART_CODE_ALREADY_EXISTS",
			)

		default:
			response.Error(
				c,
				http.StatusInternalServerError,
				"failed to update part",
				"INTERNAL_SERVER_ERROR",
			)
		}

		return
	}

	response.Success(
		c,
		http.StatusOK,
		"part updated successfully",
		PartResponse{
			ID:          part.ID,
			Code:        part.Code,
			Name:        part.Name,
			Description: part.Description,
			Unit:        part.Unit,
			Stock:       part.Stock,
			MinStock:    part.MinStock,
			UnitPrice:   part.UnitPrice,
			IsActive:    part.IsActive,
		},
	)
}
