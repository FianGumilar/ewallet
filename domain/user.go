package domain

import (
	"context"

	"fiangumilar.id/e-wallet/dto"
)

type User struct {
	ID       int64  `gorm:"primaryKey"`
	FullName string `gorm:"column:fullname"`
	Phone    string `gorm:"column:phone"`
	Username string `gorm:"column:username"`
	Password string `gorm:"column:password"`
}

type UserRepository interface {
	FindByID(ctx context.Context, id int64) (User, error)
	FindByUsername(ctx context.Context, username string) (User, error)
}

type UserService interface {
	Authenticate(ctx context.Context, req dto.UserReq) (dto.UserRes, error)
	ValidateToken(ctx context.Context, token string) (dto.UserData, error)
}
