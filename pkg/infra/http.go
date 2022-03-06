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
	Create() (*walletDomain.Wallet, error)
	Deposit(id string, amount int64) (*walletDomain.Wallet, error)
	Withdraw(id string, amount int64) (*walletDomain.Wallet, error)
}

type Router struct {
	service AppService
}

func (router *Router) create(r *gin.Engine) {
	r.POST("/wallets", func(c *gin.Context) {
		wallet, err := router.service.Create()

		if err != nil {
			fmt.Println(err)
			c.String(500, err.Error())
			return
		}

		c.JSON(201, wallet)
	})
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

type BodyWithAmount struct {
	Amount int64 `json:"amount"`
}

func (router *Router) deposit(r *gin.Engine) {
	r.POST("/wallets/:id/deposit", func(c *gin.Context) {
		id := c.Param("id")
		var body BodyWithAmount
		if err := c.ShouldBindJSON(&body); err != nil {
			c.String(400, "invalid amount")
			return
		}

		wallet, err := router.service.Deposit(id, body.Amount)
		if err != nil {
			fmt.Println(err)
			c.String(500, err.Error())
			return
		}

		c.JSON(200, wallet)
	})
}

func (router *Router) withdraw(r *gin.Engine) {
	r.POST("/wallets/:id/withdraw", func(c *gin.Context) {
		id := c.Param("id")
		var body BodyWithAmount
		if err := c.ShouldBindJSON(&body); err != nil {
			c.String(400, "invalid amount")
			return
		}

		wallet, err := router.service.Withdraw(id, body.Amount)
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
	router.create(r)
	router.deposit(r)
	router.withdraw(r)

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
