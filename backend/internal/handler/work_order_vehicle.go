package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"servixa/internal/response"
	"servixa/internal/service"
)

type WorkOrderVehicleHandler struct {
	service *service.WorkOrderVehicleService
}

func NewWorkOrderVehicleHandler(
	service *service.WorkOrderVehicleService,
) *WorkOrderVehicleHandler {
	return &WorkOrderVehicleHandler{
		service: service,
	}
}

type AssignVehicleRequest struct {
	VehicleID uint `json:"vehicle_id"`
}

type WorkOrderVehicleResponse struct {
	ID          uint      `json:"id"`
	WorkOrderID uint      `json:"work_order_id"`
	VehicleID   uint      `json:"vehicle_id"`
	AssignedAt  time.Time `json:"assigned_at"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (h *WorkOrderVehicleHandler) AssignVehicle(c *gin.Context) {
	workOrderID, err := strconv.ParseUint(c.Param("id"), 10, 64)

	if err != nil || workOrderID == 0 {
		response.Error(
			c,
			http.StatusBadRequest,
			"invalid work order ID",
			"INVALID_WORK_ORDER_ID",
		)
		return
	}

	var request AssignVehicleRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		response.Error(
			c,
			http.StatusBadRequest,
			"invalid request body",
			"INVALID_REQUEST",
		)
		return
	}

	if request.VehicleID == 0 {
		response.Error(
			c,
			http.StatusBadRequest,
			"vehicle_id is required",
			"INVALID_VEHICLE_ID",
		)
		return
	}

	assignment, err := h.service.Assign(
		uint(workOrderID),
		request.VehicleID,
	)

	if err != nil {
		switch err.Error() {
		case "work order not found":
			response.Error(
				c,
				http.StatusNotFound,
				err.Error(),
				"WORK_ORDER_NOT_FOUND",
			)

		case "vehicle not found":
			response.Error(
				c,
				http.StatusNotFound,
				err.Error(),
				"VEHICLE_NOT_FOUND",
			)

		case "work order cannot be assigned",
			"vehicle is inactive",
			"vehicle is not available",
			"work order already has a vehicle":
			response.Error(
				c,
				http.StatusConflict,
				err.Error(),
				"VEHICLE_ASSIGNMENT_CONFLICT",
			)

		default:
			response.Error(
				c,
				http.StatusInternalServerError,
				"failed to assign vehicle",
				"INTERNAL_ERROR",
			)
		}

		return
	}

	response.Success(
		c,
		http.StatusCreated,
		"vehicle assigned successfully",
		WorkOrderVehicleResponse{
			ID:          assignment.ID,
			WorkOrderID: assignment.WorkOrderID,
			VehicleID:   assignment.VehicleID,
			AssignedAt:  assignment.AssignedAt,
			CreatedAt:   assignment.CreatedAt,
			UpdatedAt:   assignment.UpdatedAt,
		},
	)
}
