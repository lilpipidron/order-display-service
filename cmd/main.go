package main

import (
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/lilpipidron/order-desplay-service/internal/config"
	"github.com/lilpipidron/order-desplay-service/internal/storage/postgresql"
	"net/http"
)

func main() {
	cfg := config.MustLoad()

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		cfg.Postgres.Host, cfg.Postgres.Port, cfg.Postgres.User, cfg.Postgres.Password, cfg.Postgres.DBName)

	_, err := postgresql.NewPostgresDB(psqlInfo, cfg.DBName)

	if err != nil {
		log.Fatal("failed to init storage: ", "err", err)
	}

	srv := &http.Server{
		Addr: cfg.HTTPServer.Address,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start http server: ", "err", err)
	}
}
