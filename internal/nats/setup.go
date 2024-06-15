package nats

import (
	"github.com/charmbracelet/log"
	"github.com/lilpipidron/order-desplay-service/internal/config"
	"github.com/nats-io/nats.go"
)

func acceptMessage(msg *nats.Msg) {
	log.Print(string(msg.Data))
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
