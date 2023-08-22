package domain

import (
	"context"

	"fiangumilar.id/e-wallet/dto"
)

type IpChecker interface {
	Query(ctx context.Context, ip string) (dto.IpChecker, error)
}
