package main

import (
	"github.com/charmbracelet/log"
	"github.com/lilpipidron/order-desplay-service/internal/config"
)

func main() {
	cnf := config.MustLoad()
	log.Print(cnf)
}
