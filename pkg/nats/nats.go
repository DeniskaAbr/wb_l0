package nats

import (
	"log"
	"time"
	"wb_l0/config"

	"github.com/nats-io/stan.go"
)

const (
	connectWait        = time.Second * 30
	pubAckWait         = time.Second * 30
	interval           = 10
	maxOut             = 5
	maxPubAcksInflight = 25
)

func NewNatsConnect(cfg *config.Config, log *log.Logger, clientID string) (stan.Conn, error) {

	return stan.Connect(
		cfg.Nats_cluster_id,
		clientID,
		stan.ConnectWait(connectWait),
		stan.PubAckWait(pubAckWait),
		stan.NatsURL("nats://"+cfg.Nats_hostname+":4222"),
		stan.Pings(interval, maxOut),
		stan.SetConnectionLostHandler(func(_ stan.Conn, reason error) {
			log.Fatalf("Connection lost, reason: %v", reason)
		}),
		stan.MaxPubAcksInflight(maxPubAcksInflight),
	)
}
