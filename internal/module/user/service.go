package user

import (
	"context"
	"encoding/json"
	"log"
	"time"

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
		UserId: user.ID,
		Token:  token,
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

// Register implements domain.UserService.
func (s service) Register(ctx context.Context, req dto.UserRegisterReq) (dto.UserRegisterRes, error) {
	// Check req username
	result, err := s.userRepository.FindByUsername(ctx, req.Username)
	if err != nil {
		return dto.UserRegisterRes{}, err
	}

	// Check if username exists
	if result != (domain.User{}) {
		return dto.UserRegisterRes{}, domain.ErrUsernameExists
	}

	user := domain.User{
		FullName: req.FullName,
		Phone:    req.Phone,
		Email:    req.Email,
		Username: req.Username,
		Password: req.Password,
	}

	err = s.userRepository.Insert(ctx, &user)
	if err != nil {
		return dto.UserRegisterRes{}, domain.ErrUsernameExists
	}

	otpCode := utils.GenerateRandomInt(6)
	referenceID := utils.GenerateRandomString(16)

	log.Printf("Your OTP code: %s", otpCode)
	_ = s.cacheRepository.Set("otp:"+referenceID, []byte(otpCode))
	_ = s.cacheRepository.Set("user-ref:"+referenceID, []byte(user.Username))

	return dto.UserRegisterRes{
		ReferenceID: referenceID,
	}, nil

}

// ValidateOtp implements domain.UserService.
func (s service) ValidateOtp(ctx context.Context, req dto.ValidateOtpReq) error {
	val, err := s.cacheRepository.Get("otp:" + req.ReferenceID)
	if err != nil {
		log.Printf("Error retrieving OTP from cache: %s", err)
		return domain.ErrOtpInvalid
	}

	otp := string(val)
	log.Printf("Retrieved OTP: %s", otp)
	if otp != req.OTP {
		return domain.ErrOtpInvalid
	}

	val, err = s.cacheRepository.Get("user-ref:" + req.ReferenceID)
	if err != nil {
		log.Printf("Error retrieving user reference from cache: %s", err)
		return domain.ErrOtpInvalid
	}

	user, err := s.userRepository.FindByUsername(ctx, string(val))
	if err != nil {
		log.Printf("Error retrieving user by username: %s", err)
		return domain.ErrOtpInvalid
	}

	user.EmailVerifiedAt = time.Now()
	_ = s.userRepository.Update(ctx, &user)
	return nil
}
