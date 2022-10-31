package repository

import (
	"context"
	"fmt"
	"log"
	"wb_l0/internal/models"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
)

type orderPGRepository struct {
	db *pgxpool.Pool
}

func NewOrderPGRepository(db *pgxpool.Pool) *orderPGRepository {
	return &orderPGRepository{db: db}
}

func (o *orderPGRepository) CreateBatch(ctx context.Context, order *models.Order) error {

	batch := &pgx.Batch{}

	batch.Queue(
		createOrderQuery,
		&order.OrderUID,
		&order.TrackNumber,
		&order.Entry,
		&order.Locale,
		&order.InternalSignature,
		&order.CustomerId,
		&order.DeliveryService,
		&order.Shardkey,
		&order.SmId,
		&order.DateCreated,
		&order.OofShard)

	batch.Queue(
		createDeliveryQuery,
		&order.OrderUID,
		&order.Name,
		&order.Phone,
		&order.Zip,
		&order.City,
		&order.Address,
		&order.Region,
		&order.Email)

	batch.Queue(
		createPaymentsQuery,
		&order.OrderUID,
		&order.Payment.Transaction,
		&order.Payment.RequestId,
		&order.Payment.Currency,
		&order.Payment.Provider,
		&order.Payment.Amount,
		&order.Payment.PaymentDt,
		&order.Payment.Bank,
		&order.Payment.DeliveryCost,
		&order.Payment.GoodsTotal,
		&order.Payment.CustomFee)

	for _, v := range order.Items {
		batch.Queue(
			createItemQuery,
			&order.OrderUID,
			&v.ChrtId,
			&v.TrackNumber,
			&v.Price,
			&v.Rid,
			&v.Name,
			&v.Sale,
			&v.Size,
			&v.TotalPrice,
			&v.NmId,
			&v.Brand,
			&v.Status)
	}

	br := o.db.SendBatch(ctx, batch)

	for i := 0; i < batch.Len(); i++ {
		_, err := br.Exec()
		if err != nil {
			return err
		}
	}

	err := br.Close()
	if err != nil {
		return err
	}

	return nil

}

func (o *orderPGRepository) Create(ctx context.Context, order *models.Order) (*models.Order, error) {

	var ord models.Order

	if err := o.db.QueryRow(
		ctx,
		createOrderQuery,
		&order.OrderUID,
		&order.TrackNumber,
		&order.Entry,
		&order.Locale,
		&order.InternalSignature,
		&order.CustomerId,
		&order.DeliveryService,
		&order.Shardkey,
		&order.SmId,
		&order.DateCreated,
		&order.OofShard,
	).Scan(
		&ord.OrderUID,
		&ord.TrackNumber,
		&ord.Entry,
		&ord.Locale,
		&ord.InternalSignature,
		&ord.CustomerId,
		&ord.DeliveryService,
		&ord.Shardkey,
		&ord.SmId,
		&ord.DateCreated,
		&ord.OofShard,

		// &odr.Delivery,
		// &odr.Payment,
		// &odr.Items,
	); err != nil {
		return nil, errors.Wrap(err, "Scan")
	}

	return &ord, nil
}

func (e *orderPGRepository) GetByUID(ctx context.Context, orderUID string) (*models.Order, error) {
	var ord models.Order
	if err := e.db.QueryRow(ctx, getOrderByOrderUIDQuery, orderUID).Scan(
		&ord.OrderUID,
		&ord.TrackNumber,
		&ord.Entry,
		&ord.Locale,
		&ord.InternalSignature,
		&ord.CustomerId,
		&ord.DeliveryService,
		&ord.Shardkey,
		&ord.SmId,
		&ord.DateCreated,
		&ord.OofShard,
	); err != nil {
		return nil, errors.Wrap(err, "Scan")
	}

	return &ord, nil
}

func (e *orderPGRepository) GetFullByUID(ctx context.Context, orderUID string) (*models.Order, error) {
	var ord models.Order

	if err := e.db.QueryRow(ctx, getOrderByOrderUIDQuery, orderUID).Scan(
		&ord.OrderUID,
		&ord.TrackNumber,
		&ord.Entry,
		&ord.Locale,
		&ord.InternalSignature,
		&ord.CustomerId,
		&ord.DeliveryService,
		&ord.Shardkey,
		&ord.SmId,
		&ord.DateCreated,
		&ord.OofShard,
	); err != nil {
		log.Println("1")
		fmt.Println(err)
		return nil, errors.Wrap(err, "Scan")
	}

	if err := e.db.QueryRow(ctx, getDeliveryByOrderUIDQuery, orderUID).Scan(
		&ord.Delivery.Name,
		&ord.Delivery.Phone,
		&ord.Delivery.Zip,
		&ord.Delivery.City,
		&ord.Delivery.Address,
		&ord.Delivery.Region,
		&ord.Delivery.Email,
	); err != nil {
		log.Println("2")
		fmt.Println(err)
		return nil, errors.Wrap(err, "Scan")
	}

	if err := e.db.QueryRow(ctx, getPaymentByOrderUIDQuery, orderUID).Scan(
		&ord.Payment.Transaction,
		&ord.Payment.RequestId,
		&ord.Payment.Currency,
		&ord.Payment.Provider,
		&ord.Payment.Amount,
		&ord.Payment.PaymentDt,
		&ord.Payment.Bank,
		&ord.Payment.DeliveryCost,
		&ord.Payment.GoodsTotal,
		&ord.Payment.CustomFee,
	); err != nil {
		log.Println("3")
		fmt.Println(err)
		return nil, errors.Wrap(err, "Scan")
	}

	rows, err := e.db.Query(ctx, getItemsByOrderUIDQuery, orderUID)
	if err != nil {
		return nil, errors.Wrap(err, "Scan")
	}

	item := models.Item{}

	for rows.Next() {
		err := rows.Scan(
			&item.ChrtId,
			&item.TrackNumber,
			&item.Price,
			&item.Rid,
			&item.Name,
			&item.Sale,
			&item.Size,
			&item.TotalPrice,
			&item.NmId,
			&item.Brand,
			&item.Status,
		)
		if err != nil {
			log.Println("4")
			fmt.Println(err)
			return nil, errors.Wrap(err, "Scan")
		}

		ord.Items = append(ord.Items, item)
	}

	return &ord, nil
}

func (e *orderPGRepository) GetOrdersCount(ctx context.Context) (int, error) {
	var count int

	if err := e.db.QueryRow(ctx, orderCountsQuery).Scan(
		&count,
	); err != nil {
		log.Println("1")
		fmt.Println(err)
		return 0, errors.Wrap(err, "Scan")
	}

	return count, nil
}

func (e *orderPGRepository) GetOrderUIDs(ctx context.Context) ([]string, error) {
	var orderUIDs []string

	rows, err := e.db.Query(ctx, orderUIDsQuery)
	if err != nil {
		return []string{}, errors.Wrap(err, "Scan")
	}

	var orderUID string

	for rows.Next() {
		err := rows.Scan(
			&orderUID,
		)
		if err != nil {
			return []string{}, errors.Wrap(err, "Scan")
		}

		orderUIDs = append(orderUIDs, orderUID)
	}

	return orderUIDs, nil
}
