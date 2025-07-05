package main

import (
	"context"
	"payment-service/grpc/gen"
	"payment-service/internal/application/usecases"
	"payment-service/internal/infrastructure/repository"
	"payment-service/internal/interfaces/httpserver"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func init() {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "02.01.2006 15:04",
		ForceColors:     true,
	})
	gin.SetMode(gin.ReleaseMode)
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOpts := options.Client().ApplyURI("mongodb://root:example@localhost:27017")
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		logrus.Fatal(err)
	}
	defer client.Disconnect(ctx)
	logrus.Infof("connected to database: mongo (%s)", clientOpts.GetURI())

	db := client.Database("paymentsdb")

	err = db.CreateCollection(ctx, "payments")
	if err != nil {
		// если коллекция уже существует — будет ошибка
		if mongo.IsDuplicateKeyError(err) {
			logrus.Warn("Collection already exists")
		} else {
			logrus.Fatal(err)
		}
	}

	grpcClient, err := grpc.NewClient("localhost:5002", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logrus.Error(err)
		return
	}
	defer grpcClient.Close()

	orderServiceClient := gen.NewOrderServiceClient(grpcClient)

	repo := repository.NewMongoRepository(db)
	usecase := usecases.NewPaymentUseCase(orderServiceClient, repo)
	httpsrv := httpserver.NewHttpServer(usecase)

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		r := gin.Default()
		r.POST("/payment/:order_id", httpsrv.PayHandler)
		logrus.Info("starting http-server proccessing...")
		if err := r.Run("0.0.0.0:8003"); err != nil {
			logrus.Error(err)
			return
		}
	}()

	wg.Wait()
}
