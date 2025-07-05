package domain

import (
	"context"
	dto "payment-service/internal/application/DTO"
	"time"

	"github.com/google/uuid"
)

type Item interface{}

type Payment struct {
	ID       string `bson:"id"`
	Order_id string `bson:"order_id"`
	Payed_at string `bson:"payed_at"`
}

func NewPayment(order_id string) Payment {
	return Payment{ID: uuid.NewString(), Order_id: order_id, Payed_at: time.Now().Format(time.RFC3339)}
}

type PaymentRepository interface {
	Pay(ctx context.Context, payment Payment) dto.PaymentStatus
}
