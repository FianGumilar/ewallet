package notification

import (
	"context"

	"fiangumilar.id/e-wallet/domain"
	"fiangumilar.id/e-wallet/dto"
)

type service struct {
	notificationRepository domain.NotificationRepository
}

func NewNotificationService(notificationRepository domain.NotificationRepository) domain.NotificationService {
	return &service{
		notificationRepository: notificationRepository,
	}
}

// FindByUser implements domain.NotificationService.
func (s service) FindByUser(ctx context.Context, user int64) ([]dto.NotificationData, error) {
	notifications, err := s.notificationRepository.FindByUser(ctx, user)
	if err != nil {
		return nil, err
	}

	var result []dto.NotificationData
	for _, v := range notifications {
		result = append(result, dto.NotificationData{
			ID:        v.ID,
			Title:     v.Title,
			Body:      v.Body,
			Status:    v.Status,
			IsRead:    v.IsRead,
			CreatedAt: v.CreatedAt,
		})
	}
	if result == nil {
		result = make([]dto.NotificationData, 0)
	}
	return result, nil
}
