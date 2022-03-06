package domain

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWallet_NewWallet(t *testing.T) {
	// Act
	wallet, err := NewWallet()

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, wallet.Balance, int64(0))
}

func TestWallet_Withdraw_NegativeAmount(t *testing.T) {
	// Arrange
	wallet := &Wallet{
		ID:      "test",
		Balance: 100,
	}
	amount := int64(-50)

	// Act
	err := wallet.Withdraw(amount)

	// Assert
	assert.NotNil(t, err)
	assert.Equal(t, wallet.Balance, int64(100))
}

func TestWallet_Withdraw_ZeroAmount(t *testing.T) {
	// Arrange
	wallet := &Wallet{
		ID:      "test",
		Balance: 100,
	}
	amount := int64(0)

	// Act
	err := wallet.Withdraw(amount)

	// Assert
	assert.NotNil(t, err)
	assert.Equal(t, wallet.Balance, int64(100))
}

func TestWallet_Withdraw_InsufficientBalance(t *testing.T) {
	// Arrange
	wallet := &Wallet{
		ID:      "test",
		Balance: 100,
	}
	amount := int64(500)

	// Act
	err := wallet.Withdraw(amount)

	// Assert
	assert.NotNil(t, err)
	assert.Equal(t, wallet.Balance, int64(100))
}

func TestWallet_Withdraw_ShouldSubtract(t *testing.T) {
	// Arrange
	wallet := &Wallet{
		ID:      "test",
		Balance: 100,
	}
	amount := int64(50)

	// Act
	err := wallet.Withdraw(amount)

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, wallet.Balance, int64(50))
}

func TestWallet_Deposit_NegativeAmount(t *testing.T) {
	// Arrange
	wallet := &Wallet{
		ID:      "test",
		Balance: 100,
	}
	amount := int64(-50)

	// Act
	err := wallet.Deposit(amount)

	// Assert
	assert.NotNil(t, err)
	assert.Equal(t, wallet.Balance, int64(100))
}

func TestWallet_Deposit_ZeroAmount(t *testing.T) {
	// Arrange
	wallet := &Wallet{
		ID:      "test",
		Balance: 100,
	}
	amount := int64(0)

	// Act
	err := wallet.Deposit(amount)

	// Assert
	assert.NotNil(t, err)
	assert.Equal(t, wallet.Balance, int64(100))
}

func TestWallet_Deposit_ShouldAdd(t *testing.T) {
	// Arrange
	wallet := &Wallet{
		ID:      "test",
		Balance: 100,
	}
	amount := int64(50)

	// Act
	err := wallet.Deposit(amount)

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, wallet.Balance, int64(150))
}
