package order

import (
	"context"
	"wb_l0/internal/models"
)

type UseCase interface {
	Create(ctx context.Context, order *models.Order) error
	CreateBatch(ctx context.Context, order *models.Order) error
	GetByUID(ctx context.Context, orderUID string) (*models.Order, error)
}
