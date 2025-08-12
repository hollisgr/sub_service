package main

import (
	"context"
	"main/internal/config"
	"main/internal/db"
	"main/internal/handler"
	"main/internal/subscription"
	"main/pkg/logger"
	"main/pkg/postgres"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
)

func main() {
	logger := logger.NewLogger()
	config := config.GetConfig()

	pgxPool, err := postgres.NewPool(context.Background(), 5, *config)

	if err != nil {
		logger.Fatalln(err)
	}
	logger.Infoln("creating new pgx pool OK")

	err = pgxPool.Ping(context.Background())

	if err != nil {
		logger.Fatalln(err)
	}
	logger.Infoln("database ping OK")

	router := gin.Default()
	storage := db.NewDataBase(pgxPool, logger)

	service := subscription.NewService(storage, logger)

	handler := handler.NewHandler(router, service, logger)
	handler.Register()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		logger.Infoln("Interrupt signal received. Exiting...")
		pgxPool.Close()
		os.Exit(0)
	}()
	router.Run(config.Listen.Addr)
}
