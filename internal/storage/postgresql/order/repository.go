package order

import (
	"database/sql"
	"github.com/lilpipidron/order-desplay-service/internal/models"
)

type Repository interface {
	AddOrder(order models.Order) error
	GetOrders() ([]models.Order, error)
}

type OrderRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (db *OrderRepository) AddOrder(order models.Order) error {
	return nil
}

func GetOrders() ([]models.Order, error) {
	return nil, nil
}
