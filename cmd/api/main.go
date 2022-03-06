package main

import (
	"github.com/joho/godotenv"
	"github.com/oguzcan-yavuz/wallet-service/internal/app"
	"github.com/oguzcan-yavuz/wallet-service/pkg/infra"
)

// @title Wallet-Service
// @version 1.0
// @description Wallet Service REST API.

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	postgresRepo := infra.NewPostgresRepository()
	redisRepo := infra.NewRedisRepository()
	service := app.NewService(postgresRepo)
	infra.InitRouter(service, redisRepo)
}
