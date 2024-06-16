package redis

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/lilpipidron/order-desplay-service/internal/models"
	"github.com/redis/go-redis/v9"
)

type Storage struct {
	client *redis.Client
}

type Repository interface {
	AddOrder(order *models.Order) error
	GetOrder(uid string) (*models.Order, error)
}

func NewRedisClient(opt *redis.Options) (*Storage, error) {
	client := redis.NewClient(opt)
	if client == nil {
		return nil, errors.New("redis client is nil")
	}
	return &Storage{client: client}, nil
}

func (repo *Storage) AddOrder(order *models.Order) error {
	jsonOrder, err := json.Marshal(order)
	if err != nil {
		return err
	}
	err = repo.client.LPush(context.Background(), order.OrderUID, jsonOrder).Err()
	if err != nil {
		return err
	}
	return nil
}

func (repo *Storage) GetOrder(uid string) (*models.Order, error) {
	return nil, nil
}
