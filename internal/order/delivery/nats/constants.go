package nats

import "time"

const (
	ackWait     = 60 * time.Second
	durableName = "microservice-dur"
	maxInflight = 25

	createOrderWorkers = 0

	createOrderSubject = "order:create"
	orderGroupName     = "order_service"
)
