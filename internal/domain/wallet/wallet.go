package wallet

import (
	"errors"
	"github.com/google/uuid"
	"github.com/oguzcan-yavuz/wallet-service/pkg/infra"
)

type Wallet struct {
	id      string
	balance int64
}

func NewWallet() (*Wallet, error) {
	id, err := uuid.NewUUID()

	if err != nil {
		return nil, err
	}

	wallet := &Wallet{
		id:      id.String(),
		balance: 0,
	}

	return wallet, nil
}

func NewWalletFromDTO(walletDTO *infra.WalletDTO) *Wallet {
	wallet := Wallet{
		id:      walletDTO.Id,
		balance: walletDTO.Balance,
	}

	return &wallet
}

func (w *Wallet) Withdraw(amount int64) error {
	if amount <= 0 {
		return errors.New("amount can't be negative or zero")
	}

	if w.balance < amount {
		return errors.New("the balance is insufficient")
	}

	w.balance -= amount

	return nil
}

func (w *Wallet) Deposit(amount int64) error {
	if amount <= 0 {
		return errors.New("amount can't be negative or zero")
	}

	w.balance += amount

	return nil
}
