package server

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"wb_l0/config"
	"wb_l0/internal/cache"
	"wb_l0/internal/order/delivery/nats"
	"wb_l0/internal/order/repository"
	"wb_l0/internal/order/usecase"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/nats-io/stan.go"

	ordersV1 "wb_l0/internal/order/delivery/http/v1"
)

const ()

type server struct {
	log      *log.Logger
	cfg      *config.Config
	natsConn stan.Conn
	pgxPool  *pgxpool.Pool
	gin      *gin.Engine
	cache    *cache.Cache
}

func NewServer(
	log *log.Logger,
	cfg *config.Config,
	natsConn stan.Conn,
	pgxPool *pgxpool.Pool,
	cache *cache.Cache,
) *server {
	gin.SetMode(gin.ReleaseMode)
	return &server{log: log, cfg: cfg, natsConn: natsConn, pgxPool: pgxPool, gin: gin.New(), cache: cache}
}

func (s *server) Run() error {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	orderPgRepo := repository.NewOrderPGRepository(s.pgxPool)
	orderCacheRepo := repository.NewOrderCacheRepository(s.cache)
	orderUC := usecase.NewOrderUseCase(s.log, orderPgRepo, orderCacheRepo)

	var validate *validator.Validate
	validate = validator.New()

	go func() {
		orderSubscriber := nats.NewOrderSubscriber(s.natsConn, s.log, orderUC, validate)
		orderSubscriber.Run(ctx)
	}()

	go func() {
		s.log.Printf("Server is listening on PORT: %s", s.cfg.Http_port)
		s.runHttpServer()
	}()

	s.gin.Use(cors.Default())
	v1 := s.gin.Group("/api/v1", cors.Default())

	orderHandlers := ordersV1.NewOrderHandlers(v1.Group("/order"), orderUC, s.log)
	orderHandlers.MapRoutes()

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
