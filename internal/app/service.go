package app

import (
	walletDomain "github.com/oguzcan-yavuz/wallet-service/internal/domain/wallet"
	"github.com/oguzcan-yavuz/wallet-service/pkg/infra"
)

type Repository interface {
	Get(id string) (*infra.WalletDTO, error)
	Create(wallet *walletDomain.Wallet) (*infra.WalletDTO, error)
	Update(wallet *walletDomain.Wallet) (*infra.WalletDTO, error)
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

func (s service) Create() (*walletDomain.Wallet, error) {
	wallet, err := walletDomain.NewWallet()
	if err != nil {
		return nil, err
	}

	_, err = s.repo.Create(wallet)
	if err != nil {
		return nil, err
	}

	return wallet, nil
}

func (s service) Withdraw(id string, amount int64) (*walletDomain.Wallet, error) {
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

func (s service) Deposit(id string, amount int64) (*walletDomain.Wallet, error) {
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
