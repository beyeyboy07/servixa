package service

import (
	"errors"
	"fmt"
	"strings"

	"servixa/internal/model"
	"servixa/internal/repository"

	"gorm.io/gorm"
)

type NotificationService struct {
	db                     *gorm.DB
	notificationRepository *repository.NotificationRepository
}

func NewNotificationService(
	db *gorm.DB,
	notificationRepository *repository.NotificationRepository,
) *NotificationService {
	return &NotificationService{
		db:                     db,
		notificationRepository: notificationRepository,
	}
}

func (s *NotificationService) Create(
	userID uint,
	notificationType model.NotificationType,
	title string,
	message string,
	reference string,
) (*model.Notification, error) {
	title = strings.TrimSpace(title)
	message = strings.TrimSpace(message)
	reference = strings.TrimSpace(reference)

	if userID == 0 {
		return nil, errors.New("user ID is required")
	}

	if title == "" {
		return nil, errors.New("notification title is required")
	}

	if message == "" {
		return nil, errors.New("notification message is required")
	}

	if len(title) > 150 {
		return nil, errors.New("notification title must not exceed 150 characters")
	}

	if len(notificationType) == 0 || len(notificationType) > 50 {
		return nil, errors.New("invalid notification type")
	}

	if len(reference) > 100 {
		return nil, errors.New("notification reference must not exceed 100 characters")
	}

	notification := &model.Notification{
		UserID:    userID,
		Type:      notificationType,
		Title:     title,
		Message:   message,
		Reference: reference,
		IsRead:    false,
	}

	err := s.db.Transaction(func(tx *gorm.DB) error {
		var user model.User

		if err := tx.
			Select("id", "is_active").
			First(&user, userID).
			Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("user not found")
			}

			return err
		}

		if !user.IsActive {
			return errors.New("user is inactive")
		}

		txRepository := s.notificationRepository.WithTx(tx)

		return txRepository.Create(notification)
	})

	if err != nil {
		return nil, err
	}

	return notification, nil
}

func (s *NotificationService) GetByUserID(
	userID uint,
) ([]model.Notification, error) {
	if userID == 0 {
		return nil, errors.New("user ID is required")
	}

	return s.notificationRepository.FindByUserID(userID)
}

func (s *NotificationService) GetUnreadByUserID(
	userID uint,
) ([]model.Notification, error) {
	if userID == 0 {
		return nil, errors.New("user ID is required")
	}

	return s.notificationRepository.FindUnreadByUserID(userID)
}

func (s *NotificationService) CountUnreadByUserID(
	userID uint,
) (int64, error) {
	if userID == 0 {
		return 0, errors.New("user ID is required")
	}

	return s.notificationRepository.CountUnreadByUserID(userID)
}

func (s *NotificationService) MarkAsRead(
	notificationID uint,
	userID uint,
) error {
	if notificationID == 0 || userID == 0 {
		return errors.New("notification ID and user ID are required")
	}

	notification, err := s.notificationRepository.FindByID(notificationID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("notification not found")
		}

		return err
	}

	if notification.UserID != userID {
		return errors.New("notification not found")
	}

	if notification.IsRead {
		return nil
	}

	return s.notificationRepository.MarkAsRead(
		notificationID,
		userID,
	)
}

func (s *NotificationService) MarkAllAsRead(
	userID uint,
) error {
	if userID == 0 {
		return errors.New("user ID is required")
	}

	return s.notificationRepository.MarkAllAsRead(userID)
}

func (s *NotificationService) CreateWorkOrderNotification(
	userID uint,
	notificationType model.NotificationType,
	workOrderID uint,
	title string,
	message string,
) (*model.Notification, error) {
	if workOrderID == 0 {
		return nil, errors.New("work order ID is required")
	}

	return s.Create(
		userID,
		notificationType,
		title,
		message,
		fmt.Sprintf("WORK_ORDER:%d", workOrderID),
	)
}
