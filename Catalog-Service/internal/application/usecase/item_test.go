package usecase

import (
	dto "catalog-service/internal/application/DTO"
	"catalog-service/internal/domain"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockItemRepository struct {
	mock.Mock
}

func NewMockItemRepository() *MockItemRepository {
	return &MockItemRepository{}
}

func (m *MockItemRepository) GetItems(categories ...string) ([]*domain.Item, *dto.Status) {
	args := m.Called(categories)
	return args.Get(0).([]*domain.Item), args.Get(1).(*dto.Status)
}

func (m *MockItemRepository) AddItem(item *domain.Item) *dto.Status {
	args := m.Called(item)
	return args.Get(0).(*dto.Status)
}

func (m *MockItemRepository) GetItem(id string) (*domain.Item, *dto.Status) {
	args := m.Called(id)
	return args.Get(0).(*domain.Item), args.Get(1).(*dto.Status)
}

func (m *MockItemRepository) DeleteItem(id string) *dto.Status {
	args := m.Called(id)
	return args.Get(0).(*dto.Status)
}

func TestItemUseCase_AddItem(t *testing.T) {
	mockRepo := new(MockItemRepository)
	uc := NewItemUC(mockRepo)

	//=============================== mock testing =============================================================
	t.Run("валидные данные", func(t *testing.T) {
		item := &dto.Request{Name: "RC Car", Category: "toys", Price: 1200}

		// Используем MatchedBy для проверки только нужных полей
		mockRepo.On("AddItem", mock.MatchedBy(func(dItem *domain.Item) bool {
			return dItem.Name == item.Name &&
				dItem.Category == item.Category &&
				dItem.Price == item.Price
		})).Return(&dto.Status{Code: 0, Status: "success", Message: "item added to catalog"})

		status := uc.AddItem(item)

		assert.Equal(t, 0, status.Code)
		assert.Equal(t, dto.Success, status.Status)

		mockRepo.AssertCalled(t, "AddItem", mock.MatchedBy(func(dItem *domain.Item) bool {
			return dItem.Name == item.Name &&
				dItem.Category == item.Category &&
				dItem.Price == item.Price
		}))
	})
	t.Run("невалидные данные", func(t *testing.T) {
		status := uc.AddItem(&dto.Request{Name: "", Category: "", Price: 0.1})

		assert.Equal(t, 1, status.Code)
		assert.Equal(t, dto.Error, status.Status)
		assert.Equal(t, errItemFieldEmpty("name").Error(), status.Message)

		mockRepo.AssertNotCalled(t, "AddItem")
	})
}

func TestItemUseCase_GetItems(t *testing.T) {
	mockRepo := NewMockItemRepository()
	uc := NewItemUC(mockRepo)

	//=============================== mock testing =============================================================
	t.Run("получение всех товаров", func(t *testing.T) {
		mockRepo.On("GetItems", []string(nil)).Return([]*domain.Item{}, &dto.Status{Code: 0, Status: dto.Success})

		_, status := uc.GetItems()

		assert.Equal(t, 0, status.Code)
		assert.Equal(t, dto.Success, status.Status)

		mockRepo.AssertCalled(t, "GetItems", []string(nil))
		mockRepo.AssertExpectations(t)
	})
	t.Run("получение товаров по категории", func(t *testing.T) {
		category := "electronics"

		mockRepo.On("GetItems", []string{category}).Return([]*domain.Item{}, &dto.Status{Code: 0, Status: dto.Success})

		_, status := uc.GetItems(category)

		assert.Equal(t, 0, status.Code)
		assert.Equal(t, dto.Success, status.Status)

		mockRepo.AssertCalled(t, "GetItems", []string{category})
		mockRepo.AssertExpectations(t)
	})
	t.Run("более 1 категории", func(t *testing.T) {
		category_1 := "electronics"
		category_2 := "toys"

		_, status := uc.GetItems(category_1, category_2)

		assert.Equal(t, 1, status.Code)
		assert.Equal(t, dto.Error, status.Status)
		assert.Equal(t, errItemCategoryOnly1Word.Error(), status.Message)

		mockRepo.AssertNotCalled(t, "GetItems")
	})
}

func TestItemUseCase_GetItem(t *testing.T) {
	mockRepo := NewMockItemRepository()
	uc := NewItemUC(mockRepo)

	//=============================== mock testing =============================================================
	t.Run("валидное получение по идентификатору", func(t *testing.T) {
		mockRepo.On("GetItem", "e22a1f92-785c-4797-8bfe-9f2d250722cc").Return(&domain.Item{}, &dto.Status{Code: 0, Status: dto.Success})

		_, status := uc.GetItem("e22a1f92-785c-4797-8bfe-9f2d250722cc")

		assert.Equal(t, 0, status.Code)
		assert.Equal(t, dto.Success, status.Status)

		mockRepo.AssertExpectations(t)
	})
	t.Run("невалидный идентификатор", func(t *testing.T) {
		_, status := uc.GetItem("1")

		assert.Equal(t, 1, status.Code)
		assert.Equal(t, dto.Error, status.Status)
		assert.Equal(t, errItemID.Error(), status.Message)

		mockRepo.AssertNotCalled(t, "GetItem")
	})
}

func TestItemUseCase_DeleteItem(t *testing.T) {
	mockRepo := NewMockItemRepository()
	uc := NewItemUC(mockRepo)

	//=============================== mock testing =============================================================
	t.Run("валидное получение по идентификатору", func(t *testing.T) {
		mockRepo.On("DeleteItem", "e22a1f92-785c-4797-8bfe-9f2d250722cc").Return(&dto.Status{Code: 0, Status: dto.Success})

		status := uc.DeleteItem("e22a1f92-785c-4797-8bfe-9f2d250722cc")

		assert.Equal(t, 0, status.Code)
		assert.Equal(t, dto.Success, status.Status)

		mockRepo.AssertExpectations(t)
	})
	t.Run("невалидный идентификатор", func(t *testing.T) {
		status := uc.DeleteItem("1")

		assert.Equal(t, 1, status.Code)
		assert.Equal(t, dto.Error, status.Status)
		assert.Equal(t, errItemID.Error(), status.Message)

		mockRepo.AssertNotCalled(t, "DeleteItem")
	})
}

func TestValidItem(t *testing.T) {
	tests := []struct {
		name    string
		item    *domain.Item
		wantErr bool
	}{
		{
			name:    "valid item",
			item:    domain.NewItem("valid name", "validcategory", 10.00),
			wantErr: false,
		},
		{
			name: "empty name",
			item: &domain.Item{
				Name:     "",
				Category: "validcategory",
				Price:    10.00,
			},
			wantErr: true,
		},
		{
			name: "too short name",
			item: &domain.Item{
				Name:     "a",
				Category: "validcategory",
				Price:    10.00,
			},
			wantErr: true,
		},
		{
			name: "too long name",
			item: &domain.Item{
				Name:     "this name is way too long for validation",
				Category: "validcategory",
				Price:    10.00,
			},
			wantErr: true,
		},
		{
			name: "empty category",
			item: &domain.Item{
				Name:     "valid name",
				Category: "",
				Price:    10.00,
			},
			wantErr: true,
		},
		{
			name: "category with space",
			item: &domain.Item{
				Name:     "valid name",
				Category: "invalid category",
				Price:    10.00,
			},
			wantErr: true,
		},
		{
			name: "too short category",
			item: &domain.Item{
				Name:     "valid name",
				Category: "cat",
				Price:    10.00,
			},
			wantErr: true,
		},
		{
			name: "too long category",
			item: &domain.Item{
				Name:     "valid name",
				Category: "thiscategoryistoolong",
				Price:    10.00,
			},
			wantErr: true,
		},
		{
			name: "category with uppercase",
			item: &domain.Item{
				Name:     "valid name",
				Category: "InvalidCategory",
				Price:    10.00,
			},
			wantErr: true,
		},
		{
			name: "price too low",
			item: &domain.Item{
				Name:     "valid name",
				Category: "validcategory",
				Price:    0.99,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validItem(tt.item); err != nil {
				if !tt.wantErr {
					t.Errorf("\nError: %s", err.Error())
				}
			}
		})
	}
}
