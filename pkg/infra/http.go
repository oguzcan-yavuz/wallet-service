package infra

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	walletDomain "github.com/oguzcan-yavuz/wallet-service/internal/domain/wallet"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

type AppService interface {
	Get(id string) (*walletDomain.Wallet, error)
}

type Router struct {
	service AppService
}

func (router *Router) get(r *gin.Engine) {
	r.GET("/wallets/:id", func(c *gin.Context) {
		id := c.Param("id")
		wallet, err := router.service.Get(id)

		if err != nil {
			fmt.Println(err)
			c.String(500, err.Error())
			return
		}

		c.JSON(200, wallet)
	})
}

func InitRouter(service AppService) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	r := gin.Default()

	// register endpoints
	router := Router{service: service}
	router.get(r)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	<-ctx.Done()

	stop()
	log.Println("shutting down gracefully, press Ctrl+C again to force")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Server exiting")
}
