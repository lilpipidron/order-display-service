package nats

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/lilpipidron/order-desplay-service/internal/config"
	"github.com/lilpipidron/order-desplay-service/internal/models"
	"github.com/lilpipidron/order-desplay-service/internal/storage/postgresql/order"
	"github.com/lilpipidron/order-desplay-service/internal/storage/redis"
	"github.com/nats-io/nats.go"
)

func acceptMessage(msg *nats.Msg, orderRepo order.Repository, redisRepo redis.Repository) {
	dec := json.NewDecoder(bytes.NewReader(msg.Data))
	dec.DisallowUnknownFields()
	var order models.Order
	err := dec.Decode(&order)
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

	go func() {
		err := orderRepo.AddOrder(order)
		if err != nil {
			log.Errorf("error adding order: %v", err)
		}
	}()

	go func() {
		err := redisRepo.AddOrder(&order)
		if err != nil {
			log.Errorf("error adding order to Redis: %v", err)
		}
	}()

	err = msg.Ack()
	if err != nil {
		log.Errorf("error acknowledging message: %v", err)
	}
}

func Setup(cfg *config.Config, orderRepo order.Repository, redisRepo redis.Repository) *nats.Conn {
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
		acceptMessage(msg, orderRepo, redisRepo)
	})

	if err != nil {
		log.Fatal("failed to subscribe: ", "err", err)
	}

	log.Info("Subscribed to nats-streaming", "subject", "wildberries")
	return nc
}
