package pub_server

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"wb_l0/config"
	"wb_l0/internal/order/delivery/nats"

	"github.com/nats-io/stan.go"
)

const (
	sendOrderSubject = "order:send"
)

type pub_server struct {
	log      *log.Logger
	cfg      *config.Config
	natsConn stan.Conn
}

func NewServer(
	log *log.Logger,
	cfg *config.Config,
	natsConn stan.Conn,
) *pub_server {
	return &pub_server{log: log, cfg: cfg, natsConn: natsConn}
}

func (ps *pub_server) Run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		orderPublisher := nats.NewPublisher(ps.natsConn, ps.log)
		orderPublisher.Run()
	}()

	////***
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	select {
	case v := <-quit:
		log.Fatalf("signal.Notify: %v", v)
	case done := <-ctx.Done():
		log.Fatalf("ctx.Done: %v", done)
	}

	log.Println("Server Exited Property")

	return nil

}
