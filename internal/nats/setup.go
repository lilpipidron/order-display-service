package nats

import (
	"encoding/json"
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/lilpipidron/order-desplay-service/internal/config"
	"github.com/lilpipidron/order-desplay-service/internal/models"
	"github.com/nats-io/nats.go"
)

func acceptMessage(msg *nats.Msg) {
	var message models.Order
	err := json.Unmarshal(msg.Data, &message)
	if err != nil {
		log.Printf("error unmarshalling message: %v", err)
		return
	}

	prettyJSON, err := json.MarshalIndent(message, "", "    ")
	if err != nil {
		log.Printf("error marshalling message to JSON: %v", err)
		return
	}

	fmt.Println(string(prettyJSON))
}

func Setup(cfg *config.Config) *nats.Conn {
	natsURL := cfg.NatsStreamingConfig.ClientAddress
	if natsURL == "" {
		natsURL = nats.DefaultURL
	}

	nc, err := nats.Connect(natsURL)
	if err != nil {
		log.Fatal("failed to connect to NATS: ", "err", err)
	}

	log.Info("Connected to NATS Streaming", "url", natsURL)

	_, err = nc.Subscribe("wildberries", acceptMessage)
	if err != nil {
		log.Fatal("failed to subscribe: ", "err", err)
	}

	log.Info("Subscribed to nats-streaming", "subject", "wildberries")
	return nc
}
