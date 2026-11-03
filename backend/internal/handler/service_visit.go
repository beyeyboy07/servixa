package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"servixa/internal/middleware"
	"servixa/internal/model"
	"servixa/internal/response"
	"servixa/internal/service"
)

type ServiceVisitHandler struct {
	service *service.ServiceVisitService
}

type CompleteServiceVisitRequest struct {
	Notes string `json:"notes"`
}

func NewServiceVisitHandler(
	service *service.ServiceVisitService,
) *ServiceVisitHandler {
	return &ServiceVisitHandler{
		service: service,
	}
}

type ServiceVisitResponse struct {
	ID           uint               `json:"id"`
	WorkOrderID  uint               `json:"work_order_id"`
	TechnicianID uint               `json:"technician_id"`
	Status       serviceVisitStatus `json:"status"`
	StartedAt    *time.Time         `json:"started_at"`
	CompletedAt  *time.Time         `json:"completed_at"`
	Notes        string             `json:"notes"`
	CreatedAt    time.Time          `json:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at"`
}

type serviceVisitStatus = model.ServiceVisitStatus

func toServiceVisitResponse(
	visit *model.ServiceVisit,
) ServiceVisitResponse {
	return ServiceVisitResponse{
		ID:           visit.ID,
		WorkOrderID:  visit.WorkOrderID,
		TechnicianID: visit.TechnicianID,
		Status:       visit.Status,
		StartedAt:    visit.StartedAt,
		CompletedAt:  visit.CompletedAt,
		Notes:        visit.Notes,
		CreatedAt:    visit.CreatedAt,
		UpdatedAt:    visit.UpdatedAt,
	}
}

func (h *ServiceVisitHandler) Create(c *gin.Context) {
	workOrderID, err := strconv.ParseUint(
		c.Param("id"),
		10,
		64,
	)

	if err != nil || workOrderID == 0 {
		response.Error(
			c,
			http.StatusBadRequest,
			"invalid work order ID",
			"INVALID_WORK_ORDER_ID",
		)
		return
	}

	technicianID, err := middleware.GetUserID(c)

	if err != nil {
		response.Error(
			c,
			http.StatusUnauthorized,
			"unauthorized",
			"UNAUTHORIZED",
		)
		return
	}

	visit, err := h.service.Create(
		uint(workOrderID),
		technicianID,
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

		case "technician is not assigned to this work order",
			"work order must be scheduled",
			"service visit already exists":
			response.Error(
				c,
				http.StatusConflict,
				err.Error(),
				"SERVICE_VISIT_CONFLICT",
			)

		default:
			response.Error(
				c,
				http.StatusInternalServerError,
				"failed to create service visit",
				"INTERNAL_ERROR",
			)
		}

		return
	}

	response.Success(
		c,
		http.StatusCreated,
		"service visit created successfully",
		toServiceVisitResponse(visit),
	)
}

func (h *ServiceVisitHandler) Start(c *gin.Context) {
	visitID, err := strconv.ParseUint(
		c.Param("id"),
		10,
		64,
	)

	if err != nil || visitID == 0 {
		response.Error(
			c,
			http.StatusBadRequest,
			"invalid service visit ID",
			"INVALID_SERVICE_VISIT_ID",
		)
		return
	}

	technicianID, err := middleware.GetUserID(c)

	if err != nil {
		response.Error(
			c,
			http.StatusUnauthorized,
			"unauthorized",
			"UNAUTHORIZED",
		)
		return
	}

	visit, err := h.service.Start(
		uint(visitID),
		technicianID,
	)

	if err != nil {
		switch err.Error() {
		case "service visit not found":
			response.Error(
				c,
				http.StatusNotFound,
				err.Error(),
				"SERVICE_VISIT_NOT_FOUND",
			)

		case "technician is not assigned to this visit",
			"service visit cannot be started":
			response.Error(
				c,
				http.StatusConflict,
				err.Error(),
				"SERVICE_VISIT_CONFLICT",
			)

		default:
			response.Error(
				c,
				http.StatusInternalServerError,
				"failed to start service visit",
				"INTERNAL_ERROR",
			)
		}

		return
	}

	response.Success(
		c,
		http.StatusOK,
		"service visit started successfully",
		toServiceVisitResponse(visit),
	)
}

func (h *ServiceVisitHandler) Complete(c *gin.Context) {
	visitID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || visitID == 0 {
		response.Error(
			c,
			http.StatusBadRequest,
			"invalid service visit id",
			"INVALID_ID",
		)
		return
	}

	technicianID, err := middleware.GetUserID(c)
	if err != nil {
		response.Error(
			c,
			http.StatusUnauthorized,
			"unauthorized",
			"UNAUTHORIZED",
		)
		return
	}

	var req CompleteServiceVisitRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(
			c,
			http.StatusBadRequest,
			"invalid request body",
			"INVALID_REQUEST",
		)
		return
	}

	visit, err := h.service.Complete(
		uint(visitID),
		technicianID,
		req.Notes,
	)

	if err != nil {
		switch err.Error() {
		case "service visit not found":
			response.Error(c, http.StatusNotFound, err.Error(), "NOT_FOUND")

		case "technician is not assigned to this visit":
			response.Error(c, http.StatusForbidden, err.Error(), "FORBIDDEN")

		case "service visit cannot be completed",
			"work order cannot be completed":
			response.Error(c, http.StatusConflict, err.Error(), "CONFLICT")

		default:
			response.Error(
				c,
				http.StatusInternalServerError,
				"failed to complete service visit",
				"SERVICE_VISIT_COMPLETE_FAILED",
			)
		}

		return
	}

	response.Success(
		c,
		http.StatusOK,
		"service visit completed successfully",
		toServiceVisitResponse(visit),
	)
}
