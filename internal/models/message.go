package models

type Message struct {
	Delivery `json:"delivery"`
	Items    []Item  `json:"items"`
	Payment  Payment `json:"payment"`
	Order    `json:"order"`
}
