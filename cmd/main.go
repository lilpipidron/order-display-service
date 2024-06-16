package main

import (
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/lilpipidron/order-desplay-service/internal/config"
	ns "github.com/lilpipidron/order-desplay-service/internal/nats"
	"github.com/lilpipidron/order-desplay-service/internal/storage/postgresql"
	"github.com/lilpipidron/order-desplay-service/internal/storage/postgresql/order"
	"github.com/nats-io/nats.go"
	"net/http"
)

func main() {
	cfg := config.MustLoad()

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		cfg.Postgres.Host, cfg.Postgres.Port, cfg.Postgres.User, cfg.Postgres.Password, cfg.Postgres.DBName)

	storage, err := postgresql.NewPostgresDB(psqlInfo, cfg.DBName)
	if err != nil {
		log.Fatal("failed to init storage: ", "err", err)
	}

	orderRepo := order.NewRepository(storage.DB)

	nc := ns.Setup(cfg, orderRepo)

	defer func(nc *nats.Conn) {
		err := nc.Drain()
		if err != nil {
			log.Fatal("failed to drain connection: ", "err", err)
		}
	}(nc)

	srv := &http.Server{
		Addr: cfg.HTTPServer.Address,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start http server: ", "err", err)
	}
}
