package domain

import "github.com/google/uuid"

type Order struct {
	ID      string
	Item_id string
	Payed   bool
}

func NewOrder(item_id string) *Order {
	return &Order{
		ID:      uuid.New().String(),
		Item_id: item_id,
	}
}

type OrderRepository interface {
	Add(*Order) error
	Pay(id string) error
}
