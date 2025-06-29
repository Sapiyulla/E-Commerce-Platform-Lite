package domain

import (
	dto "catalog-service/internal/application/DTO"

	"github.com/google/uuid"
)

type Item struct {
	ID       string
	Name     string
	Category string
	Price    float32
}

func NewItem(name, category string, price float32) *Item {
	id := uuid.New().String()
	return &Item{ID: id, Name: name, Category: category, Price: price}
}

type ItemRepository interface {
	AddItem(*Item) *dto.Status
	GetItem(id string) (*Item, *dto.Status)
	GetItems(category ...string) ([]*Item, *dto.Status)
	DeleteItem(string) *dto.Status
}
