package main

import (
	"log"
	"os"
	"wb_l0/config"
	"wb_l0/internal/cache"
	"wb_l0/internal/postgresql"
	"wb_l0/internal/server"
	"wb_l0/pkg/nats"
)

// func init() {
// 	// only for  testing
// 	os.Setenv("POSTGRES_HOST", "localhost")
// 	os.Setenv("POSTGRES_PORT", "5432")

// 	os.Setenv("POSTGRES_PASSWORD", "market_password")
// 	os.Setenv("POSTGRES_USER", "market_user")
// 	os.Setenv("POSTGRES_DB", "market_db")

// 	os.Setenv("NATS_CLUSTER_ID", "test-cluster")
// 	os.Setenv("NATS_HOSTNAME", "wb_l0-nats-streaming-server")

// }

func main() {

	var config config.Config
	err := config.InitFromEnv()

	if err != nil {
		log.Fatal(err)
	}

	appLogger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	appLogger.Println("Starting subscriber microservice")
	appLogger.Println("Config is loaded")

	natsConn, err := nats.NewNatsConnect(&config, appLogger, "subscriber")
	if err != nil {
		log.Fatalf("NewNatsConnect: %+v", err)
	}

	var ch cache.Cache

	pgxPool, err := postgresql.NewPgxConn(&config)
	if err != nil {
		log.Fatalf("NewPgxConn: %+v", err)
	}

	s := server.NewServer(appLogger, &config, natsConn, pgxPool, &ch)

	log.Fatal(s.Run())
}
