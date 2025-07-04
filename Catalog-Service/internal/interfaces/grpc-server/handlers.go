package grpcserver

import (
	"catalog-service/grpc/gen"
	"catalog-service/grpc/gen/common"
	dto "catalog-service/internal/application/DTO"
	"catalog-service/internal/application/usecase"
	"context"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GrpcServer struct {
	gen.UnimplementedCatalogServiceServer
	usecase *usecase.ItemUseCase
}

func NewGrpcServer(uc *usecase.ItemUseCase) *GrpcServer {
	return &GrpcServer{usecase: uc}
}

func (g *GrpcServer) GetItem(ctx context.Context, id *common.ID) (*gen.Item, error) {
	select {
	case <-ctx.Done():
		deadline, has := ctx.Deadline()
		logrus.Infof("deadline: %d has: %t", deadline.Second(), has)
		return nil, status.Error(codes.DeadlineExceeded, "REQUEST TIMEOUT")
	default:
		item, sts := g.usecase.GetItem(id.GetItemId())
		switch sts.Code {
		case 0:
			return &gen.Item{Id: item.ID, Name: item.Name, Category: item.Category, Price: item.Price}, nil
		case 1:
			if sts.Status == dto.NotFound {
				return nil, status.Error(codes.NotFound, "NOT FOUND")
			}
			return nil, status.Error(codes.InvalidArgument, "INVALID ARGUMENT")
		default:
			return nil, status.Error(codes.Internal, "SERVER_ERROR")
		}
	}
}
