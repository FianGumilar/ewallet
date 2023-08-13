package notification

import (
	"context"
	"database/sql"

	"fiangumilar.id/e-wallet/domain"
)

type repository struct {
	db *sql.DB
}

func NewRepository(con *sql.DB) domain.NotificationRepository {
	return &repository{db: con}
}

// FindByUser implements domain.NotificationRepositry.
func (r repository) FindByUser(ctx context.Context, user int64) (notifications []domain.Notification, err error) {
	query := `SELECT * FROM notifications WHERE user_id = ? ORDER_BY created_at DESC LIMIT 15;`

	rows, err := r.db.QueryContext(ctx, query, user)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var notification domain.Notification
		if err := rows.Scan(
			&notification.ID,
			&notification.Title,
			&notification.Body,
			&notification.Status,
			&notification.IsRead,
			&notification.CreatedAt,
		); err != nil {
			return nil, err
		}
		notifications = append(notifications, notification)
	}
	return notifications, nil
}

// Insert implements domain.NotificationRepositry.
func (r repository) Insert(ctx context.Context, notification *domain.Notification) error {
	query := `INSERT INTO notifications (user_id, title, body, status, is_read, created_at) VALUES (?, ?, ?, ?, ?, ?);`

	_, err := r.db.ExecContext(ctx, query, notification.UserID, notification.Title, notification.Body, notification.Status, notification.IsRead, notification.CreatedAt)
	return err
}

// Update implements domain.NotificationRepositry.
func (r repository) Update(ctx context.Context, notification *domain.Notification) error {
	query := `UPDATE notifications SET title = ?, body = ?, status = ?, isRead = ? WHERE id = ?;`

	_, err := r.db.ExecContext(ctx, query, notification.Title, notification.Body, notification.Status, notification.IsRead, notification.ID)
	return err
}
