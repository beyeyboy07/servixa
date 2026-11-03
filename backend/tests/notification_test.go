package tests

import (
	"testing"

	"servixa/internal/config"
	"servixa/internal/database"
	"servixa/internal/model"
	"servixa/internal/repository"
	"servixa/internal/service"
)

func TestNotificationService(t *testing.T) {
	cfg := config.Load()
	db := database.Connect(cfg)

	if err := db.AutoMigrate(&model.Notification{}); err != nil {
		t.Fatalf("failed to migrate notification: %v", err)
	}

	const (
		userID      uint = 999991
		otherUserID uint = 999992
	)

	// Cleanup test data before and after the test.
	db.Where("user_id IN ?", []uint{userID, otherUserID}).
		Delete(&model.Notification{})

	defer db.Where("user_id IN ?", []uint{userID, otherUserID}).
		Delete(&model.Notification{})

	repo := repository.NewNotificationRepository(db)
	notificationService := service.NewNotificationService(db, repo)

	notifications := []model.Notification{
		{
			UserID:    userID,
			Type:      model.NotificationWorkOrderCreated,
			Title:     "Work order created",
			Message:   "Your work order has been created",
			Reference: "WORK_ORDER:999991",
		},
		{
			UserID:    userID,
			Type:      model.NotificationWorkStarted,
			Title:     "Work started",
			Message:   "Technician has started the work",
			Reference: "WORK_ORDER:999992",
		},
		{
			UserID:    otherUserID,
			Type:      model.NotificationWorkCompleted,
			Title:     "Work completed",
			Message:   "Another user's work is completed",
			Reference: "WORK_ORDER:999993",
		},
	}

	if err := db.Create(&notifications).Error; err != nil {
		t.Fatalf("failed to create test notifications: %v", err)
	}

	firstID := notifications[0].ID
	secondID := notifications[1].ID
	otherID := notifications[2].ID

	t.Run("GetByUserID only returns own notifications", func(t *testing.T) {
		result, err := notificationService.GetByUserID(userID)
		if err != nil {
			t.Fatalf("failed to get notifications: %v", err)
		}

		if len(result) != 2 {
			t.Fatalf("expected 2 notifications, got %d", len(result))
		}

		for _, item := range result {
			if item.UserID != userID {
				t.Fatalf("returned notification belonging to user %d", item.UserID)
			}
		}
	})

	t.Run("CountUnreadByUserID", func(t *testing.T) {
		count, err := notificationService.CountUnreadByUserID(userID)
		if err != nil {
			t.Fatalf("failed to count unread notifications: %v", err)
		}

		if count != 2 {
			t.Fatalf("expected 2 unread notifications, got %d", count)
		}
	})

	t.Run("MarkAsRead", func(t *testing.T) {
		if err := notificationService.MarkAsRead(firstID, userID); err != nil {
			t.Fatalf("failed to mark notification as read: %v", err)
		}

		result, err := repo.FindByID(firstID)
		if err != nil {
			t.Fatalf("failed to retrieve notification: %v", err)
		}

		if !result.IsRead {
			t.Fatal("expected notification to be read")
		}

		if result.ReadAt == nil {
			t.Fatal("expected read_at to be set")
		}

		count, err := notificationService.CountUnreadByUserID(userID)
		if err != nil {
			t.Fatalf("failed to count unread notifications: %v", err)
		}

		if count != 1 {
			t.Fatalf("expected 1 unread notification, got %d", count)
		}
	})

	t.Run("Cannot mark another user's notification as read", func(t *testing.T) {
		err := notificationService.MarkAsRead(otherID, userID)
		if err == nil {
			t.Fatal("expected marking another user's notification to fail")
		}

		result, err := repo.FindByID(otherID)
		if err != nil {
			t.Fatalf("failed to retrieve other user's notification: %v", err)
		}

		if result.IsRead {
			t.Fatal("expected other user's notification to remain unread")
		}
	})

	t.Run("MarkAllAsRead only affects own notifications", func(t *testing.T) {
		if err := notificationService.MarkAllAsRead(userID); err != nil {
			t.Fatalf("failed to mark all notifications as read: %v", err)
		}

		count, err := notificationService.CountUnreadByUserID(userID)
		if err != nil {
			t.Fatalf("failed to count unread notifications: %v", err)
		}

		if count != 0 {
			t.Fatalf("expected 0 unread notifications, got %d", count)
		}

		otherNotification, err := repo.FindByID(otherID)
		if err != nil {
			t.Fatalf("failed to retrieve other user's notification: %v", err)
		}

		if otherNotification.IsRead {
			t.Fatal("expected other user's notification to remain unread")
		}

		remaining, err := repo.FindByID(secondID)
		if err != nil {
			t.Fatalf("failed to retrieve second notification: %v", err)
		}

		if !remaining.IsRead {
			t.Fatal("expected second notification to be read")
		}
	})
}
