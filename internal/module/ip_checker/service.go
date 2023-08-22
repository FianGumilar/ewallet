package ipchecker

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"fiangumilar.id/e-wallet/domain"
	"fiangumilar.id/e-wallet/dto"
)

type service struct {
}

func NewIpCheckerService() domain.IpChecker {
	return &service{}
}

// Query implements domain.IpChecker.
func (s service) Query(ctx context.Context, ip string) (checker dto.IpChecker, err error) {
	// Create rest client when req service
	url := fmt.Sprintf("http://ip-api.com/json/%s?fields=status,lat,lon,timezone,query", ip)

	resp, err := http.Get(url)
	if err != nil {
		return dto.IpChecker{}, err
	}
	defer resp.Body.Close()

	// Get Body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return dto.IpChecker{}, err
	}

	err = json.Unmarshal(body, &checker)
	return
}
