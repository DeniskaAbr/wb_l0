package order

import (
	"context"
	"wb_l0/internal/models"
)

type PGRepository interface {
	Create(ctx context.Context, order *models.Order) (*models.Order, error)
	CreateBatch(ctx context.Context, order *models.Order) error
	GetByUID(ctx context.Context, orderUID string) (*models.Order, error)
	GetFullByUID(ctx context.Context, orderUID string) (*models.Order, error)
	GetOrdersCount(ctx context.Context) (int, error)
	GetOrderUIDs(ctx context.Context) ([]string, error)
}

type CacheRepository interface {
	Set(ctx context.Context, order *models.Order) error
	GetByUID(ctx context.Context, orderUID string) (*models.Order, error)
	Delete(ctx context.Context, orderUID string)
	GetSize() int
}
