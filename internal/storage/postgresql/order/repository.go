package order

import (
	"database/sql"
	"github.com/charmbracelet/log"
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

func (repo *OrderRepository) AddOrder(order models.Order) error {
	tx, err := repo.db.Begin()
	if err != nil {
		return err
	}
	query := "INSERT INTO orders (order_uid, track_number, entry, locale," +
		"internal_signature, customer_id, delivery_service, shardkey, sm_id," +
		"date_created, oof_shard) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)"

	_, err = tx.Exec(query, order.OrderUID, order.TrackNumber, order.Entry,
		order.Locale, order.InternalSignature, order.CustomerID, order.DeliveryService,
		order.Shardkey, order.SmID, order.DateCreated, order.OofShard)

	if err != nil {
		err := tx.Rollback()
		if err != nil {
			return err
		}
		return err
	}

	for _, item := range order.Items {
		query = "INSERT INTO items (chrt_id, track_number, price, rid, name, sale," +
			"size, total_price, nm_id, brand, status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)"

		_, err = tx.Exec(query, item.ChrtID, item.TrackNumber, item.Price, item.Rid, item.Name,
			item.Sale, item.Size, item.TotalPrice, item.NmID, item.Brand, item.Status)

		if err != nil {
			err := tx.Rollback()
			if err != nil {
				return err
			}
			return err
		}
	}

	query = "INSERT INTO deliveries (order_uid, name, phone, zip, city," +
		"address, region, email) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)"

	_, err = tx.Exec(query, order.OrderUID, order.Delivery.Name, order.Delivery.Phone, order.Delivery.Zip,
		order.Delivery.City, order.Delivery.Address, order.Delivery.Region, order.Delivery.Email)

	if err != nil {
		err = tx.Rollback()
		if err != nil {
			return err
		}
	}

	query = "INSERT INTO payments (order_uid, request_id, currency, provider," +
		"amount, payment_dt, bank, delivery_cost, goods_total, custom_fee, transaction) " +
		"VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)"

	_, err = tx.Exec(query, order.OrderUID, order.Payment.RequestID, order.Payment.Currency,
		order.Payment.Provider, order.Payment.Amount, order.Payment.PaymentDt,
		order.Payment.Bank, order.Payment.DeliveryCost, order.Payment.GoodsTotal,
		order.Payment.CustomFee, order.Payment.Transaction)

	if err != nil {
		err := tx.Rollback()
		if err != nil {
			return err
		}
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (repo *OrderRepository) GetOrders() ([]models.Order, error) {
	var orders []models.Order
	query := "SELECT * from orders"
	row, err := repo.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Error("failed to close rows: ", "error", err)
		}
	}(row)

	for row.Next() {
		var uid string
		o := &models.Order{}
		err = row.Scan(&o.OrderUID, o.TrackNumber, &o.Entry, &o.Locale, &o.InternalSignature,
			&o.CustomerID, &o.DeliveryService, &o.Shardkey, &o.SmID, &o.DateCreated, &o.OofShard)
		if err != nil {
			return nil, err
		}

		query = "SELECT * from deliveries WHERE order_uid = $1"
		deliveryRow, err := repo.db.Query(query, o.OrderUID)
		if err != nil {
			return nil, err
		}
		deliveryRow.Next()
		err = row.Scan(&uid, &o.Delivery.Name, &o.Delivery.Phone, &o.Delivery.Zip,
			&o.Delivery.City, &o.Delivery.Address, &o.Delivery.Region, &o.Delivery.Email)
		if err != nil {
			return nil, err
		}
		err = deliveryRow.Close()
		if err != nil {
			return nil, err
		}

		query = "SELECT * from payments WHERE order_uid = $1"
		paymentRow, err := repo.db.Query(query, o.OrderUID)
		if err != nil {
			return nil, err
		}
		paymentRow.Next()
		err = paymentRow.Scan(&uid, &o.Payment.RequestID, &o.Payment.Currency,
			&o.Payment.Provider, &o.Payment.Amount, &o.Payment.PaymentDt,
			&o.Payment.Bank, &o.Payment.DeliveryCost, &o.Payment.GoodsTotal,
			&o.Payment.CustomFee, &o.Payment.Transaction)
		if err != nil {
			return nil, err
		}
		err = paymentRow.Close()
		if err != nil {
			return nil, err
		}

		query = "SELECT * from items WHERE track_number = $1"
		itemRow, err := repo.db.Query(query, o.TrackNumber)
		if err != nil {
			return nil, err
		}

		var items []models.Item
		for itemRow.Next() {
			item := &models.Item{}
			err := itemRow.Scan(&item.ChrtID, &item.TrackNumber,
				&item.Price, &item.Rid, &item.Name, &item.Sale, &item.Size,
				&item.TotalPrice, &item.NmID, &item.Brand, &item.Status)
			if err != nil {
				return nil, err
			}
			items = append(items, *item)
		}
		err = itemRow.Close()
		if err != nil {
			return nil, err
		}
		o.Items = items
		orders = append(orders, *o)
	}
	return orders, nil
}
