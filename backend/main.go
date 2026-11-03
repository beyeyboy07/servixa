package main

import (
	"log"

	"servixa/internal/config"
	"servixa/internal/database"
	"servixa/internal/handler"
	"servixa/internal/model"
	"servixa/internal/repository"
	"servixa/internal/router"
	"servixa/internal/service"
)

func main() {
	cfg := config.Load()

	db := database.Connect(cfg)

	if db == nil {
		log.Fatal("database connection is nil")
	}

	if err := db.AutoMigrate(
		&model.User{},
		&model.RefreshToken{},
		&model.Customer{},
		&model.Technician{},
		&model.Vehicle{},
		&model.Part{},
		&model.ServiceRequest{},
		&model.WorkOrder{},
		&model.WorkOrderAssignment{},
		&model.WorkOrderVehicle{},
		&model.WorkOrderSchedule{},
		&model.ServiceVisit{},
		&model.CustomerConfirmation{},
		&model.Notification{},
	); err != nil {
		log.Fatal("failed to migrate database:", err)
	}

	healthHandler := handler.NewHealthHandler()

	// Repositories
	userRepository := repository.NewUserRepository(db)
	customerRepository := repository.NewCustomerRepository(db)
	technicianRepository := repository.NewTechnicianRepository(db)
	vehicleRepository := repository.NewVehicleRepository(db)
	partRepository := repository.NewPartRepository(db)
	refreshTokenRepository := repository.NewRefreshTokenRepository(db)
	serviceRequestRepository := repository.NewServiceRequestRepository(db)
	workOrderRepository := repository.NewWorkOrderRepository(db)
	notificationRepository := repository.NewNotificationRepository(db)
	workOrderAssignmentRepository := repository.NewWorkOrderAssignmentRepository(db)
	workOrderVehicleRepository := repository.NewWorkOrderVehicleRepository(db)
	workOrderScheduleRepository := repository.NewWorkOrderScheduleRepository(db)
	serviceVisitRepository := repository.NewServiceVisitRepository(db)
	workCompletionRepository := repository.NewWorkCompletionRepository(db)
	customerConfirmationRepository := repository.NewCustomerConfirmationRepository(db)

	// Customer
	customerService := service.NewCustomerService(
		db,
		customerRepository,
		userRepository,
	)

	customerHandler := handler.NewCustomerHandler(
		customerService,
	)

	// Technician
	technicianService := service.NewTechnicianService(
		technicianRepository,
		userRepository,
	)

	technicianHandler := handler.NewTechnicianHandler(
		technicianService,
	)

	// Vehicle
	vehicleService := service.NewVehicleService(
		db,
		vehicleRepository,
	)
	vehicleHandler := handler.NewVehicleHandler(
		vehicleService,
	)

	// Authentication
	refreshTokenService := service.NewRefreshTokenService(
		refreshTokenRepository,
		cfg,
	)

	authService := service.NewAuthService(
		userRepository,
		refreshTokenService,
		cfg,
	)

	authHandler := handler.NewAuthHandler(
		authService,
	)

	// Part
	partService := service.NewPartService(
		partRepository,
	)

	partHandler := handler.NewPartHandler(
		partService,
	)

	// Service Request
	serviceRequestService := service.NewServiceRequestService(
		serviceRequestRepository,
		customerRepository,
		userRepository,
	)

	serviceRequestHandler := handler.NewServiceRequestHandler(
		serviceRequestService,
	)

	// Notification
	notificationService := service.NewNotificationService(
		db,
		notificationRepository,
	)

	notificationHandler := handler.NewNotificationHandler(
		notificationService,
	)

	// Work Order
	workOrderService := service.NewWorkOrderService(
		db,
		workOrderRepository,
		serviceRequestRepository,
		customerRepository,
		notificationRepository,
		workOrderAssignmentRepository,
		workOrderVehicleRepository,
		vehicleRepository,
		technicianRepository,
	)

	workOrderHandler := handler.NewWorkOrderHandler(
		workOrderService,
	)

	// Work Order Assignment
	workOrderAssignmentService := service.NewWorkOrderAssignmentService(
		db,
		workOrderAssignmentRepository,
		workOrderRepository,
		technicianRepository,
		notificationRepository,
	)

	workOrderAssignmentHandler := handler.NewWorkOrderAssignmentHandler(
		workOrderAssignmentService,
	)

	// Work Order Vehicle
	workOrderVehicleService := service.NewWorkOrderVehicleService(
		db,
		workOrderVehicleRepository,
		workOrderRepository,
		vehicleRepository,
	)

	workOrderVehicleHandler := handler.NewWorkOrderVehicleHandler(
		workOrderVehicleService,
	)

	// Work Order Schedule
	workOrderScheduleService := service.NewWorkOrderScheduleService(
		db,
		workOrderScheduleRepository,
		workOrderRepository,
		workOrderAssignmentRepository,
		workOrderVehicleRepository,
		technicianRepository,
		notificationRepository,
	)

	workOrderScheduleHandler := handler.NewWorkOrderScheduleHandler(
		workOrderScheduleService,
	)

	// Service Visit
	serviceVisitService := service.NewServiceVisitService(
		db,
		serviceVisitRepository,
		workOrderRepository,
		workOrderAssignmentRepository,
		technicianRepository,
		serviceRequestRepository,
		customerRepository,
		notificationRepository,
	)

	serviceVisitHandler := handler.NewServiceVisitHandler(
		serviceVisitService,
	)

	// Work Completion
	workCompletionService := service.NewWorkCompletionService(
		db,
		workCompletionRepository,
		workOrderRepository,
		workOrderVehicleRepository,
		vehicleRepository,
		workOrderAssignmentRepository,
		technicianRepository,
		serviceVisitRepository,
		serviceRequestRepository,
		customerRepository,
		notificationRepository,
	)

	workCompletionHandler := handler.NewWorkCompletionHandler(
		workCompletionService,
	)

	// Customer Confirmation
	customerConfirmationService := service.NewCustomerConfirmationService(
		db,
		customerConfirmationRepository,
		workOrderRepository,
		workCompletionRepository,
		customerRepository,
		workOrderAssignmentRepository,
		technicianRepository,
		notificationRepository,
	)

	customerConfirmationHandler := handler.NewCustomerConfirmationHandler(
		customerConfirmationService,
	)

	// Router
	r := router.SetupRouter(
		healthHandler,
		authHandler,
		customerHandler,
		technicianHandler,
		vehicleHandler,
		partHandler,
		serviceRequestHandler,
		workOrderHandler,
		workOrderAssignmentHandler,
		workOrderVehicleHandler,
		workOrderScheduleHandler,
		serviceVisitHandler,
		workCompletionHandler,
		customerConfirmationHandler,
		notificationHandler,
		cfg,
	)

	if err := r.Run(":" + cfg.AppPort); err != nil {
		log.Fatal(err)
	}
}
