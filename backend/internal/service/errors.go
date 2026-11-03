package service

import "errors"

var (
	ErrTechnicianNotFound = errors.New("technician not found")
	ErrTechnicianInactive = errors.New("technician is inactive")

	ErrWorkOrderNotFound          = errors.New("work order not found")
	ErrWorkOrderMustBeScheduled   = errors.New("work order must be scheduled")
	ErrWorkOrderMustBeInProgress  = errors.New("work order must be in progress")
	ErrWorkOrderCannotBeStarted   = errors.New("work order cannot be started")
	ErrWorkOrderCannotBeCompleted = errors.New("work order cannot be completed")

	ErrTechnicianAssignmentNotFound     = errors.New("technician assignment not found")
	ErrTechnicianNotAssignedToWorkOrder = errors.New("technician is not assigned to this work order")
	ErrTechnicianNotAssignedToVisit     = errors.New("technician is not assigned to this visit")

	ErrServiceVisitNotFound          = errors.New("service visit not found")
	ErrServiceVisitAlreadyExists     = errors.New("service visit already exists")
	ErrServiceVisitCannotBeStarted   = errors.New("service visit cannot be started")
	ErrServiceVisitCannotBeCompleted = errors.New("service visit cannot be completed")
	ErrServiceVisitMustBeCompleted   = errors.New("service visit must be completed first")

	ErrWorkCompletionAlreadyExists   = errors.New("work completion already exists")
	ErrWorkCompletionSummaryRequired = errors.New("summary is required")
	ErrWorkPerformedRequired         = errors.New("work performed is required")

	ErrServiceRequestNotFound = errors.New("service request not found")
	ErrCustomerNotFound       = errors.New("customer not found")
)
