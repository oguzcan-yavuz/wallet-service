package main

import (
	"github.com/oguzcan-yavuz/wallet-service/internal/app"
	"github.com/oguzcan-yavuz/wallet-service/pkg/infra"
)

func main() {
	repo := infra.NewPostgresRepository()
	service := app.NewService(repo)
	infra.InitRouter(service)
}
