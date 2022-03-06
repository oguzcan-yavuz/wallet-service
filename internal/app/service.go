package app

import (
	"github.com/oguzcan-yavuz/wallet-service/internal/domain"
	"github.com/oguzcan-yavuz/wallet-service/pkg/infra"
)

type Repository interface {
	Get(id string) (*infra.WalletDTO, error)
	Create(wallet *domain.Wallet) (*infra.WalletDTO, error)
	Update(wallet *domain.Wallet) (*infra.WalletDTO, error)
}

type service struct {
	repo Repository
}

func NewService(repository Repository) *service {
	return &service{repo: repository}
}

func (s service) Get(id string) (*domain.Wallet, error) {
	walletDTO, err := s.repo.Get(id)
	if err != nil {
		return nil, err
	}

	wallet := walletDTO.ToDomain()

	return wallet, nil
}

func (s service) Create() (*domain.Wallet, error) {
	wallet, err := domain.NewWallet()
	if err != nil {
		return nil, err
	}

	_, err = s.repo.Create(wallet)
	if err != nil {
		return nil, err
	}

	return wallet, nil
}

func (s service) Withdraw(id string, amount int64) (*domain.Wallet, error) {
	walletDTO, err := s.repo.Get(id)
	if err != nil {
		return nil, err
	}

	wallet := walletDTO.ToDomain()

	err = wallet.Withdraw(amount)
	if err != nil {
		return nil, err
	}

	_, err = s.repo.Update(wallet)
	if err != nil {
		return nil, err
	}

	return wallet, nil
}

func (s service) Deposit(id string, amount int64) (*domain.Wallet, error) {
	walletDTO, err := s.repo.Get(id)
	if err != nil {
		return nil, err
	}

	wallet := walletDTO.ToDomain()

	err = wallet.Deposit(amount)
	if err != nil {
		return nil, err
	}

	_, err = s.repo.Update(wallet)
	if err != nil {
		return nil, err
	}

	return wallet, nil
}
