package main

import (
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/lilpipidron/order-desplay-service/internal/config"
	ns "github.com/lilpipidron/order-desplay-service/internal/nats"
	"github.com/lilpipidron/order-desplay-service/internal/storage/postgresql"
	"github.com/lilpipidron/order-desplay-service/internal/storage/postgresql/order"
	rds "github.com/lilpipidron/order-desplay-service/internal/storage/redis"
	"github.com/nats-io/nats.go"
	"github.com/redis/go-redis/v9"
	"net/http"
	"strconv"
)

func main() {
	cfg := config.MustLoad()

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		cfg.Postgres.Host, cfg.Postgres.Port, cfg.Postgres.User, cfg.Postgres.Password, cfg.Postgres.DBName)

	storage, err := postgresql.NewPostgresDB(psqlInfo, cfg.Postgres.DBName)
	if err != nil {
		log.Fatal("failed to init storage: ", "err", err)
	}

	orderRepo := order.NewRepository(storage.DB)

	addr := cfg.Redis.Host + ":" + strconv.Itoa(cfg.Redis.Port)
	opt := redis.Options{
		Addr:     addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
		Protocol: cfg.Redis.Protocol,
	}
	redisRepo, err := rds.NewRedisRepo(&opt)
	if err != nil {
		log.Fatal("failed to get redis repo: ", "err", err)
	}

	nc := ns.Setup(cfg, orderRepo, redisRepo)

	defer func(nc *nats.Conn) {
		err := nc.Drain()
		if err != nil {
			log.Fatal("failed to drain connection: ", "err", err)
		}
	}(nc)

	orders, err := orderRepo.GetOrders()
	if err != nil {
		log.Fatal("failed to get orders: ", "err", err)
	}

	for _, odr := range orders {
		err := redisRepo.AddOrder(&odr)
		if err != nil {
			log.Fatal("failed to add order: ", "err", err)
		}
	}

	srv := &http.Server{
		Addr: cfg.HTTPServer.Address,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start http server: ", "err", err)
	}
}
