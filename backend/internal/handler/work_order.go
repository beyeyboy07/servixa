package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"servixa/internal/middleware"
	"servixa/internal/response"
	"servixa/internal/service"
)

type WorkOrderHandler struct {
	workOrderService *service.WorkOrderService
}

func NewWorkOrderHandler(
	workOrderService *service.WorkOrderService,
) *WorkOrderHandler {
	return &WorkOrderHandler{
		workOrderService: workOrderService,
	}
}

type CreateWorkOrderRequest struct {
	ServiceRequestID uint `json:"service_request_id"`
}

type WorkOrderResponse struct {
	ID                uint                         `json:"id"`
	ServiceRequestID  uint                         `json:"service_request_id"`
	CustomerID        uint                         `json:"customer_id"`
	Title             string                       `json:"title"`
	Description       string                       `json:"description"`
	ServiceAddress    string                       `json:"service_address"`
	Status            string                       `json:"status"`
	ScheduledStart    *time.Time                   `json:"scheduled_start"`
	ScheduledEnd      *time.Time                   `json:"scheduled_end"`
	Assignment        *WorkOrderAssignmentResponse `json:"assignment"`
	VehicleAssignment *WorkOrderVehicleResponse    `json:"vehicle_assignment"`
	CreatedAt         time.Time                    `json:"created_at"`
	UpdatedAt         time.Time                    `json:"updated_at"`
}

func (h *WorkOrderHandler) Create(c *gin.Context) {
	var request CreateWorkOrderRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		response.Error(
			c,
			http.StatusBadRequest,
			"invalid request body",
			"VALIDATION_ERROR",
		)
		return
	}

	if request.ServiceRequestID == 0 {
		response.Error(
			c,
			http.StatusBadRequest,
			"service_request_id is required",
			"VALIDATION_ERROR",
		)
		return
	}

	workOrder, err := h.workOrderService.Create(
		request.ServiceRequestID,
	)

	if err != nil {
		switch err.Error() {
		case "service request not found":
			response.Error(
				c,
				http.StatusNotFound,
				err.Error(),
				"SERVICE_REQUEST_NOT_FOUND",
			)

		case "customer not found":
			response.Error(
				c,
				http.StatusNotFound,
				err.Error(),
				"CUSTOMER_NOT_FOUND",
			)

		case "work order already exists":
			response.Error(
				c,
				http.StatusConflict,
				err.Error(),
				"WORK_ORDER_ALREADY_EXISTS",
			)

		case "service request must be reviewed before creating work order":
			response.Error(
				c,
				http.StatusConflict,
				err.Error(),
				"INVALID_STATUS_TRANSITION",
			)

		default:
			response.Error(
				c,
				http.StatusInternalServerError,
				"failed to create work order",
				"INTERNAL_SERVER_ERROR",
			)
		}

		return
	}

	response.Success(
		c,
		http.StatusCreated,
		"work order created successfully",
		WorkOrderResponse{
			ID:               workOrder.ID,
			ServiceRequestID: workOrder.ServiceRequestID,
			CustomerID:       workOrder.CustomerID,
			Title:            workOrder.Title,
			Description:      workOrder.Description,
			ServiceAddress:   workOrder.ServiceAddress,
			Status:           string(workOrder.Status),
			ScheduledStart:   workOrder.ScheduledStart,
			ScheduledEnd:     workOrder.ScheduledEnd,
			CreatedAt:        workOrder.CreatedAt,
			UpdatedAt:        workOrder.UpdatedAt,
		},
	)
}

func (h *WorkOrderHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		response.Error(
			c,
			http.StatusBadRequest,
			"invalid work order id",
			"VALIDATION_ERROR",
		)
		return
	}

	workOrder, assignment, vehicleAssignment, err := h.workOrderService.GetByID(uint(id))
	if err != nil {
		if err.Error() == "work order not found" {
			response.Error(
				c,
				http.StatusNotFound,
				"work order not found",
				"WORK_ORDER_NOT_FOUND",
			)
			return
		}

		response.Error(
			c,
			http.StatusInternalServerError,
			"failed to retrieve work order",
			"INTERNAL_SERVER_ERROR",
		)
		return
	}

	var assignmentResponse *WorkOrderAssignmentResponse

	if assignment != nil {
		assignmentResponse = &WorkOrderAssignmentResponse{
			ID:           assignment.ID,
			WorkOrderID:  assignment.WorkOrderID,
			TechnicianID: assignment.TechnicianID,
			AssignedAt:   assignment.AssignedAt,
			CreatedAt:    assignment.CreatedAt,
			UpdatedAt:    assignment.UpdatedAt,
		}
	}

	var vehicleAssignmentResponse *WorkOrderVehicleResponse

	if vehicleAssignment != nil {
		vehicleAssignmentResponse = &WorkOrderVehicleResponse{
			ID:          vehicleAssignment.ID,
			WorkOrderID: vehicleAssignment.WorkOrderID,
			VehicleID:   vehicleAssignment.VehicleID,
			AssignedAt:  vehicleAssignment.AssignedAt,
			CreatedAt:   vehicleAssignment.CreatedAt,
			UpdatedAt:   vehicleAssignment.UpdatedAt,
		}
	}

	response.Success(
		c,
		http.StatusOK,
		"work order retrieved successfully",
		WorkOrderResponse{
			ID:                workOrder.ID,
			ServiceRequestID:  workOrder.ServiceRequestID,
			CustomerID:        workOrder.CustomerID,
			Title:             workOrder.Title,
			Description:       workOrder.Description,
			ServiceAddress:    workOrder.ServiceAddress,
			Status:            string(workOrder.Status),
			ScheduledStart:    workOrder.ScheduledStart,
			ScheduledEnd:      workOrder.ScheduledEnd,
			Assignment:        assignmentResponse,
			CreatedAt:         workOrder.CreatedAt,
			UpdatedAt:         workOrder.UpdatedAt,
			VehicleAssignment: vehicleAssignmentResponse,
		},
	)
}

func (h *WorkOrderHandler) GetAll(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "unauthorized", "UNAUTHORIZED")
		return
	}

	userRole, err := middleware.GetUserRole(c)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "unauthorized", "UNAUTHORIZED")
		return
	}

	workOrders, err := h.workOrderService.GetAll(userID, userRole)
	if err != nil {
		response.Error(
			c,
			http.StatusInternalServerError,
			"failed to retrieve work orders",
			"INTERNAL_SERVER_ERROR",
		)
		return
	}

	data := make([]WorkOrderResponse, 0, len(workOrders))

	for _, workOrder := range workOrders {
		_, assignment, vehicleAssignment, err :=
			h.workOrderService.GetByID(workOrder.ID)

		if err != nil {
			response.Error(
				c,
				http.StatusInternalServerError,
				"failed to retrieve work order assignments",
				"INTERNAL_SERVER_ERROR",
			)
			return
		}

		var assignmentResponse *WorkOrderAssignmentResponse

		if assignment != nil {
			assignmentResponse = &WorkOrderAssignmentResponse{
				ID:           assignment.ID,
				WorkOrderID:  assignment.WorkOrderID,
				TechnicianID: assignment.TechnicianID,
				AssignedAt:   assignment.AssignedAt,
				CreatedAt:    assignment.CreatedAt,
				UpdatedAt:    assignment.UpdatedAt,
			}
		}

		var vehicleAssignmentResponse *WorkOrderVehicleResponse

		if vehicleAssignment != nil {
			vehicleAssignmentResponse = &WorkOrderVehicleResponse{
				ID:          vehicleAssignment.ID,
				WorkOrderID: vehicleAssignment.WorkOrderID,
				VehicleID:   vehicleAssignment.VehicleID,
				AssignedAt:  vehicleAssignment.AssignedAt,
				CreatedAt:   vehicleAssignment.CreatedAt,
				UpdatedAt:   vehicleAssignment.UpdatedAt,
			}
		}

		data = append(data, WorkOrderResponse{
			ID:                workOrder.ID,
			ServiceRequestID:  workOrder.ServiceRequestID,
			CustomerID:        workOrder.CustomerID,
			Title:             workOrder.Title,
			Description:       workOrder.Description,
			ServiceAddress:    workOrder.ServiceAddress,
			Status:            string(workOrder.Status),
			ScheduledStart:    workOrder.ScheduledStart,
			ScheduledEnd:      workOrder.ScheduledEnd,
			Assignment:        assignmentResponse,
			VehicleAssignment: vehicleAssignmentResponse,
			CreatedAt:         workOrder.CreatedAt,
			UpdatedAt:         workOrder.UpdatedAt,
		})
	}

	response.Success(
		c,
		http.StatusOK,
		"work orders retrieved successfully",
		data,
	)
}

func (h *WorkOrderHandler) Cancel(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		response.Error(
			c,
			http.StatusBadRequest,
			"invalid work order id",
			"VALIDATION_ERROR",
		)
		return
	}

	workOrder, err := h.workOrderService.Cancel(uint(id))
	if err != nil {
		switch err.Error() {
		case "work order not found":
			response.Error(
				c,
				http.StatusNotFound,
				err.Error(),
				"WORK_ORDER_NOT_FOUND",
			)
		case "work order cannot be cancelled":
			response.Error(
				c,
				http.StatusConflict,
				err.Error(),
				"INVALID_STATUS_TRANSITION",
			)
		case "customer not found",
			"assigned technician not found":
			response.Error(
				c,
				http.StatusInternalServerError,
				err.Error(),
				"RELATED_RESOURCE_NOT_FOUND",
			)
		default:
			response.Error(
				c,
				http.StatusInternalServerError,
				"failed to cancel work order",
				"INTERNAL_SERVER_ERROR",
			)
		}
		return
	}

	response.Success(
		c,
		http.StatusOK,
		"work order cancelled successfully",
		WorkOrderResponse{
			ID:               workOrder.ID,
			ServiceRequestID: workOrder.ServiceRequestID,
			CustomerID:       workOrder.CustomerID,
			Title:            workOrder.Title,
			Description:      workOrder.Description,
			ServiceAddress:   workOrder.ServiceAddress,
			Status:           string(workOrder.Status),
			ScheduledStart:   workOrder.ScheduledStart,
			ScheduledEnd:     workOrder.ScheduledEnd,
			CreatedAt:        workOrder.CreatedAt,
			UpdatedAt:        workOrder.UpdatedAt,
		},
	)
}
