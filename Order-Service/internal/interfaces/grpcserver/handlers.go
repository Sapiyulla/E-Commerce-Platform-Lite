package grpcserver

import (
	"context"
	"order-service/grpc/gen"
	"order-service/grpc/gen/common"
	"order-service/internal/application/usecase"
	"order-service/internal/infrastructure/repository"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GrpcServer struct {
	gen.UnimplementedOrderServiceServer
	catalogClient gen.CatalogServiceClient
	uc            *usecase.OrderUseCase
}

func NewGrpcServer(catalogClient gen.CatalogServiceClient, usecase *usecase.OrderUseCase) *GrpcServer {
	return &GrpcServer{catalogClient: catalogClient, uc: usecase}
}

func (g *GrpcServer) CreateOrder(ctx context.Context, id *common.ID) (*gen.Order, error) {
	select {
	case <-ctx.Done():
		return nil, status.Error(codes.DeadlineExceeded, "REQUEST TIMEOUT")
	default:
		// Используем контекст запроса вместо глобального контекста
		catalogCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		item, err := g.catalogClient.GetItem(catalogCtx, &common.ID{ItemId: id.GetItemId()})
		respStatus, ok := status.FromError(err)
		if !ok {
			logrus.Warnf("Catalog-Service: %s", err.Error())
			return nil, status.Error(codes.Internal, "INTERNAL SERVER")
		}
		if err != nil {
			if respStatus.Code() == codes.NotFound {
				return nil, status.Error(codes.NotFound, "NOT FOUND")
			} else if respStatus.Code() == codes.Unimplemented {
				return nil, status.Error(codes.Unavailable, "UNAVIABLE")
			}
		}
		if id, err := g.uc.AddNewOrder(id.ItemId); err == nil {
			return &gen.Order{Id: id, ItemId: item.Id}, nil
		} else {
			if err.Error() == repository.ErrOrderAlreadyExists.Error() {
				return nil, status.Error(codes.AlreadyExists, "ALREADY EXISTS")
			}
			logrus.Warnf("internal: %s", err.Error())
			return nil, status.Error(codes.Internal, "INTERNAL SERVER")
		}

	}
}

func (g *GrpcServer) PayOrder(ctx context.Context, id *common.ID) (*gen.Order, error) {
	select {
	case <-ctx.Done():
		return nil, status.Error(codes.DeadlineExceeded, "REQUEST TIMEOUT")
	default:
		if err := g.uc.PayOrder(id.GetOrderId()); err != nil {
			if err.Error() == usecase.ErrInvalidID.Error() {
				return nil, status.Error(codes.InvalidArgument, "BAD REQUEST")
			}
			logrus.Errorf("internal: %s", err.Error())
			return nil, status.Error(codes.Internal, "INTERNAL SERVER")
		}
		return &gen.Order{}, status.Error(codes.OK, "OK")
	}
}
