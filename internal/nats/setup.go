package nats

import (
	"encoding/json"
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/lilpipidron/order-desplay-service/internal/config"
	"github.com/lilpipidron/order-desplay-service/internal/models"
	"github.com/lilpipidron/order-desplay-service/internal/storage/postgresql/order"
	"github.com/nats-io/nats.go"
)

func acceptMessage(msg *nats.Msg, orderRepo order.Repository) {
	var order models.Order
	err := json.Unmarshal(msg.Data, &order)
	if err != nil {
		log.Printf("error unmarshalling message: %v", err)
		return
	}

	prettyJSON, err := json.MarshalIndent(order, "", "    ")
	if err != nil {
		log.Printf("error marshalling message to JSON: %v", err)
		return
	}

	fmt.Println(string(prettyJSON))

	err = orderRepo.AddOrder(order)
	if err != nil {
		log.Errorf("error adding order: %v", err)
	}
}

func Setup(cfg *config.Config, orderRepo order.Repository) *nats.Conn {
	natsURL := cfg.NatsStreamingConfig.ClientAddress
	if natsURL == "" {
		natsURL = nats.DefaultURL
	}

	nc, err := nats.Connect(natsURL)
	if err != nil {
		log.Fatal("failed to connect to NATS: ", "err", err)
	}

	log.Info("Connected to NATS Streaming", "url", natsURL)

	_, err = nc.Subscribe("wildberries", func(msg *nats.Msg) {
		acceptMessage(msg, orderRepo)
	})

	if err != nil {
		log.Fatal("failed to subscribe: ", "err", err)
	}

	log.Info("Subscribed to nats-streaming", "subject", "wildberries")
	return nc
}
