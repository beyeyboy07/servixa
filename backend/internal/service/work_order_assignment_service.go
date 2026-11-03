package service

import (
	"errors"
	"fmt"
	"time"

	"servixa/internal/model"
	"servixa/internal/repository"

	"gorm.io/gorm"
)

type WorkOrderAssignmentService struct {
	db                            *gorm.DB
	workOrderAssignmentRepository *repository.WorkOrderAssignmentRepository
	workOrderRepository           *repository.WorkOrderRepository
	technicianRepository          *repository.TechnicianRepository
	notificationRepository        *repository.NotificationRepository
}

func NewWorkOrderAssignmentService(
	db *gorm.DB,
	workOrderAssignmentRepository *repository.WorkOrderAssignmentRepository,
	workOrderRepository *repository.WorkOrderRepository,
	technicianRepository *repository.TechnicianRepository,
	notificationRepository *repository.NotificationRepository,

) *WorkOrderAssignmentService {
	return &WorkOrderAssignmentService{
		db:                            db,
		workOrderAssignmentRepository: workOrderAssignmentRepository,
		workOrderRepository:           workOrderRepository,
		technicianRepository:          technicianRepository,
		notificationRepository:        notificationRepository,
	}
}

func (s *WorkOrderAssignmentService) Assign(
	workOrderID uint,
	technicianID uint,
) (*model.WorkOrderAssignment, error) {

	var assignment *model.WorkOrderAssignment

	err := s.db.Transaction(func(tx *gorm.DB) error {
		txWorkOrderRepository :=
			s.workOrderRepository.WithTx(tx)

		txTechnicianRepository :=
			s.technicianRepository.WithTx(tx)

		txAssignmentRepository :=
			s.workOrderAssignmentRepository.WithTx(tx)

		workOrder, err :=
			txWorkOrderRepository.FindByIDForUpdate(workOrderID)

		if err != nil {
			return errors.New("work order not found")
		}

		if workOrder.Status == model.WorkOrderCompleted ||
			workOrder.Status == model.WorkOrderCustomerConfirmed ||
			workOrder.Status == model.WorkOrderClosed ||
			workOrder.Status == model.WorkOrderCancelled {
			return errors.New("work order cannot be assigned")
		}

		technician, err := txTechnicianRepository.FindByIDForUpdate(technicianID)
		if err != nil {
			return errors.New("technician not found")
		}

		fmt.Printf(
			"DEBUG technicianID=%d IsActive=%t IsAvailable=%t\n",
			technician.ID,
			technician.IsActive,
			technician.IsAvailable,
		)

		if !technician.IsActive {
			return errors.New("technician is inactive")
		}

		if !technician.IsAvailable {
			return errors.New("technician is not available")
		}

		existingAssignment, err :=
			txAssignmentRepository.FindByWorkOrderID(workOrderID)

		if err == nil && existingAssignment != nil {
			return errors.New("work order already has a technician")
		}

		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		assignment = &model.WorkOrderAssignment{
			WorkOrderID:  workOrderID,
			TechnicianID: technicianID,
			AssignedAt:   time.Now(),
		}

		if err := txAssignmentRepository.Create(assignment); err != nil {
			return err
		}

		workOrder.Status = model.WorkOrderAssigned

		if err := txWorkOrderRepository.Update(workOrder); err != nil {
			return err
		}

		txNotificationRepository :=
			s.notificationRepository.WithTx(tx)

		notification := &model.Notification{
			UserID:    technician.UserID,
			Type:      model.NotificationWorkOrderAssigned,
			Title:     "Work Order Assigned",
			Message:   "You have been assigned to work order: " + workOrder.Title,
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

	return assignment, nil
}

func (s *WorkOrderAssignmentService) Reassign(
	workOrderID uint,
	newTechnicianID uint,
) (*model.WorkOrderAssignment, error) {
	var assignment *model.WorkOrderAssignment

	err := s.db.Transaction(func(tx *gorm.DB) error {
		txWorkOrderRepository :=
			s.workOrderRepository.WithTx(tx)

		txTechnicianRepository :=
			s.technicianRepository.WithTx(tx)

		txAssignmentRepository :=
			s.workOrderAssignmentRepository.WithTx(tx)

		txNotificationRepository :=
			s.notificationRepository.WithTx(tx)

		// 1. Lock work order.
		workOrder, err :=
			txWorkOrderRepository.FindByIDForUpdate(workOrderID)
		if err != nil {
			return errors.New("work order not found")
		}

		// 2. Work order yang sudah final tidak boleh dipindahkan.
		if workOrder.Status == model.WorkOrderCompleted ||
			workOrder.Status == model.WorkOrderCustomerConfirmed ||
			workOrder.Status == model.WorkOrderClosed ||
			workOrder.Status == model.WorkOrderCancelled {
			return errors.New("work order cannot be reassigned")
		}

		// 3. Ambil assignment yang sudah ada.
		existingAssignment, err :=
			txAssignmentRepository.FindByWorkOrderID(workOrderID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("work order has no technician assigned")
			}
			return err
		}

		// 4. Tidak perlu reassignment ke teknisi yang sama.
		if existingAssignment.TechnicianID == newTechnicianID {
			return errors.New("technician is already assigned to this work order")
		}

		// 5. Lock dan validasi teknisi baru.
		newTechnician, err :=
			txTechnicianRepository.FindByIDForUpdate(newTechnicianID)
		if err != nil {
			return errors.New("technician not found")
		}

		if !newTechnician.IsActive {
			return errors.New("technician is inactive")
		}

		if !newTechnician.IsAvailable {
			return errors.New("technician is not available")
		}

		// 6. Ambil teknisi lama untuk notifikasi.
		oldTechnician, err :=
			txTechnicianRepository.FindByIDForUpdate(
				existingAssignment.TechnicianID,
			)
		if err != nil {
			return errors.New("previous technician not found")
		}

		// 7. Perbarui assignment yang sama.
		existingAssignment.TechnicianID = newTechnician.ID
		existingAssignment.AssignedAt = time.Now()

		if err := txAssignmentRepository.Update(existingAssignment); err != nil {
			return err
		}

		// 8. Status kembali menjadi ASSIGNED.
		workOrder.Status = model.WorkOrderAssigned

		if err := txWorkOrderRepository.Update(workOrder); err != nil {
			return err
		}

		// 9. Notifikasi kepada teknisi lama.
		oldNotification := &model.Notification{
			UserID:    oldTechnician.UserID,
			Type:      model.NotificationWorkOrderCancelled,
			Title:     "Work Order Reassigned",
			Message:   "You have been removed from work order: " + workOrder.Title,
			Reference: fmt.Sprintf("WORK_ORDER:%d", workOrder.ID),
			IsRead:    false,
		}

		if err := txNotificationRepository.Create(oldNotification); err != nil {
			return err
		}

		// 10. Notifikasi kepada teknisi baru.
		newNotification := &model.Notification{
			UserID:    newTechnician.UserID,
			Type:      model.NotificationWorkOrderAssigned,
			Title:     "Work Order Assigned",
			Message:   "You have been reassigned to work order: " + workOrder.Title,
			Reference: fmt.Sprintf("WORK_ORDER:%d", workOrder.ID),
			IsRead:    false,
		}

		if err := txNotificationRepository.Create(newNotification); err != nil {
			return err
		}

		assignment = existingAssignment
		return nil
	})

	if err != nil {
		return nil, err
	}

	return assignment, nil
}
