package repository

import (
	"context"
	"encoding/json"
	"log"
	"wb_l0/internal/cache"
	"wb_l0/internal/models"

	"github.com/pkg/errors"
)

type orderCacheRepository struct {
	cache *cache.Cache
}

func (o *orderCacheRepository) GetSize() int {
	return len(o.cache.Data)
}

func NewOrderCacheRepository(cache *cache.Cache) *orderCacheRepository {
	data := make(map[string]string)
	cache.Data = data

	return &orderCacheRepository{cache: cache}
}

func (o *orderCacheRepository) Set(ctx context.Context, order *models.Order) error {
	orderBytes, err := json.Marshal(order)
	if err != nil {
		return errors.Wrap(err, "orderCacheRepository.Marshal")
	}

	o.cache.Put(order.OrderUID, string(orderBytes))
	log.Println("Cache len:")
	log.Println(len(o.cache.Data))
	return nil
}

func (o *orderCacheRepository) GetByUID(ctx context.Context, orderUID string) (*models.Order, error) {
	result, err := o.cache.Get(orderUID)
	if err != nil {
		return nil, errors.Wrap(err, "orderCacheRepository.cache.Get")
	}

	var res models.Order
	if err := json.Unmarshal([]byte(result), &res); err != nil {
		return nil, errors.Wrap(err, "json.Unmarshal")
	}

	return &res, nil
}

func (o *orderCacheRepository) Delete(ctx context.Context, orderUID string) {
	o.cache.Del(orderUID)
}
