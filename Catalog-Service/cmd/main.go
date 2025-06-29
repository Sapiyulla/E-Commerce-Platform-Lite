package main

import (
	"catalog-service/grpc/gen"
	"catalog-service/internal/application/usecase"
	"catalog-service/internal/infrastructure/repository"
	grpcserver "catalog-service/internal/interfaces/grpc-server"
	httpserver "catalog-service/internal/interfaces/http-server"
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func init() {
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:     true,
		TimestampFormat: "02.01.2006 15:04",
		FullTimestamp:   true,
	})
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // гарантированное завершение при выходе

	pool, err := pgxpool.New(ctx, "postgres://postgres:admin@localhost:5432/catalog")
	if err != nil {
		logrus.Fatal("Postgres connection error:", err)
	}
	defer pool.Close()

	// Горутина для проверки соединения
	go func() {
		ticker := time.NewTicker(5 * time.Second) // проверка каждые 5 сек
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if err := pool.Ping(ctx); err != nil {
					logrus.Error("Postgres ping error: ", err)
					cancel() // инициируем завершение
					return
				}
			case <-ctx.Done():
				return // выход при отмене контекста
			}
		}
	}()

	//создадим прослушиватель порта для grpc(server)
	listener, err := net.Listen("tcp", "0.0.0.0:5001")
	if err != nil {
		logrus.Fatalf("grpc server initialized error: %s", err.Error())
	}
	defer listener.Close()

	// Инициализация зависимостей
	repo := repository.NewPostgresRepository(pool)
	uc := usecase.NewItemUC(repo)
	httpsrv := httpserver.NewHttpServer(uc)
	grpcsrv := grpcserver.NewGrpcServer(uc)

	//инициализируем grpc(server)
	grpcServer := grpc.NewServer()

	//зарегистрируем его на listener
	gen.RegisterCatalogServiceServer(grpcServer, grpcsrv)

	//запустим в горутине чтобы не мешал основному потоку
	go func() {
		if err := grpcServer.Serve(listener); err != nil {
			logrus.Errorf("grpc server starting error: %s", err.Error())
		}
		logrus.Info("grpc server stoped")
	}()

	// Настройка HTTP-сервера
	router := gin.Default()
	router.POST("/items", httpsrv.AddItemHandler)
	router.GET("/items", httpsrv.GetItemsHandler)
	router.GET("/items/:item_id", httpsrv.GetItemsHandler)
	router.DELETE("/items/:item_id", httpsrv.DeleteItemHandler)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	// Запуск сервера в горутине с graceful shutdown
	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			logrus.Error("HTTP server error:", err)
		}
	}()

	go func() {

	}()

	// Ожидание сигналов завершения
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-quit:
		grpcServer.GracefulStop()
		logrus.Info("Shutting down server...")
	case <-ctx.Done():
		logrus.Info("Shutting down due to database error...")
	}

	// Graceful shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logrus.Error("Server shutdown error:", err)
	}

	logrus.Info("Service stopped")
}
