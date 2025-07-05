package usecases

import (
	"context"
	"payment-service/grpc/gen"
	"payment-service/grpc/gen/common"
	dto "payment-service/internal/application/DTO"
	"payment-service/internal/domain"

	"github.com/google/uuid"
)

type PaymentUseCase struct {
	orderServiceClient gen.OrderServiceClient
	repo               domain.PaymentRepository
}

func NewPaymentUseCase(orderServiceClient gen.OrderServiceClient, repo domain.PaymentRepository) *PaymentUseCase {
	return &PaymentUseCase{orderServiceClient: orderServiceClient, repo: repo}
}

func (uc *PaymentUseCase) Pay(ctx context.Context, order_id string) dto.PaymentStatus {
	if err := uuid.Validate(order_id); err != nil {
		return dto.PaymentStatus{Status: "error", Message: "invalid uuid"}
	}
	_, err := uc.orderServiceClient.PayOrder(ctx, &common.ID{OrderId: order_id})
	if err != nil {
		return dto.PaymentStatus{Status: "error", Message: "internal server error"}
	}
	return uc.repo.Pay(ctx, domain.NewPayment(order_id))
}
