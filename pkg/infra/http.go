package infra

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/oguzcan-yavuz/wallet-service/api/docs"
	"github.com/oguzcan-yavuz/wallet-service/internal/domain"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

type AppService interface {
	Get(id string) (*domain.Wallet, error)
	Create() (*domain.Wallet, error)
	Deposit(id string, amount int64) (*domain.Wallet, error)
	Withdraw(id string, amount int64) (*domain.Wallet, error)
}

type IdempotencyMiddleware interface {
	Handler() gin.HandlerFunc
}

type Router struct {
	service     AppService
	idempotency IdempotencyMiddleware
}

//go:generate swag init --parseDependency -g ../../cmd/api/main.go --output=../../api/docs

type WalletResponseDTO struct {
	ID      string `json:"id" example:"5ec7ebf4-9d72-11ec-9802-acde48001122"`
	Balance int64  `json:"balance" example:"100"`
}

// @Description Creates a wallet
// @Tags wallets
// @Accept json
// @Produce json
// @Param Idempotency-Key header string true "UUID for idempotency key"
// @Success 201 {object} WalletResponseDTO "success response"
// @Failure 500 {string} string "error message"
// @Router /wallets [post]
func (router *Router) create(r *gin.Engine) {
	r.POST("/wallets", router.idempotency.Handler(), func(c *gin.Context) {
		wallet, err := router.service.Create()

		if err != nil {
			fmt.Println(err)
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		response := &WalletResponseDTO{
			ID:      wallet.ID,
			Balance: wallet.Balance,
		}
		c.JSON(http.StatusCreated, response)
	})
}

// @Description Gets a wallet
// @Tags wallets
// @Accept json
// @Produce json
// @Param id path string true "Wallet ID"
// @Success 200 {object} WalletResponseDTO "success response"
// @Failure 500 {string} string "error message"
// @Router /wallets/{id} [get]
func (router *Router) get(r *gin.Engine) {
	r.GET("/wallets/:id", func(c *gin.Context) {
		id := c.Param("id")
		wallet, err := router.service.Get(id)
		if err != nil {
			fmt.Println(err)
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		response := &WalletResponseDTO{
			ID:      wallet.ID,
			Balance: wallet.Balance,
		}

		c.JSON(http.StatusOK, response)
	})
}

type UpdateWalletRequestDTO struct {
	Amount int64 `json:"amount" example:"100"`
}

// @Description Deposits into a wallet
// @Tags wallets
// @Accept json
// @Produce json
// @Param Idempotency-Key header string true "UUID for idempotency key"
// @Param id path string true "Wallet ID"
// @Param DepositWalletRequest body UpdateWalletRequestDTO true "deposit amount"
// @Success 200 {object} WalletResponseDTO "success response"
// @Failure 400 {string} string "error message"
// @Failure 500 {string} string "error message"
// @Router /wallets/{id}/deposit [post]
func (router *Router) deposit(r *gin.Engine) {
	r.POST("/wallets/:id/deposit", router.idempotency.Handler(), func(c *gin.Context) {
		id := c.Param("id")
		var body UpdateWalletRequestDTO
		if err := c.ShouldBindJSON(&body); err != nil {
			c.String(http.StatusBadRequest, "invalid amount")
			return
		}

		wallet, err := router.service.Deposit(id, body.Amount)
		if err != nil {
			fmt.Println(err)
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		response := WalletResponseDTO{
			ID:      wallet.ID,
			Balance: wallet.Balance,
		}

		c.JSON(http.StatusOK, response)
	})
}

// @Description Withdraws from a wallet
// @Tags wallets
// @Accept json
// @Produce json
// @Param Idempotency-Key header string true "UUID for idempotency key"
// @Param id path string true "Wallet ID"
// @Param WithdrawWalletRequest body UpdateWalletRequestDTO true "withdraw amount"
// @Success 200 {object} WalletResponseDTO "success response"
// @Failure 400 {string} string "error message"
// @Failure 500 {string} string "error message"
// @Router /wallets/{id}/withdraw [post]
func (router *Router) withdraw(r *gin.Engine) {
	r.POST("/wallets/:id/withdraw", router.idempotency.Handler(), func(c *gin.Context) {
		id := c.Param("id")
		var body UpdateWalletRequestDTO
		if err := c.ShouldBindJSON(&body); err != nil {
			c.String(http.StatusBadRequest, "invalid amount")
			return
		}

		wallet, err := router.service.Withdraw(id, body.Amount)
		if err != nil {
			fmt.Println(err)
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		response := WalletResponseDTO{
			ID:      wallet.ID,
			Balance: wallet.Balance,
		}

		c.JSON(http.StatusOK, response)
	})
}

func (router *Router) docs(r *gin.Engine) {
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	docs.SwaggerInfo.Host = "localhost:8080"
}

func InitRouter(service AppService, repository Repository) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	r := gin.Default()
	idempotency := NewIdempotencyMiddleware(repository)

	// register endpoints
	router := Router{service: service, idempotency: idempotency}
	router.get(r)
	router.create(r)
	router.deposit(r)
	router.withdraw(r)
	router.docs(r)

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
