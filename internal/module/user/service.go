package user

import (
	"context"
	"encoding/json"

	"fiangumilar.id/e-wallet/domain"
	"fiangumilar.id/e-wallet/dto"
	"fiangumilar.id/e-wallet/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

type service struct {
	userRepository  domain.UserRepository
	cacheRepository domain.CacheRepository
}

func NewUserService(userRepository domain.UserRepository, cacheRepository domain.CacheRepository) domain.UserService {
	return &service{
		userRepository:  userRepository,
		cacheRepository: cacheRepository,
	}
}

// Authenticate implements domain.UserService.
func (s service) Authenticate(ctx context.Context, req dto.UserReq) (dto.UserRes, error) {
	// Check user exist
	user, err := s.userRepository.FindByUsername(ctx, req.Username)
	if err != nil {
		return dto.UserRes{}, domain.ErrUserNotFound
	}

	if user == (domain.User{}) {
		return dto.UserRes{}, domain.ErrUserNotFound
	}

	//Compare req & user password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return dto.UserRes{}, domain.ErrAuthFailed
	}

	token := utils.GenerateRandomString(16)

	//Set token as cache
	userJson, _ := json.Marshal(user)
	_ = s.cacheRepository.Set("user:"+token, userJson)

	return dto.UserRes{
		Token: token,
	}, nil

}

// ValidateToken implements domain.UserService.
func (s service) ValidateToken(ctx context.Context, token string) (dto.UserData, error) {
	data, err := s.cacheRepository.Get("user:" + token)
	if err != nil {
		return dto.UserData{}, domain.ErrAuthFailed
	}

	var user domain.User
	_ = json.Unmarshal(data, &user)

	return dto.UserData{
		ID:       user.ID,
		FullName: user.FullName,
		Phone:    user.Phone,
		Username: user.Username,
	}, nil
}
