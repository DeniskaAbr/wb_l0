package nats

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
	"wb_l0/internal/mock"

	"github.com/nats-io/stan.go"
)

type Publisher interface {
	Publish(subject string, data []byte) error
	PublishAsync(subject string, data []byte, ah stan.AckHandler) (string, error)
}

type publisher struct {
	stanConn stan.Conn
	log      *log.Logger
}

// NewPublisher Nats publisher constructor
func NewPublisher(stanConn stan.Conn, log *log.Logger) *publisher {
	return &publisher{stanConn: stanConn, log: log}
}

// Publish Publish will publish to the cluster and wait for an ACK
func (p *publisher) Publish(subject string, data []byte) error {
	log.Printf("Publish data: %v to subject: %v", string(data), subject)
	return p.stanConn.Publish(subject, data)
}

// PublishAsync PublishAsync will publish to the cluster and asynchronously process the ACK or error state.
// It will return the GUID for the message being sent.
func (p *publisher) PublishAsync(subject string, data []byte, ah stan.AckHandler) (string, error) {
	log.Printf("Publish data: %v to subject: %v", string(data), subject)
	return p.stanConn.PublishAsync(subject, data, ah)
}

func (p *publisher) Run() {
	for {

		order := mock.CreateRandomOrder()
		orderBytes, _ := json.Marshal(*order)

		log.Println("Publish new random order")
		err := p.stanConn.Publish("order:create", orderBytes)

		if err != nil {
			fmt.Println(err)
		}

		time.Sleep(3000 * time.Millisecond)
	}

}
