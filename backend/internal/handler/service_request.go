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

type ServiceRequestHandler struct {
	serviceRequestService *service.ServiceRequestService
}

func NewServiceRequestHandler(
	serviceRequestService *service.ServiceRequestService,
) *ServiceRequestHandler {
	return &ServiceRequestHandler{
		serviceRequestService: serviceRequestService,
	}
}

type CreateServiceRequestRequest struct {
	Title          string     `json:"title"`
	Description    string     `json:"description"`
	ServiceAddress string     `json:"service_address"`
	PreferredDate  *time.Time `json:"preferred_date"`
}

type ServiceRequestResponse struct {
	ID             uint       `json:"id"`
	CustomerID     uint       `json:"customer_id"`
	Title          string     `json:"title"`
	Description    string     `json:"description"`
	ServiceAddress string     `json:"service_address"`
	PreferredDate  *time.Time `json:"preferred_date"`
	Status         string     `json:"status"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

// Create - CUSTOMER
func (h *ServiceRequestHandler) Create(c *gin.Context) {
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

	var request CreateServiceRequestRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		response.Error(
			c,
			http.StatusBadRequest,
			"invalid request body",
			"VALIDATION_ERROR",
		)
		return
	}

	serviceRequest, err := h.serviceRequestService.Create(
		userID,
		request.Title,
		request.Description,
		request.ServiceAddress,
		request.PreferredDate,
	)

	if err != nil {
		switch err.Error() {
		case "title is required",
			"description is required",
			"service address is required":
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

		case "only customer users can create service request":
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

		case "customer profile not found":
			response.Error(
				c,
				http.StatusNotFound,
				err.Error(),
				"CUSTOMER_NOT_FOUND",
			)

		case "customer profile is inactive":
			response.Error(
				c,
				http.StatusForbidden,
				err.Error(),
				"CUSTOMER_INACTIVE",
			)

		default:
			response.Error(
				c,
				http.StatusInternalServerError,
				"failed to create service request",
				"INTERNAL_SERVER_ERROR",
			)
		}
		return
	}

	response.Success(
		c,
		http.StatusCreated,
		"service request created successfully",
		toServiceRequestResponse(serviceRequest),
	)
}

// GetAll - CUSTOMER / ADMIN / DISPATCHER
func (h *ServiceRequestHandler) GetAll(c *gin.Context) {
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

	userRole, err := middleware.GetUserRole(c)
	if err != nil {
		response.Error(
			c,
			http.StatusUnauthorized,
			"unauthorized",
			"UNAUTHORIZED",
		)
		return
	}

	var serviceRequests []model.ServiceRequest

	if userRole == model.RoleCustomer {
		serviceRequests, err = h.serviceRequestService.GetByCustomer(userID)
	} else {
		serviceRequests, err = h.serviceRequestService.GetAll(userID)
	}

	if err != nil {
		switch err.Error() {
		case "user not found":
			response.Error(
				c,
				http.StatusNotFound,
				err.Error(),
				"USER_NOT_FOUND",
			)

		case "user is inactive":
			response.Error(
				c,
				http.StatusForbidden,
				err.Error(),
				"USER_INACTIVE",
			)

		case "only customer users can access service requests",
			"only admin and dispatcher users can access service requests":
			response.Error(
				c,
				http.StatusForbidden,
				err.Error(),
				"FORBIDDEN",
			)

		case "customer profile not found":
			response.Error(
				c,
				http.StatusNotFound,
				err.Error(),
				"CUSTOMER_NOT_FOUND",
			)

		case "customer profile is inactive":
			response.Error(
				c,
				http.StatusForbidden,
				err.Error(),
				"CUSTOMER_INACTIVE",
			)

		default:
			response.Error(
				c,
				http.StatusInternalServerError,
				"failed to get service requests",
				"INTERNAL_SERVER_ERROR",
			)
		}
		return
	}

	data := make([]ServiceRequestResponse, 0, len(serviceRequests))

	for i := range serviceRequests {
		data = append(data, toServiceRequestResponse(&serviceRequests[i]))
	}

	response.Success(
		c,
		http.StatusOK,
		"service requests retrieved successfully",
		data,
	)
}

// GetByID - CUSTOMER (own request only)
func (h *ServiceRequestHandler) GetByID(c *gin.Context) {
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

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		response.Error(
			c,
			http.StatusBadRequest,
			"invalid service request id",
			"VALIDATION_ERROR",
		)
		return
	}

	serviceRequest, err := h.serviceRequestService.GetByID(
		userID,
		uint(id),
	)

	if err != nil {
		switch err.Error() {
		case "user not found":
			response.Error(
				c,
				http.StatusNotFound,
				err.Error(),
				"USER_NOT_FOUND",
			)

		case "only customer users can access service request":
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

		case "customer profile not found":
			response.Error(
				c,
				http.StatusNotFound,
				err.Error(),
				"CUSTOMER_NOT_FOUND",
			)

		case "customer profile is inactive":
			response.Error(
				c,
				http.StatusForbidden,
				err.Error(),
				"CUSTOMER_INACTIVE",
			)

		case "service request not found":
			response.Error(
				c,
				http.StatusNotFound,
				err.Error(),
				"SERVICE_REQUEST_NOT_FOUND",
			)

		default:
			response.Error(
				c,
				http.StatusInternalServerError,
				"failed to get service request",
				"INTERNAL_SERVER_ERROR",
			)
		}
		return
	}

	response.Success(
		c,
		http.StatusOK,
		"service request retrieved successfully",
		toServiceRequestResponse(serviceRequest),
	)
}

// Review - ADMIN / DISPATCHER
func (h *ServiceRequestHandler) Review(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		response.Error(
			c,
			http.StatusBadRequest,
			"invalid service request id",
			"VALIDATION_ERROR",
		)
		return
	}

	serviceRequest, err := h.serviceRequestService.Review(uint(id))

	if err != nil {
		switch err.Error() {
		case "service request not found":
			response.Error(
				c,
				http.StatusNotFound,
				err.Error(),
				"SERVICE_REQUEST_NOT_FOUND",
			)

		case "service request cannot be reviewed":
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
				"failed to review service request",
				"INTERNAL_SERVER_ERROR",
			)
		}
		return
	}

	response.Success(
		c,
		http.StatusOK,
		"service request reviewed successfully",
		toServiceRequestResponse(serviceRequest),
	)
}

// Reject - ADMIN / DISPATCHER
func (h *ServiceRequestHandler) Reject(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		response.Error(
			c,
			http.StatusBadRequest,
			"invalid service request id",
			"VALIDATION_ERROR",
		)
		return
	}

	serviceRequest, err := h.serviceRequestService.Reject(uint(id))

	if err != nil {
		switch err.Error() {
		case "service request not found":
			response.Error(
				c,
				http.StatusNotFound,
				err.Error(),
				"SERVICE_REQUEST_NOT_FOUND",
			)

		case "service request cannot be rejected":
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
				"failed to reject service request",
				"INTERNAL_SERVER_ERROR",
			)
		}
		return
	}

	response.Success(
		c,
		http.StatusOK,
		"service request rejected successfully",
		toServiceRequestResponse(serviceRequest),
	)
}

// Convert model to API response
func toServiceRequestResponse(
	serviceRequest *model.ServiceRequest,
) ServiceRequestResponse {
	return ServiceRequestResponse{
		ID:             serviceRequest.ID,
		CustomerID:     serviceRequest.CustomerID,
		Title:          serviceRequest.Title,
		Description:    serviceRequest.Description,
		ServiceAddress: serviceRequest.ServiceAddress,
		PreferredDate:  serviceRequest.PreferredDate,
		Status:         string(serviceRequest.Status),
		CreatedAt:      serviceRequest.CreatedAt,
		UpdatedAt:      serviceRequest.UpdatedAt,
	}
}
