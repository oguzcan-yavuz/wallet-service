package wallet

import (
	"github.com/oguzcan-yavuz/wallet-service/pkg/infra"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWallet_NewWallet(t *testing.T) {
	// Act
	wallet, err := NewWallet()

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, wallet.balance, int64(0))
}

func TestWallet_Withdraw_NegativeAmount(t *testing.T) {
	// Arrange
	wallet := &Wallet{
		id:      "test",
		balance: 100,
	}
	amount := int64(-50)

	// Act
	err := wallet.Withdraw(amount)

	// Assert
	assert.NotNil(t, err)
	assert.Equal(t, wallet.balance, int64(100))
}

func TestWallet_Withdraw_ZeroAmount(t *testing.T) {
	// Arrange
	wallet := &Wallet{
		id:      "test",
		balance: 100,
	}
	amount := int64(0)

	// Act
	err := wallet.Withdraw(amount)

	// Assert
	assert.NotNil(t, err)
	assert.Equal(t, wallet.balance, int64(100))
}

func TestWallet_Withdraw_InsufficientBalance(t *testing.T) {
	// Arrange
	wallet := &Wallet{
		id:      "test",
		balance: 100,
	}
	amount := int64(500)

	// Act
	err := wallet.Withdraw(amount)

	// Assert
	assert.NotNil(t, err)
	assert.Equal(t, wallet.balance, int64(100))
}

func TestWallet_Withdraw_ShouldSubtract(t *testing.T) {
	wallet := &Wallet{
		id:      "test",
		balance: 100,
	}
	amount := int64(50)

	err := wallet.Withdraw(amount)

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, wallet.balance, int64(50))
}

func TestWallet_Deposit_NegativeAmount(t *testing.T) {
	// Arrange
	wallet := &Wallet{
		id:      "test",
		balance: 100,
	}
	amount := int64(-50)

	// Act
	err := wallet.Deposit(amount)

	// Assert
	assert.NotNil(t, err)
	assert.Equal(t, wallet.balance, int64(100))
}

func TestWallet_Deposit_ZeroAmount(t *testing.T) {
	// Arrange
	wallet := &Wallet{
		id:      "test",
		balance: 100,
	}
	amount := int64(0)

	// Act
	err := wallet.Deposit(amount)

	// Assert
	assert.NotNil(t, err)
	assert.Equal(t, wallet.balance, int64(100))
}

func TestWallet_Deposit_ShouldAdd(t *testing.T) {
	wallet := &Wallet{
		id:      "test",
		balance: 100,
	}
	amount := int64(50)

	err := wallet.Deposit(amount)

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, wallet.balance, int64(150))
}

func TestWallet_NewWalletFromDTO(t *testing.T) {
	walletDTO := &infra.WalletDTO{
		Id:      "test",
		Balance: 100,
	}

	wallet := NewWalletFromDTO(walletDTO)

	assert.Equal(t, wallet.id, walletDTO.Id)
	assert.Equal(t, wallet.balance, walletDTO.Balance)
}
