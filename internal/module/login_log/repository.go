package loginlog

import (
	"context"
	"database/sql"

	"fiangumilar.id/e-wallet/domain"
)

type repositoy struct {
	db *sql.DB
}

func NewLoginLogRepository(con *sql.DB) domain.LoginLogRepository {
	return &repositoy{db: con}
}

// FindLastAuthorized implements domain.LoginLogRepository.
func (r repositoy) FindLastAuthorized(ctx context.Context, userId int64) (loginlog domain.LoginLog, err error) {
	query := `
	SELECT * FROM login_log 
	WHERE user_id = ? AND is_authorized = true 
	ORDER BY id DESC LIMIT 1;
	`

	rows := r.db.QueryRowContext(ctx, query, userId)
	err = rows.Scan(
		&loginlog.ID,
		&loginlog.UserID,
		&loginlog.IsAuthorized,
		&loginlog.IpAddress,
		&loginlog.Timezone,
		&loginlog.Lat,
		&loginlog.Lon,
		&loginlog.AccessTime,
	)
	if err != nil {
		return loginlog, err
	}
	return
}

// Save implements domain.LoginLogRepository.
func (r repositoy) Save(ctx context.Context, login *domain.LoginLog) error {
	query := `INSERT INTO login_log (user_id, is_authorized, ip_address, timezone, lat, lon, access_time)
	VALUES (?, ?, ?, ?, ?, ?, ?)
	RETURNING id;
	`

	rows := r.db.QueryRowContext(ctx, query, login.UserID, login.IsAuthorized, login.IpAddress, login.Timezone, login.Lat, login.Lon, login.AccessTime)
	return rows.Scan(&login.ID)
}
