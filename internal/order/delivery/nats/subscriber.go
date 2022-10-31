package nats

import (
	"context"
	"encoding/json"
	"log"
	"sync"

	"wb_l0/internal/models"
	"wb_l0/internal/order"

	"github.com/go-playground/validator/v10"
	"github.com/nats-io/stan.go"
)

const ()

type orderSubscriber struct {
	log      *log.Logger
	stanConn stan.Conn
	orderUC  order.UseCase
	validate *validator.Validate
}

func NewOrderSubscriber(stanConn stan.Conn, log *log.Logger, orderUC order.UseCase, validate *validator.Validate) *orderSubscriber {
	return &orderSubscriber{stanConn: stanConn, log: log, orderUC: orderUC, validate: validate}
}

func (s *orderSubscriber) Subscribe(subject, qgroup string, workerNum int, cb stan.MsgHandler) {
	s.log.Printf("Subscribing to Subject: %v, group: %v", subject, qgroup)
	wg := &sync.WaitGroup{}

	for i := 0; i <= workerNum; i++ {
		wg.Add(1)
		go s.runWorker(
			wg,
			i,
			s.stanConn,
			subject,
			qgroup,
			cb,
			// stan.SetManualAckMode(),
			// stan.AckWait(ackWait),
			// stan.DurableName(durableName),
			// stan.MaxInflight(maxInflight),
			// stan.DeliverAllAvailable(),
		)
	}
	wg.Wait()
}

func (s *orderSubscriber) runWorker(
	wg *sync.WaitGroup,
	workerID int,
	conn stan.Conn,
	subject string,
	qgroup string,
	cb stan.MsgHandler,
	opts ...stan.SubscriptionOption,
) {
	s.log.Printf("Subscribing worker: %v, subject: %v, qgroup: %v", workerID, subject, qgroup)
	defer wg.Done()

	_, err := conn.QueueSubscribe(subject, qgroup, cb, opts...)
	if err != nil {
		s.log.Printf("WorkerID: %v, QueueSubscribe: %v", workerID, err)
		if err := conn.Close(); err != nil {
			s.log.Printf("WorkerID: %v, conn.Close error: %v", workerID, err)
		}
	}

}

func (s *orderSubscriber) Run(ctx context.Context) {
	go s.Subscribe(createOrderSubject, orderGroupName, createOrderWorkers, s.processCreateOrder(ctx))
}

func (s *orderSubscriber) processCreateOrder(ctx context.Context) stan.MsgHandler {
	return func(msg *stan.Msg) {
		s.log.Printf("subscriber process Create Order: %s", msg.String())

		var m models.Order

		// avlidator
		err := s.validate.Struct(m)

		if err != nil {
			log.Print("Data validate error")
			if _, ok := err.(*validator.InvalidValidationError); ok {
				log.Println(err)
				return
			}
		}

		if err := json.Unmarshal(msg.Data, &m); err != nil {
			s.log.Printf("json.Unmarshal: %v", err)
			return
		}

		// if err := s.orderUC.Create(ctx, &m); err != nil {
		// 	s.log.Printf("orderUC.Create : %v", err)
		// 	return
		// }

		if err := s.orderUC.CreateBatch(ctx, &m); err != nil {
			s.log.Printf("orderUC.Create : %v", err)
			return
		}
	}
}
