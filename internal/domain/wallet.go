package domain

import (
	"errors"
	"github.com/google/uuid"
)

type Wallet struct {
	ID      string
	Balance int64
}

func NewWallet() (*Wallet, error) {
	id, err := uuid.NewUUID()

	if err != nil {
		return nil, err
	}

	wallet := &Wallet{
		ID:      id.String(),
		Balance: 0,
	}

	return wallet, nil
}

func (w *Wallet) Withdraw(amount int64) error {
	if amount <= 0 {
		return errors.New("amount can't be negative or zero")
	}

	if w.Balance < amount {
		return errors.New("the balance is insufficient")
	}

	w.Balance -= amount

	return nil
}

func (w *Wallet) Deposit(amount int64) error {
	if amount <= 0 {
		return errors.New("amount can't be negative or zero")
	}

	w.Balance += amount

	return nil
}
