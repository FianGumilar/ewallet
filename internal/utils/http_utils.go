package utils

import (
	"errors"

	"fiangumilar.id/e-wallet/domain"
)

func GetHttpStatus(err error) int {
	switch {
	case errors.Is(err, domain.ErrPinInvalid):
		return 400

	default:
		return 500
	}
}
