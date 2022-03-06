package main

import (
	"github.com/oguzcan-yavuz/wallet-service/internal/app"
	"github.com/oguzcan-yavuz/wallet-service/pkg/infra"
)

func main() {
	postgresRepo := infra.NewPostgresRepository()
	redisRepo := infra.NewRedisRepository()
	service := app.NewService(postgresRepo)
	infra.InitRouter(service, redisRepo)
}
