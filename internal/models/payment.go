package models

type Payment struct {
	RequestID    string  `json:"request_id"`
	Currency     string  `json:"currency"`
	Provider     string  `json:"provider"`
	Amount       float64 `json:"amount"`
	PaymentDt    int64   `json:"payment_dt"`
	Bank         string  `json:"bank"`
	DeliveryCost float64 `json:"delivery_cost"`
	GoodsTotal   float64 `json:"goods_total"`
	CustomFee    float64 `json:"custom_fee"`
	Transaction  string  `json:"transaction"`
}
