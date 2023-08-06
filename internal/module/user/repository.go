package user

import (
	"context"
	"database/sql"

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

// Register implements domain.UserRepository.
func (r repository) Insert(ctx context.Context, user *domain.User) error {
	executor := r.db.Debug().WithContext(ctx).Create(user)
	if executor.Error != nil {
		return executor.Error
	}
	return nil
}

// Update implements domain.UserRepository.
func (r repository) Update(ctx context.Context, user *domain.User) error {
	user.EmailVerifiedAtDB = sql.NullTime{
		Time:  user.EmailVerifiedAt,
		Valid: true,
	}

	executor := r.db.Debug().
		WithContext(ctx).
		Model(&user).
		Where("id = ?", user.ID).
		Updates(
			map[string]interface{}{
				"fullname":          user.FullName,
				"phone":             user.Phone,
				"email":             user.Email,
				"username":          user.Username,
				"password":          user.Password,
				"email_verified_at": user.EmailVerifiedAtDB,
			})
	if executor.Error != nil {
		return executor.Error
	}
	return nil
}
