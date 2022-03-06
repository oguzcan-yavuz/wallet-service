package infra

import (
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	walletDomain "github.com/oguzcan-yavuz/wallet-service/internal/domain/wallet"
)

type WalletDTO struct {
	Id      string `json:"id"`
	Balance int64  `json:"balance"`
}

func (w *WalletDTO) ToDomain() *walletDomain.Wallet {
	return &walletDomain.Wallet{
		Id:      w.Id,
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
		Addr:     "localhost:5432",
	})
	// Normally you would run migrations manually
	err := db.Model(&WalletDTO{}).CreateTable(&orm.CreateTableOptions{
		Temp: true,
	})
	if err != nil {
		panic(err)
	}

	return &PostgresRepository{db: db}
}

func (repo *PostgresRepository) Get(id string) (*WalletDTO, error) {
	wallet := &WalletDTO{Id: id}
	err := repo.db.Model(wallet).WherePK().Select()

	return wallet, err
}
