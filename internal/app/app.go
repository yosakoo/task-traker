package app

import (
    "context"
    "fmt"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/yosakoo/task-traker/internal/config"
    "github.com/yosakoo/task-traker/internal/delivery/http"
    "github.com/yosakoo/task-traker/internal/repository"
    "github.com/yosakoo/task-traker/internal/service"
    "github.com/yosakoo/task-traker/pkg/auth"
    "github.com/yosakoo/task-traker/pkg/hash"
    "github.com/yosakoo/task-traker/pkg/logger"
    "github.com/yosakoo/task-traker/pkg/postgres"
    "github.com/yosakoo/task-traker/pkg/httpserver"
    "github.com/yosakoo/task-traker/pkg/rabbitmq"
)

func Run(cfg *config.Config) {
    l := logger.New(cfg.Log.Level)
    l.Info("start server")

    pg, err := postgres.New(cfg.PG.URL, l)
    if err != nil {
        l.Fatal(fmt.Errorf("app - Run - postgres.New: %w", err))
    }
    defer pg.Close()

    hashSecret := os.Getenv("HASH_SECRET")
    tokenManagerSecret := os.Getenv("TOKEN_MANAGER_SECRET")

    hasher := hash.NewSHA1Hasher(hashSecret)
    tokenManager, err := auth.NewManager(tokenManagerSecret)
    if err != nil {
        l.Error(err)
        return
    }
    rmqConn, err := rabbitmq.New(rabbitmq.Config{
		URL:      cfg.RabbitMQ.URL,
		WaitTime: 5 * time.Second,
		Attempts: 10,
        Exchange: cfg.RabbitMQ.Exchange,
        ExchangeType: cfg.RabbitMQ.ExchangeType,
		Queue: cfg.RabbitMQ.Queue,
	})
	if err != nil {
		l.Fatal(fmt.Errorf("failed to create RabbitMQ connection: %w", err))
	}
    defer rmqConn.Close()

    l.Info("RabbitMQ connected")

    repos := repo.NewRepositories(pg)
    services := service.NewServices(service.Deps{
        Repos:           repos,
        Log:             l,
        Hasher:          hasher,
        TokenManager:    tokenManager,
        QueueConn: rmqConn,
        AccessTokenTTL:  time.Minute * 1,
        RefreshTokenTTL: time.Hour * 24 * 7,
    })

    handlers := http.NewHandler(services, tokenManager)
    srv := server.NewServer(cfg, handlers.Init(l))
    

    go func() {
        if err := srv.Run(); err != nil {
            l.Error(err)
        }
    }()

    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
    <-quit

    const timeout = 5 * time.Second
    ctx, cancel := context.WithTimeout(context.Background(), timeout)
    defer cancel()

    if err := srv.Stop(ctx); err != nil {
        l.Error(fmt.Errorf("failed to stop server: %v", err))
    }
}
