package user

import (
	"context"

	"fiangumilar.id/e-wallet/domain"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

func NewUserRepository(con *gorm.DB) domain.UserRepository {
	return &repository{db: con}
}

// FindByID implements domain.UserRepository.
func (r repository) FindByID(ctx context.Context, id int64) (user domain.User, err error) {
	dataset := r.db.Where("id = ?", id).First(&user)
	if dataset.Error != nil {
		return user, nil
	}
	return
}

// FindByUsername implements domain.UserRepository.
func (r repository) FindByUsername(ctx context.Context, username string) (user domain.User, err error) {
	dataset := r.db.Where("username = ?", username).First(&user)
	if dataset.Error != nil {
		return user, nil
	}
	return
}
