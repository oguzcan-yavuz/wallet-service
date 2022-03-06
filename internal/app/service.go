package app

import (
	walletDomain "github.com/oguzcan-yavuz/wallet-service/internal/domain/wallet"
	"github.com/oguzcan-yavuz/wallet-service/pkg/infra"
)

type Repository interface {
	Get(id string) (*infra.WalletDTO, error)
}

type service struct {
	repo Repository
}

func NewService(repository Repository) *service {
	return &service{repo: repository}
}

func (s service) Get(id string) (*walletDomain.Wallet, error) {
	walletDTO, err := s.repo.Get(id)

	if err != nil {
		return nil, err
	}

	wallet := walletDTO.ToDomain()

	return wallet, nil
}
