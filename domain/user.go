package domain

import (
	"context"
	"database/sql"
	"time"

	"fiangumilar.id/e-wallet/dto"
)

type User struct {
	ID                int64        `gorm:"primaryKey"`
	FullName          string       `gorm:"column:fullname"`
	Phone             string       `gorm:"column:phone"`
	Email             string       `gorm:"column:email"`
	Username          string       `gorm:"column:username"`
	Password          string       `gorm:"column:password"`
	EmailVerifiedAtDB sql.NullTime `gorm:"column:email_verified_at"`
	EmailVerifiedAt   time.Time    `gorm:"-"`
}

type UserRepository interface {
	FindByID(ctx context.Context, id int64) (User, error)
	FindByUsername(ctx context.Context, username string) (User, error)
	Insert(ctx context.Context, user *User) error
	Update(ctx context.Context, user *User) error
}

type UserService interface {
	Authenticate(ctx context.Context, req dto.UserReq) (dto.UserRes, error)
	ValidateToken(ctx context.Context, token string) (dto.UserData, error)
	Register(ctx context.Context, req dto.UserRegisterReq) (dto.UserRegisterRes, error)
	ValidateOtp(ctx context.Context, req dto.ValidateOtpReq) error
}
