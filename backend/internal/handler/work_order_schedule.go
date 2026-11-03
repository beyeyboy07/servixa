package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"servixa/internal/response"
	"servixa/internal/service"
)

type WorkOrderScheduleHandler struct {
	service *service.WorkOrderScheduleService
}

func NewWorkOrderScheduleHandler(
	service *service.WorkOrderScheduleService,
) *WorkOrderScheduleHandler {
	return &WorkOrderScheduleHandler{
		service: service,
	}
}

type ScheduleWorkOrderRequest struct {
	ScheduledStart time.Time `json:"scheduled_start"`
	ScheduledEnd   time.Time `json:"scheduled_end"`
}

type WorkOrderScheduleResponse struct {
	ID             uint      `json:"id"`
	WorkOrderID    uint      `json:"work_order_id"`
	ScheduledStart time.Time `json:"scheduled_start"`
	ScheduledEnd   time.Time `json:"scheduled_end"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func (h *WorkOrderScheduleHandler) Schedule(c *gin.Context) {
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

	var request ScheduleWorkOrderRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		response.Error(
			c,
			http.StatusBadRequest,
			"invalid request body",
			"INVALID_REQUEST",
		)
		return
	}

	schedule, err := h.service.Schedule(
		uint(workOrderID),
		request.ScheduledStart,
		request.ScheduledEnd,
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

		case "technician assignment not found":
			response.Error(
				c,
				http.StatusConflict,
				err.Error(),
				"TECHNICIAN_ASSIGNMENT_NOT_FOUND",
			)

		case "vehicle assignment not found":
			response.Error(
				c,
				http.StatusConflict,
				err.Error(),
				"VEHICLE_ASSIGNMENT_NOT_FOUND",
			)

		case "work order must have technician assigned",
			"work order already has a schedule",
			"technician has overlapping schedule",
			"vehicle has overlapping schedule",
			"scheduled start must be before scheduled end":
			response.Error(
				c,
				http.StatusConflict,
				err.Error(),
				"SCHEDULE_CONFLICT",
			)

		default:
			response.Error(
				c,
				http.StatusInternalServerError,
				"failed to schedule work order",
				"INTERNAL_ERROR",
			)
		}

		return
	}

	response.Success(
		c,
		http.StatusCreated,
		"work order scheduled successfully",
		WorkOrderScheduleResponse{
			ID:             schedule.ID,
			WorkOrderID:    schedule.WorkOrderID,
			ScheduledStart: schedule.ScheduledStart,
			ScheduledEnd:   schedule.ScheduledEnd,
			CreatedAt:      schedule.CreatedAt,
			UpdatedAt:      schedule.UpdatedAt,
		},
	)
}
