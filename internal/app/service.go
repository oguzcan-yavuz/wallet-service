package app

import (
	"github.com/oguzcan-yavuz/wallet-service/internal/domain/wallet"
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

func (s service) Get(id string) (*wallet.Wallet, error) {
	walletDTO, err := s.repo.Get(id)

	if err != nil {
		return nil, err
	}

	w := wallet.NewWalletFromDTO(walletDTO)

	return w, nil
}
