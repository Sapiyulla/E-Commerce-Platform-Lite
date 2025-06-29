package httpserver

import (
	dto "catalog-service/internal/application/DTO"
	"catalog-service/internal/application/usecase"
	"catalog-service/internal/domain"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HttpServer struct {
	usecase *usecase.ItemUseCase
}

func NewHttpServer(usecase *usecase.ItemUseCase) *HttpServer {
	return &HttpServer{usecase: usecase}
}

func (s *HttpServer) AddItemHandler(c *gin.Context) {
	var item *dto.Request
	if err := c.BindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	status := s.usecase.AddItem(item)
	switch status.Code {
	case 0:
		c.JSON(http.StatusCreated, status)
		return
	case 1:
		c.JSON(http.StatusBadRequest, status)
		return
	case 2:
		c.JSON(http.StatusInternalServerError, status)
		return
	}
}

func (s *HttpServer) GetItemsHandler(c *gin.Context) {
	if item_id := c.Param("item_id"); item_id != "" {
		item, status := s.usecase.GetItem(item_id)
		c.JSON(http.StatusOK, struct {
			Info *domain.Item `json:"info"`
			*dto.Status
		}{Info: item, Status: status})
		return
	}
	if category := c.Query("category"); category != "" {
		items, status := s.usecase.GetItems(category)
		switch status.Code {
		case 0:
			c.JSON(http.StatusOK, struct {
				Info []*domain.Item `json:"info"`
				*dto.Status
			}{
				Info:   items,
				Status: status,
			})
			return
		case 1:
			c.JSON(http.StatusBadRequest, status)
			return
		case 2:
			c.JSON(http.StatusInternalServerError, status)
			return
		}
	}
	items, status := s.usecase.GetItems()
	switch status.Code {
	case 0:
		c.JSON(http.StatusOK, struct {
			Info []*domain.Item `json:"info"`
			*dto.Status
		}{
			Info:   items,
			Status: status,
		})
		return
	case 1:
		c.JSON(http.StatusBadRequest, status)
		return
	case 2:
		c.JSON(http.StatusInternalServerError, status)
		return
	}
}

func (s *HttpServer) DeleteItemHandler(c *gin.Context) {
	if item_id := c.Param("item_id"); item_id != "" {
		status := s.usecase.DeleteItem(item_id)
		switch status.Code {
		case 0:
			c.JSON(http.StatusOK, status)
			return
		case 1:
			c.JSON(http.StatusBadRequest, status)
			return
		case 2:
			c.JSON(http.StatusInternalServerError, status)
			return
		}
	}
}
