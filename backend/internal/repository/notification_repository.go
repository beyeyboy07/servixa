package repository

import (
	"servixa/internal/model"

	"gorm.io/gorm"
)

type NotificationRepository struct {
	*BaseRepository
}

func NewNotificationRepository(db *gorm.DB) *NotificationRepository {
	return &NotificationRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

func (r *NotificationRepository) Create(
	notification *model.Notification,
) error {
	return r.db.Create(notification).Error
}

func (r *NotificationRepository) FindByID(
	id uint,
) (*model.Notification, error) {
	var notification model.Notification

	err := r.db.
		Where("id = ?", id).
		First(&notification).
		Error

	if err != nil {
		return nil, err
	}

	return &notification, nil
}

func (r *NotificationRepository) FindByUserID(
	userID uint,
) ([]model.Notification, error) {
	var notifications []model.Notification

	err := r.db.
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&notifications).
		Error

	if err != nil {
		return nil, err
	}

	return notifications, nil
}

func (r *NotificationRepository) FindUnreadByUserID(
	userID uint,
) ([]model.Notification, error) {
	var notifications []model.Notification

	err := r.db.
		Where("user_id = ? AND is_read = ?", userID, false).
		Order("created_at DESC").
		Find(&notifications).
		Error

	if err != nil {
		return nil, err
	}

	return notifications, nil
}

func (r *NotificationRepository) CountUnreadByUserID(
	userID uint,
) (int64, error) {
	var count int64

	err := r.db.
		Model(&model.Notification{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Count(&count).
		Error

	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *NotificationRepository) MarkAsRead(
	id uint,
	userID uint,
) error {
	return r.db.
		Model(&model.Notification{}).
		Where("id = ? AND user_id = ? AND is_read = ?", id, userID, false).
		Updates(map[string]interface{}{
			"is_read": true,
			"read_at": gorm.Expr("CURRENT_TIMESTAMP"),
		}).
		Error
}

func (r *NotificationRepository) MarkAllAsRead(
	userID uint,
) error {
	return r.db.
		Model(&model.Notification{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Updates(map[string]interface{}{
			"is_read": true,
			"read_at": gorm.Expr("CURRENT_TIMESTAMP"),
		}).
		Error
}

func (r *NotificationRepository) WithTx(
	tx *gorm.DB,
) *NotificationRepository {
	return NewNotificationRepository(tx)
}
