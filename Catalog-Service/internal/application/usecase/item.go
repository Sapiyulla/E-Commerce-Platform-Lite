package usecase

import (
	dto "catalog-service/internal/application/DTO"
	"catalog-service/internal/domain"
	"errors"
	"strconv"
	"strings"
	"unicode"

	"github.com/google/uuid"
)

var (
	errItemCategoryLowerCase = errors.New("category must contain lowercase letters")
	errItemCategoryOnly1Word = errors.New("category can have only 1 word")
	errItemPrice             = errors.New("price must be at least 1 unit")
	errItemID                = errors.New("item ID is invalid")
	errItemFieldEmpty        = func(_field string) error { return errors.New(_field + " cannot be empty") }
	errItemFieldMinLen       = func(_field string, _len int) error {
		return errors.New(_field + " length cannot be less than " + strconv.Itoa(_len) + " charasters")
	}
	errItemFieldMaxLen = func(_field string, _len int) error {
		return errors.New(_field + " length cannot be more than " + strconv.Itoa(_len) + " charasters")
	}
)

type ItemUseCase struct {
	repo domain.ItemRepository
}

func NewItemUC(repo domain.ItemRepository) *ItemUseCase {
	return &ItemUseCase{repo: repo}
}

func (uc *ItemUseCase) AddItem(item *dto.Request) *dto.Status {
	domainItem := domain.NewItem(item.Name, item.Category, item.Price)
	if err := validItem(domainItem); err != nil {
		return &dto.Status{Code: 1, Status: "error", Message: err.Error()}
	}
	return uc.repo.AddItem(domainItem)
}

func (uc *ItemUseCase) GetItems(category ...string) ([]*domain.Item, *dto.Status) {
	if len(category) > 1 {
		return []*domain.Item{}, &dto.Status{Code: 1, Status: dto.Error, Message: errItemCategoryOnly1Word.Error()}
	}
	return uc.repo.GetItems(category...)
}

func (uc *ItemUseCase) GetItem(id string) (*domain.Item, *dto.Status) {
	if err := uuid.Validate(id); err != nil {
		return nil, &dto.Status{Code: 1, Status: dto.Error, Message: errItemID.Error()}
	}
	return uc.repo.GetItem(id)
}

func (uc *ItemUseCase) DeleteItem(id string) *dto.Status {
	if err := uuid.Validate(id); err != nil {
		return &dto.Status{Code: 1, Status: dto.Error, Message: errItemID.Error()}
	}
	return uc.repo.DeleteItem(id)
}

func validItem(item *domain.Item) error {
	if err := validName(item.Name); err != nil {
		return err
	}
	if err := validCategory(item.Category); err != nil {
		return err
	}
	if item.Price < 10 {
		return errItemPrice
	}
	return nil
}

func validName(name string) error {
	rune_name := []rune(name)
	minLen := 2
	maxLen := 20
	if name == "" {
		return errItemFieldEmpty("name")
	}
	if len(rune_name) < minLen {
		return errItemFieldMinLen("name", minLen)
	}
	if len(rune_name) > maxLen {
		return errItemFieldMaxLen("name", maxLen)
	}
	return nil
}

func validCategory(category string) error {
	rune_category := []rune(category)
	minLen := 4
	maxLen := 15
	if category == "" {
		return errItemFieldEmpty("category")
	}
	if len(strings.Split(category, " ")) > 1 {
		return errItemCategoryOnly1Word
	}
	if len(rune_category) < minLen {
		return errItemFieldMinLen("category", minLen)
	}
	if len(rune_category) > maxLen {
		return errItemFieldMaxLen("category", maxLen)
	}
	for _, char := range category {
		if !unicode.IsLower(char) {
			return errItemCategoryLowerCase
		}
	}
	return nil
}
