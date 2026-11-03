package repository

import (
	"time"

	"servixa/internal/model"

	"gorm.io/gorm"
)

type WorkOrderScheduleRepository struct {
	*BaseRepository
}

func NewWorkOrderScheduleRepository(
	db *gorm.DB,
) *WorkOrderScheduleRepository {
	return &WorkOrderScheduleRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

func (r *WorkOrderScheduleRepository) Create(
	schedule *model.WorkOrderSchedule,
) error {
	return r.db.Create(schedule).Error
}

func (r *WorkOrderScheduleRepository) FindByWorkOrderID(
	workOrderID uint,
) (*model.WorkOrderSchedule, error) {
	var schedule model.WorkOrderSchedule

	err := r.db.
		Where("work_order_id = ?", workOrderID).
		First(&schedule).
		Error

	if err != nil {
		return nil, err
	}

	return &schedule, nil
}

func (r *WorkOrderScheduleRepository) FindTechnicianOverlap(
	technicianID uint,
	start time.Time,
	end time.Time,
) ([]model.WorkOrderSchedule, error) {
	var schedules []model.WorkOrderSchedule

	err := r.db.
		Joins(
			"JOIN work_order_assignments ON work_order_assignments.work_order_id = work_order_schedules.work_order_id",
		).
		Where("work_order_assignments.technician_id = ?", technicianID).
		Where(
			"scheduled_start < ? AND scheduled_end > ?",
			end,
			start,
		).
		Find(&schedules).
		Error

	if err != nil {
		return nil, err
	}

	return schedules, nil
}

func (r *WorkOrderScheduleRepository) FindVehicleOverlap(
	vehicleID uint,
	start time.Time,
	end time.Time,
) ([]model.WorkOrderSchedule, error) {
	var schedules []model.WorkOrderSchedule

	err := r.db.
		Joins(
			"JOIN work_order_vehicles ON work_order_vehicles.work_order_id = work_order_schedules.work_order_id",
		).
		Where("work_order_vehicles.vehicle_id = ?", vehicleID).
		Where(
			"scheduled_start < ? AND scheduled_end > ?",
			end,
			start,
		).
		Find(&schedules).
		Error

	if err != nil {
		return nil, err
	}

	return schedules, nil
}

func (r *WorkOrderScheduleRepository) WithTx(
	tx *gorm.DB,
) *WorkOrderScheduleRepository {
	return NewWorkOrderScheduleRepository(tx)
}
