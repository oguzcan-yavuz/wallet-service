package infra

import (
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/oguzcan-yavuz/wallet-service/internal/domain"
	"os"
)

type WalletDTO struct {
	Id      string `json:"id"`
	Balance int64  `json:"balance"`
}

func (w *WalletDTO) ToDomain() *domain.Wallet {
	return &domain.Wallet{
		ID:      w.Id,
		Balance: w.Balance,
	}
}

type PostgresRepository struct {
	db *pg.DB
}

func NewPostgresRepository() *PostgresRepository {
	db := pg.Connect(&pg.Options{
		User:     "postgres",
		Password: "postgres",
		Database: "postgres",
		Addr:     os.Getenv("POSTGRES_ADDR"),
	})
	err := db.Model(&WalletDTO{}).CreateTable(&orm.CreateTableOptions{
		Temp: true,
	})
	if err != nil {
		panic(err)
	}

	return &PostgresRepository{db: db}
}

func (repo *PostgresRepository) Get(id string) (*WalletDTO, error) {
	walletDTO := &WalletDTO{Id: id}
	err := repo.db.Model(walletDTO).WherePK().Select()

	return walletDTO, err
}

func (repo *PostgresRepository) Create(wallet *domain.Wallet) (*WalletDTO, error) {
	walletDTO := &WalletDTO{Id: wallet.ID, Balance: wallet.Balance}
	_, err := repo.db.Model(walletDTO).Insert()

	return walletDTO, err
}

func (repo *PostgresRepository) Update(wallet *domain.Wallet) (*WalletDTO, error) {
	walletDTO := &WalletDTO{Id: wallet.ID, Balance: wallet.Balance}
	_, err := repo.db.Model(walletDTO).WherePK().Update()

	return walletDTO, err
}
