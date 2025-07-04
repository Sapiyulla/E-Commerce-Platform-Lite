package usecase

import (
	"errors"
	"order-service/internal/domain"

	"github.com/google/uuid"
)

type OrderUseCase struct {
	repo domain.OrderRepository
}

var (
	ErrInvalidID = errors.New("item_id is invalid")
)

func NewOrderUC(repo domain.OrderRepository) *OrderUseCase {
	return &OrderUseCase{repo: repo}
}

func (uc *OrderUseCase) AddNewOrder(item_id string) (string, error) {
	if uuid.Validate(item_id) != nil {
		return "", ErrInvalidID
	}
	order := domain.NewOrder(item_id)
	if err := uc.repo.Add(order); err != nil {
		return "", err
	}
	return order.ID, nil

}

func (uc *OrderUseCase) PayOrder(order_id string) error {
	if uuid.Validate(order_id) != nil {
		return ErrInvalidID
	}
	return uc.repo.Pay(order_id)
}
