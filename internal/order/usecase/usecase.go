package usecase

import (
	"context"
	"fmt"
	"log"
	"wb_l0/internal/cache"
	"wb_l0/internal/models"
	"wb_l0/internal/order"

	"github.com/pkg/errors"
)

type orderUseCase struct {
	log         *log.Logger
	orderPGRepo order.PGRepository
	cacheRepo   order.CacheRepository
}

func NewOrderUseCase(log *log.Logger, orderPGRepo order.PGRepository, cacheRepo order.CacheRepository) *orderUseCase {
	return &orderUseCase{
		log:         log,
		orderPGRepo: orderPGRepo,
		cacheRepo:   cacheRepo,
	}
}

func (o *orderUseCase) Create(ctx context.Context, order *models.Order) error {
	_, err := o.orderPGRepo.Create(ctx, order)
	if err != nil {
		return errors.Wrap(err, "orderPGRepo.Create")
	}
	return nil
}

func (o *orderUseCase) CreateBatch(ctx context.Context, order *models.Order) error {
	err := o.orderPGRepo.CreateBatch(ctx, order)
	if err != nil {
		return errors.Wrap(err, "orderPGRepo.CreateBatch")
	}
	return nil
}

func (o *orderUseCase) GetByUID(ctx context.Context, orderUID string) (*models.Order, error) {

	if o.cacheRepo.GetSize() == 0 {

		orderCounts, err := o.orderPGRepo.GetOrdersCount(ctx)
		if err != nil {
			return nil, errors.Wrap(err, "orderPGRepo.GetByUID")
		}

		fmt.Printf("order count: %v \n", orderCounts)

		oUIDs, err := o.orderPGRepo.GetOrderUIDs(ctx)

		if err != nil {
			return nil, errors.Wrap(err, "orderPGRepo.GetByUID")
		}

		for _, uid := range oUIDs {
			order, err := o.orderPGRepo.GetFullByUID(ctx, uid)
			if err != nil {
				return nil, errors.Wrap(err, "orderPGRepo.GetByUID")
			}
			if err := o.cacheRepo.Set(ctx, order); err != nil {
				o.log.Printf("cache.SetOrder: %v", err)
			}
		}
	}

	cached, err := o.cacheRepo.GetByUID(ctx, orderUID)
	if err != nil && err != cache.Nil {
		o.log.Printf("cacherepo,GetOrderByUID: %v", err)
	}

	if cached != nil {
		return cached, nil
	}

	order, err := o.orderPGRepo.GetFullByUID(ctx, orderUID)
	if err != nil {
		return nil, errors.Wrap(err, "orderPGRepo.GetByUID")
	}
	if err := o.cacheRepo.Set(ctx, order); err != nil {
		o.log.Printf("cache.SetOrder: %v", err)
	}

	return order, nil
}
