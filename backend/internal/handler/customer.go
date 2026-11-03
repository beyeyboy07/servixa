package handler

import (
	"net/http"
	"strconv"

	"servixa/internal/middleware"
	"servixa/internal/response"
	"servixa/internal/service"

	"github.com/gin-gonic/gin"
)

type CustomerHandler struct {
	customerService *service.CustomerService
}

type AdminCreateCustomerRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Phone    string `json:"phone"`
	Address  string `json:"address"`
}

func NewCustomerHandler(customerService *service.CustomerService) *CustomerHandler {
	return &CustomerHandler{
		customerService: customerService,
	}
}

type CreateCustomerRequest struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
	Address string `json:"address"`
}

type CustomerResponse struct {
	ID       uint   `json:"id"`
	UserID   uint   `json:"user_id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Address  string `json:"address"`
	IsActive bool   `json:"is_active"`
}

func (h *CustomerHandler) Create(c *gin.Context) {
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

	var request CreateCustomerRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		response.Error(
			c,
			http.StatusBadRequest,
			"invalid request body",
			"INVALID_REQUEST",
		)
		return
	}

	customer, err := h.customerService.Create(
		userID,
		request.Name,
		request.Email,
		request.Phone,
		request.Address,
	)

	if err != nil {
		switch err.Error() {
		case "name is required",
			"email is required",
			"phone is required",
			"address is required":
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

		case "only customer users can create customer profile":
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

		case "customer profile already exists":
			response.Error(
				c,
				http.StatusConflict,
				err.Error(),
				"CUSTOMER_ALREADY_EXISTS",
			)

		default:
			response.Error(
				c,
				http.StatusInternalServerError,
				"failed to create customer profile",
				"INTERNAL_SERVER_ERROR",
			)
		}

		return
	}

	response.Success(
		c,
		http.StatusCreated,
		"customer profile created successfully",
		CustomerResponse{
			ID:       customer.ID,
			UserID:   customer.UserID,
			Name:     customer.Name,
			Email:    customer.Email,
			Phone:    customer.Phone,
			Address:  customer.Address,
			IsActive: customer.IsActive,
		},
	)
}

func (h *CustomerHandler) GetMe(c *gin.Context) {
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

	customer, err := h.customerService.GetByUserID(userID)
	if err != nil {
		switch err.Error() {
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
				"failed to get customer profile",
				"INTERNAL_SERVER_ERROR",
			)
		}

		return
	}

	response.Success(
		c,
		http.StatusOK,
		"customer profile retrieved successfully",
		CustomerResponse{
			ID:       customer.ID,
			UserID:   customer.UserID,
			Name:     customer.Name,
			Email:    customer.Email,
			Phone:    customer.Phone,
			Address:  customer.Address,
			IsActive: customer.IsActive,
		},
	)
}

func (h *CustomerHandler) GetAll(c *gin.Context) {
	customers, err := h.customerService.GetAll()
	if err != nil {
		response.Error(
			c,
			http.StatusInternalServerError,
			"failed to retrieve customers",
			"INTERNAL_SERVER_ERROR",
		)
		return
	}

	data := make([]CustomerResponse, 0, len(customers))

	for _, customer := range customers {
		data = append(data, CustomerResponse{
			ID:       customer.ID,
			UserID:   customer.UserID,
			Name:     customer.Name,
			Email:    customer.Email,
			Phone:    customer.Phone,
			Address:  customer.Address,
			IsActive: customer.IsActive,
		})
	}

	response.Success(
		c,
		http.StatusOK,
		"customers retrieved successfully",
		data,
	)
}

func (h *CustomerHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		response.Error(
			c,
			http.StatusBadRequest,
			"invalid customer id",
			"VALIDATION_ERROR",
		)
		return
	}

	customer, err := h.customerService.GetByID(uint(id))
	if err != nil {
		if err.Error() == "customer not found" {
			response.Error(
				c,
				http.StatusNotFound,
				"customer not found",
				"CUSTOMER_NOT_FOUND",
			)
			return
		}

		response.Error(
			c,
			http.StatusInternalServerError,
			"failed to retrieve customer",
			"INTERNAL_SERVER_ERROR",
		)
		return
	}

	response.Success(
		c,
		http.StatusOK,
		"customer retrieved successfully",
		CustomerResponse{
			ID:       customer.ID,
			UserID:   customer.UserID,
			Name:     customer.Name,
			Email:    customer.Email,
			Phone:    customer.Phone,
			Address:  customer.Address,
			IsActive: customer.IsActive,
		},
	)
}

func (h *CustomerHandler) CreateByAdmin(c *gin.Context) {
	var request AdminCreateCustomerRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		response.Error(
			c,
			http.StatusBadRequest,
			"invalid request body",
			"INVALID_REQUEST",
		)
		return
	}

	customer, err := h.customerService.CreateByAdmin(
		request.Name,
		request.Email,
		request.Password,
		request.Phone,
		request.Address,
	)

	if err != nil {
		switch err.Error() {
		case "name is required",
			"email is required",
			"password is required",
			"phone is required",
			"address is required",
			"password must be at least 8 characters":
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
				"failed to create customer",
				"INTERNAL_SERVER_ERROR",
			)
		}
		return
	}

	response.Success(
		c,
		http.StatusCreated,
		"customer created successfully",
		CustomerResponse{
			ID:       customer.ID,
			UserID:   customer.UserID,
			Name:     customer.Name,
			Email:    customer.Email,
			Phone:    customer.Phone,
			Address:  customer.Address,
			IsActive: customer.IsActive,
		},
	)
}
