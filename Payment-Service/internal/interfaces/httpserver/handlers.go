package httpserver

import (
	"context"
	"net/http"
	dto "payment-service/internal/application/DTO"
	"payment-service/internal/application/usecases"
	"time"

	"github.com/gin-gonic/gin"
)

type HttpServer struct {
	usecase *usecases.PaymentUseCase
}

func NewHttpServer(uc *usecases.PaymentUseCase) *HttpServer {
	return &HttpServer{usecase: uc}
}

func (h *HttpServer) PayHandler(c *gin.Context) {
	order_id := c.Param("order_id")
	if order_id == "" {
		c.JSON(http.StatusBadRequest, dto.New("error", "order_id param not found"))
		return
	}
	reqCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	c.JSON(200, h.usecase.Pay(reqCtx, order_id))
}
