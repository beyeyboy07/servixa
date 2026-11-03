package router

import (
	"servixa/internal/config"
	"servixa/internal/handler"
	"servixa/internal/middleware"
	"servixa/internal/model"

	"github.com/gin-gonic/gin"
)

func SetupRouter(
	healthHandler *handler.HealthHandler,
	authHandler *handler.AuthHandler,
	customerHandler *handler.CustomerHandler,
	technicianHandler *handler.TechnicianHandler,
	vehicleHandler *handler.VehicleHandler,
	partHandler *handler.PartHandler,
	serviceRequestHandler *handler.ServiceRequestHandler,
	workOrderHandler *handler.WorkOrderHandler,
	workOrderAssignmentHandler *handler.WorkOrderAssignmentHandler,
	workOrderVehicleHandler *handler.WorkOrderVehicleHandler,
	workOrderScheduleHandler *handler.WorkOrderScheduleHandler,
	serviceVisitHandler *handler.ServiceVisitHandler,
	workCompletionHandler *handler.WorkCompletionHandler,
	customerConfirmationHandler *handler.CustomerConfirmationHandler,
	notificationHandler *handler.NotificationHandler,
	cfg config.Config,
) *gin.Engine {

	router := gin.Default()
	router.Use(CORSMiddleware())

	// Health
	router.GET("/health", healthHandler.Check)

	// Auth
	api := router.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.Refresh)
			auth.POST("/logout", authHandler.Logout)
		}
	}

	// Protected routes
	protected := router.Group("/api/v1")
	protected.Use(middleware.Auth(cfg))
	{
		// Customer - CUSTOMER
		customer := protected.Group("/customers")
		customer.Use(middleware.RequireRoles(model.RoleCustomer))
		{
			customer.POST("", customerHandler.Create)
			customer.GET("/me", customerHandler.GetMe)
		}

		// Vehicle - ADMIN / DISPATCHER
		vehicle := protected.Group("/vehicles")
		vehicle.Use(
			middleware.RequireRoles(
				model.RoleAdmin,
				model.RoleDispatcher,
			),
		)
		{
			vehicle.POST("", vehicleHandler.Create)
			vehicle.GET("/all", vehicleHandler.GetAll)
			vehicle.GET("", vehicleHandler.GetAvailable)

			vehicle.PATCH("/:id/activate", vehicleHandler.Activate)
			vehicle.PATCH("/:id/deactivate", vehicleHandler.Deactivate)

			vehicle.GET("/:id", vehicleHandler.GetByID)
			vehicle.PUT("/:id", vehicleHandler.Update)
		}

		// Parts - ADMIN / DISPATCHER
		part := protected.Group("/parts")
		part.Use(
			middleware.RequireRoles(
				model.RoleAdmin,
				model.RoleDispatcher,
			),
		)
		{
			part.POST("", partHandler.Create)
			part.GET("/:id", partHandler.GetByID)
			part.PUT("/:id", partHandler.Update)
		}

		// Service Requests - ALL ALLOWED ROLES
		// GET uses one handler to determine data based on role.
		serviceRequest := protected.Group("/service-requests")
		serviceRequest.Use(
			middleware.RequireRoles(
				model.RoleCustomer,
				model.RoleAdmin,
				model.RoleDispatcher,
			),
		)
		{
			serviceRequest.GET("", serviceRequestHandler.GetAll)
		}

		// Service Requests - CUSTOMER
		serviceRequestCustomer := protected.Group("/service-requests")
		serviceRequestCustomer.Use(
			middleware.RequireRoles(model.RoleCustomer),
		)
		{
			serviceRequestCustomer.POST("", serviceRequestHandler.Create)
			serviceRequestCustomer.GET("/:id", serviceRequestHandler.GetByID)
		}

		// Service Request Management - ADMIN / DISPATCHER
		serviceRequestManagement := protected.Group("/service-requests")
		serviceRequestManagement.Use(
			middleware.RequireRoles(
				model.RoleAdmin,
				model.RoleDispatcher,
			),
		)
		{
			serviceRequestManagement.PATCH(
				"/:id/review",
				serviceRequestHandler.Review,
			)

			serviceRequestManagement.PATCH(
				"/:id/reject",
				serviceRequestHandler.Reject,
			)
		}

		// Work Orders - ALL AUTHENTICATED ROLES (READ)
		workOrderRead := protected.Group("/work-orders")
		workOrderRead.Use(
			middleware.RequireRoles(
				model.RoleAdmin,
				model.RoleDispatcher,
				model.RoleTechnician,
				model.RoleCustomer,
			),
		)
		{
			workOrderRead.GET("", workOrderHandler.GetAll)
		}

		// Work Order Management - ADMIN / DISPATCHER
		workOrder := protected.Group("/work-orders")
		workOrder.Use(
			middleware.RequireRoles(
				model.RoleAdmin,
				model.RoleDispatcher,
			),
		)
		{
			workOrder.POST("", workOrderHandler.Create)
			workOrder.GET("/:id", workOrderHandler.GetByID)

			workOrder.POST(
				"/:id/assign-technician",
				workOrderAssignmentHandler.AssignTechnician,
			)

			workOrder.POST(
				"/:id/assign-vehicle",
				workOrderVehicleHandler.AssignVehicle,
			)

			workOrder.POST(
				"/:id/schedule",
				workOrderScheduleHandler.Schedule,
			)

			workOrder.PATCH(
				"/:id/cancel",
				workOrderHandler.Cancel,
			)

			workOrder.PATCH(
				"/:id/reassign-technician",
				workOrderAssignmentHandler.ReassignTechnician,
			)
		}

		// Work Order Completion - TECHNICIAN
		workOrderCompletion := protected.Group("/work-orders")
		workOrderCompletion.Use(
			middleware.RequireRoles(model.RoleTechnician),
		)
		{
			workOrderCompletion.POST(
				"/:id/completion",
				workCompletionHandler.Create,
			)
		}

		// Technician - TECHNICIAN
		technician := protected.Group("/technicians")
		technician.Use(
			middleware.RequireRoles(model.RoleTechnician),
		)
		{
			technician.POST("", technicianHandler.Create)
			technician.GET("/me", technicianHandler.GetMe)
		}

		// Technician Operations - TECHNICIAN
		technicianOperations := protected.Group("")
		technicianOperations.Use(
			middleware.RequireRoles(model.RoleTechnician),
		)
		{
			technicianOperations.POST(
				"/work-orders/:id/visit",
				serviceVisitHandler.Create,
			)

			technicianOperations.POST(
				"/service-visits/:id/start",
				serviceVisitHandler.Start,
			)

			technicianOperations.POST(
				"/service-visits/:id/complete",
				serviceVisitHandler.Complete,
			)
		}

		// Customer Confirmation - CUSTOMER
		customerConfirmation := protected.Group("/work-orders")
		customerConfirmation.Use(
			middleware.RequireRoles(model.RoleCustomer),
		)
		{
			customerConfirmation.POST(
				"/:id/confirmation",
				customerConfirmationHandler.Confirm,
			)
		}

		// Notifications - ALL AUTHENTICATED ROLES
		notification := protected.Group("/notifications")
		{
			notification.GET("", notificationHandler.GetByUserID)
			notification.GET(
				"/unread",
				notificationHandler.GetUnreadByUserID,
			)
			notification.GET(
				"/unread/count",
				notificationHandler.CountUnreadByUserID,
			)
			notification.PATCH(
				"/read-all",
				notificationHandler.MarkAllAsRead,
			)
			notification.PATCH(
				"/:id/read",
				notificationHandler.MarkAsRead,
			)
		}

		// Technician Management - ADMIN / DISPATCHER
		technicianManagement := protected.Group("/technicians")
		technicianManagement.Use(
			middleware.RequireRoles(
				model.RoleAdmin,
				model.RoleDispatcher,
			),
		)
		{
			technicianManagement.GET(
				"",
				technicianHandler.GetAvailable,
			)
		}

		// Customer Management - ADMIN / DISPATCHER
		customerManagement := protected.Group("/customers")
		customerManagement.Use(
			middleware.RequireRoles(
				model.RoleAdmin,
				model.RoleDispatcher,
			),
		)
		{
			customerManagement.POST(
				"/manage",
				customerHandler.CreateByAdmin,
			)

			customerManagement.GET(
				"",
				customerHandler.GetAll,
			)

			customerManagement.GET(
				"/:id",
				customerHandler.GetByID,
			)
		}
	}

	return router
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		if origin == "http://localhost:5173" {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Vary", "Origin")
			c.Header(
				"Access-Control-Allow-Methods",
				"GET, POST, PUT, PATCH, DELETE, OPTIONS",
			)
			c.Header(
				"Access-Control-Allow-Headers",
				"Origin, Content-Type, Accept, Authorization",
			)
		}

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
