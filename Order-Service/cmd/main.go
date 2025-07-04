package main

import (
	"context"
	"log"
	"net"
	"order-service/grpc/gen"
	"order-service/internal/application/usecase"
	"order-service/internal/infrastructure/repository"
	"order-service/internal/interfaces/grpcserver"
	"sync"

	"time"

	"github.com/sirupsen/logrus"
	yugapool "github.com/yugabyte/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func init() {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "02.01.2006 15:04",
		ForceColors:     true,
	})
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	wg := sync.WaitGroup{}

	// Initialize the YugabyteDB connection pool
	pool, err := yugapool.New(ctx, "postgres://yugabyte:yugabyte@localhost:5433/yugabyte")
	if err != nil {
		log.Fatalf("Failed to connect to YugabyteDB: %v", err)
	}
	defer pool.Close()

	pingCtx, cancelPing := context.WithTimeout(context.Background(), 7*time.Second)
	defer cancelPing()

	dbReadyCtx, dbReady := context.WithCancel(context.Background())
	defer dbReady()

	wg.Add(1)
	go func() {
		defer wg.Done()
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if err := pool.Ping(context.Background()); err != nil {
					logrus.Errorf("database ping: %s", err.Error())
					cancelPing() // Отменяем pingCtx при ошибке
					return
				}
				dbReady() // Сигнализируем, что БД готова
				return
			case <-pingCtx.Done():
				return
			}
		}
	}()
	logrus.Info("database connected")

	wg.Add(1)
	go func() {
		defer wg.Done()
		<-dbReadyCtx.Done()
		repository.InitMigrate(pool)
	}()

	grpcConn, err := grpc.NewClient("0.0.0.0:5001", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithDefaultCallOptions(grpc.WaitForReady(true)))
	if err != nil {
		logrus.Errorf("grpc: connection error: %s", err.Error())
		return
	}
	defer grpcConn.Close()
	logrus.Info("connected to grpc-protocol with Catalog-Service")

	//создадим прослушиватель порта для grpc(server)
	listener, err := net.Listen("tcp", "0.0.0.0:5002")
	if err != nil {
		logrus.Fatalf("grpc server initialized error: %s", err.Error())
	}
	defer listener.Close()

	repo := repository.NewYugabyteDBRepository(2*time.Second, pool)
	uc := usecase.NewOrderUC(repo)
	grpcClient := gen.NewCatalogServiceClient(grpcConn)
	grpcsrv := grpcserver.NewGrpcServer(grpcClient, uc)

	//инициализируем grpc(server)
	grpcServer := grpc.NewServer()

	//зарегистрируем его на listener
	gen.RegisterOrderServiceServer(grpcServer, grpcsrv)

	wg.Add(1)
	//запустим в горутине чтобы не мешал основному потоку
	go func() {
		defer wg.Done()
		if err := grpcServer.Serve(listener); err != nil {
			logrus.Errorf("grpc server starting error: %s", err.Error())
		}
		defer grpcServer.GracefulStop()
	}()
	logrus.Info("starting listen grpc-server on ", listener.Addr())

	wg.Wait()
	time.Sleep(5 * time.Second)
}
