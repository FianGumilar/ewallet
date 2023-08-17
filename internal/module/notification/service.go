package notification

import (
	"bytes"
	"context"
	"html/template"
	"time"

	"fiangumilar.id/e-wallet/domain"
	"fiangumilar.id/e-wallet/dto"
)

type service struct {
	notificationRepository domain.NotificationRepository
	templateRepositry      domain.TemplateRepository
	hub                    *dto.Hub
}

func NewNotificationService(
	notificationRepository domain.NotificationRepository,
	templateRepository domain.TemplateRepository,
	hub *dto.Hub) domain.NotificationService {
	return &service{
		notificationRepository: notificationRepository,
		templateRepositry:      templateRepository,
		hub:                    hub,
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

// Insert implements domain.NotificationService.
func (s service) Insert(ctx context.Context, userID int64, code string, data map[string]string) error {
	tmpl, err := s.templateRepositry.FindByCode(ctx, code)
	if err != nil {
		return err
	}

	if tmpl == (domain.Template{}) {
		return domain.ErrCodeNotFound
	}

	// buffer for save message notification
	body := new(bytes.Buffer)
	t := template.Must(template.New("notification").Parse(tmpl.Body))
	err = t.Execute(body, data)
	if err != nil {
		return err
	}

	notification := domain.Notification{
		UserID:    userID,
		Title:     tmpl.Title,
		Body:      body.String(),
		Status:    1,
		IsRead:    0,
		CreatedAt: time.Now(),
	}
	err = s.notificationRepository.Insert(ctx, &notification)
	if err != nil {
		return err
	}
	if channel, ok := s.hub.NotificationChannel[userID]; ok {
		channel <- dto.NotificationData{
			ID:        notification.ID,
			Title:     notification.Title,
			Body:      notification.Body,
			Status:    notification.Status,
			IsRead:    notification.IsRead,
			CreatedAt: notification.CreatedAt,
		}
	}
	return nil
}
