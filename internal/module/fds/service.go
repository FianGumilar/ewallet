package fds

import (
	"context"
	"time"

	"fiangumilar.id/e-wallet/domain"
	"fiangumilar.id/e-wallet/dto"
	"fiangumilar.id/e-wallet/internal/utils"
)

type service struct {
	ipCheckerService   domain.IpChecker
	loginLogRepository domain.LoginLogRepository
}

func NewFdsService(ipCheckerService domain.IpChecker, loginLogRepository domain.LoginLogRepository) domain.FdsService {
	return &service{
		ipCheckerService:   ipCheckerService,
		loginLogRepository: loginLogRepository,
	}
}

// IsAuthorized implements domain.FdsService.
func (s service) IsAuthorized(ctx context.Context, ip string, userId int64) bool {
	locationCheck, err := s.ipCheckerService.Query(ctx, ip)
	if err != nil || locationCheck == (dto.IpChecker{}) {
		return false
	}

	newAccess := domain.LoginLog{
		UserID:       userId,
		IsAuthorized: false,
		IpAddress:    ip,
		Timezone:     locationCheck.Timezone,
		Lat:          locationCheck.Lat,
		Lon:          locationCheck.Lon,
		AccessTime:   time.Now(),
	}

	// Fetch Last Login
	lastLogin, err := s.loginLogRepository.FindLastAuthorized(ctx, userId)
	if err != nil {
		_ = s.loginLogRepository.Save(ctx, &newAccess)
		return false
	}

	// Check if never login berfore
	if lastLogin == (domain.LoginLog{}) {
		newAccess.IsAuthorized = true
		_ = s.loginLogRepository.Save(ctx, &newAccess)
		return true
	}

	// Check distance or displacement in kilometers
	distanceHour := newAccess.AccessTime.Sub(lastLogin.AccessTime)
	distanceChange := utils.GetDistance(lastLogin.Lat, lastLogin.Lon, newAccess.Lat, newAccess.Lon)

	// Check if distance/time(V= S/T) trick over 400
	if distanceChange/distanceHour.Hours() > 400 {
		_ = s.loginLogRepository.Save(ctx, &newAccess)
		return false
	}

	// Give Access
	newAccess.IsAuthorized = true
	_ = s.loginLogRepository.Save(ctx, &newAccess)
	return true

}
