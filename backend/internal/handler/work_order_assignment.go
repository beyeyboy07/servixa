package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"servixa/internal/response"
	"servixa/internal/service"
)

type WorkOrderAssignmentHandler struct {
	service *service.WorkOrderAssignmentService
}

func NewWorkOrderAssignmentHandler(
	service *service.WorkOrderAssignmentService,
) *WorkOrderAssignmentHandler {
	return &WorkOrderAssignmentHandler{
		service: service,
	}
}

type AssignTechnicianRequest struct {
	TechnicianID uint `json:"technician_id"`
}

type WorkOrderAssignmentResponse struct {
	ID           uint      `json:"id"`
	WorkOrderID  uint      `json:"work_order_id"`
	TechnicianID uint      `json:"technician_id"`
	AssignedAt   time.Time `json:"assigned_at"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (h *WorkOrderAssignmentHandler) AssignTechnician(c *gin.Context) {
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

	var request AssignTechnicianRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		response.Error(
			c,
			http.StatusBadRequest,
			"invalid request body",
			"INVALID_REQUEST",
		)
		return
	}

	if request.TechnicianID == 0 {
		response.Error(
			c,
			http.StatusBadRequest,
			"technician_id is required",
			"INVALID_TECHNICIAN_ID",
		)
		return
	}

	assignment, err := h.service.Assign(
		uint(workOrderID),
		request.TechnicianID,
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
		case "technician not found":
			response.Error(
				c,
				http.StatusNotFound,
				err.Error(),
				"TECHNICIAN_NOT_FOUND",
			)
		case "work order cannot be assigned",
			"technician is inactive",
			"technician is not available",
			"work order already has a technician":
			response.Error(
				c,
				http.StatusConflict,
				err.Error(),
				"ASSIGNMENT_CONFLICT",
			)
		default:
			response.Error(
				c,
				http.StatusInternalServerError,
				"failed to assign technician",
				"INTERNAL_ERROR",
			)
		}

		return
	}

	response.Success(
		c,
		http.StatusCreated,
		"technician assigned successfully",
		WorkOrderAssignmentResponse{
			ID:           assignment.ID,
			WorkOrderID:  assignment.WorkOrderID,
			TechnicianID: assignment.TechnicianID,
			AssignedAt:   assignment.AssignedAt,
			CreatedAt:    assignment.CreatedAt,
			UpdatedAt:    assignment.UpdatedAt,
		},
	)
}

func (h *WorkOrderAssignmentHandler) ReassignTechnician(c *gin.Context) {
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

	var request AssignTechnicianRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		response.Error(
			c,
			http.StatusBadRequest,
			"invalid request body",
			"INVALID_REQUEST",
		)
		return
	}

	if request.TechnicianID == 0 {
		response.Error(
			c,
			http.StatusBadRequest,
			"technician_id is required",
			"INVALID_TECHNICIAN_ID",
		)
		return
	}

	assignment, err := h.service.Reassign(
		uint(workOrderID),
		request.TechnicianID,
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
		case "technician not found",
			"previous technician not found":
			response.Error(
				c,
				http.StatusNotFound,
				err.Error(),
				"TECHNICIAN_NOT_FOUND",
			)
		case "work order has no technician assigned":
			response.Error(
				c,
				http.StatusConflict,
				err.Error(),
				"ASSIGNMENT_NOT_FOUND",
			)
		case "work order cannot be reassigned",
			"technician is already assigned to this work order",
			"technician is inactive",
			"technician is not available":
			response.Error(
				c,
				http.StatusConflict,
				err.Error(),
				"REASSIGNMENT_CONFLICT",
			)
		default:
			response.Error(
				c,
				http.StatusInternalServerError,
				"failed to reassign technician",
				"INTERNAL_ERROR",
			)
		}

		return
	}

	response.Success(
		c,
		http.StatusOK,
		"technician reassigned successfully",
		WorkOrderAssignmentResponse{
			ID:           assignment.ID,
			WorkOrderID:  assignment.WorkOrderID,
			TechnicianID: assignment.TechnicianID,
			AssignedAt:   assignment.AssignedAt,
			CreatedAt:    assignment.CreatedAt,
			UpdatedAt:    assignment.UpdatedAt,
		},
	)
}
