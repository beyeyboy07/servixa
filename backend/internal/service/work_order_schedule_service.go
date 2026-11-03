package service

import (
	"errors"
	"fmt"
	"time"

	"servixa/internal/model"
	"servixa/internal/repository"

	"gorm.io/gorm"
)

type WorkOrderScheduleService struct {
	db                      *gorm.DB
	scheduleRepository      *repository.WorkOrderScheduleRepository
	workOrderRepository     *repository.WorkOrderRepository
	workOrderAssignmentRepo *repository.WorkOrderAssignmentRepository
	workOrderVehicleRepo    *repository.WorkOrderVehicleRepository
	technicianRepository    *repository.TechnicianRepository
	notificationRepository  *repository.NotificationRepository
}

func NewWorkOrderScheduleService(
	db *gorm.DB,
	scheduleRepository *repository.WorkOrderScheduleRepository,
	workOrderRepository *repository.WorkOrderRepository,
	workOrderAssignmentRepo *repository.WorkOrderAssignmentRepository,
	workOrderVehicleRepo *repository.WorkOrderVehicleRepository,
	technicianRepository *repository.TechnicianRepository,
	notificationRepository *repository.NotificationRepository,
) *WorkOrderScheduleService {
	return &WorkOrderScheduleService{
		db:                      db,
		scheduleRepository:      scheduleRepository,
		workOrderRepository:     workOrderRepository,
		workOrderAssignmentRepo: workOrderAssignmentRepo,
		workOrderVehicleRepo:    workOrderVehicleRepo,
		technicianRepository:    technicianRepository,
		notificationRepository:  notificationRepository,
	}
}

func (s *WorkOrderScheduleService) Schedule(
	workOrderID uint,
	start time.Time,
	end time.Time,
) (*model.WorkOrderSchedule, error) {

	if !start.Before(end) {
		return nil, errors.New("scheduled start must be before scheduled end")
	}

	var schedule *model.WorkOrderSchedule

	err := s.db.Transaction(func(tx *gorm.DB) error {
		txScheduleRepository := s.scheduleRepository.WithTx(tx)
		txWorkOrderRepository := s.workOrderRepository.WithTx(tx)
		txAssignmentRepository := s.workOrderAssignmentRepo.WithTx(tx)
		txVehicleRepository := s.workOrderVehicleRepo.WithTx(tx)
		txTechnicianRepository := s.technicianRepository.WithTx(tx)
		txNotificationRepository := s.notificationRepository.WithTx(tx)

		workOrder, err := txWorkOrderRepository.FindByIDForUpdate(workOrderID)
		if err != nil {
			return errors.New("work order not found")
		}

		if workOrder.Status != model.WorkOrderAssigned {
			return errors.New("work order must have technician assigned")
		}

		existingSchedule, err :=
			txScheduleRepository.FindByWorkOrderID(workOrderID)

		if err == nil && existingSchedule != nil {
			return errors.New("work order already has a schedule")
		}

		assignment, err :=
			txAssignmentRepository.FindByWorkOrderID(workOrderID)
		if err != nil {
			return errors.New("technician assignment not found")
		}

		vehicleAssignment, err :=
			txVehicleRepository.FindByWorkOrderID(workOrderID)
		if err != nil {
			return errors.New("vehicle assignment not found")
		}

		technicianSchedules, err :=
			txScheduleRepository.FindTechnicianOverlap(
				assignment.TechnicianID,
				start,
				end,
			)
		if err != nil {
			return err
		}

		if len(technicianSchedules) > 0 {
			return errors.New("technician has overlapping schedule")
		}

		vehicleSchedules, err :=
			txScheduleRepository.FindVehicleOverlap(
				vehicleAssignment.VehicleID,
				start,
				end,
			)
		if err != nil {
			return err
		}

		if len(vehicleSchedules) > 0 {
			return errors.New("vehicle has overlapping schedule")
		}

		technician, err := txTechnicianRepository.FindByID(
			assignment.TechnicianID,
		)
		if err != nil {
			return errors.New("technician not found")
		}

		schedule = &model.WorkOrderSchedule{
			WorkOrderID:    workOrderID,
			ScheduledStart: start,
			ScheduledEnd:   end,
		}

		if err := txScheduleRepository.Create(schedule); err != nil {
			return err
		}

		workOrder.Status = model.WorkOrderScheduled

		if err := txWorkOrderRepository.Update(workOrder); err != nil {
			return err
		}

		notification := &model.Notification{
			UserID: technician.UserID,
			Type:   model.NotificationWorkOrderScheduled,
			Title:  "Work Order Scheduled",
			Message: fmt.Sprintf(
				"Work order %s has been scheduled for %s.",
				workOrder.Title,
				start.Format("02 Jan 2006 15:04"),
			),
			Reference: fmt.Sprintf("WORK_ORDER:%d", workOrder.ID),
			IsRead:    false,
		}

		if err := txNotificationRepository.Create(notification); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return schedule, nil
}
