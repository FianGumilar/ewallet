package factor

import (
	"context"

	"fiangumilar.id/e-wallet/domain"
	"fiangumilar.id/e-wallet/dto"
	"golang.org/x/crypto/bcrypt"
)

type service struct {
	factorRepository domain.FactorRepository
}

func NewFactorService(factorRepositry domain.FactorRepository) domain.FactorService {
	return &service{
		factorRepository: factorRepositry,
	}
}

// ValidatePIN implements domain.FactorService.
func (s service) ValidatePIN(ctx context.Context, req dto.ValidatePinReq) error {
	factor, err := s.factorRepository.FindByUser(ctx, req.UserId)
	if err != nil {
		return err
	}

	if factor == (domain.Factor{}) {
		return domain.ErrPinInvalid
	}

	err = bcrypt.CompareHashAndPassword([]byte(factor.PIN), []byte(req.PIN))
	if err != nil {
		return domain.ErrPinInvalid
	}
	return nil
}
