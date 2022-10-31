package models

type Order struct {
	OrderUID          string `json:"order_uid," fake:"{regex:[0-9a-z]{19}}" validate:"required,min=19,max=19"`
	TrackNumber       string `json:"track_number" fake:"{regex:[A-Z]{14}}" validate:"required,min=14,max=14"`
	Entry             string `json:"entry" fake:"{regex:[A-Z]{4}}" validate:"required,min=4,max=4 "`
	Locale            string `json:"locale" fake:"{randomstring:[ru,en,it]}" validate:"oneof=ru en it"`
	InternalSignature string `json:"internal_signature" fake:"{lettern:8}" validate:"required,min=8,max=8"`
	CustomerId        string `json:"customer_id" fake:"{lettern:4}"   validate:"required,min=4,max=4"`
	DeliveryService   string `json:"delivery_service" fake:"{lettern:5}"   validate:"required,min=5,max=5"`
	Shardkey          string `json:"shardkey" fake:"{regex:[0-9]{2}}" validate:"required,max=2"`
	SmId              int    `json:"sm_id" fake:"{number:1, 5000}" validate:"required,gte=0,lte=5000"`
	DateCreated       string `json:"date_created" fake:"{year}-{month}-{day}T{hour}:{minute}:{second}Z" format:"2006-01-02T06:22:19Z" validate:"required"`
	OofShard          string `json:"oof_shard" fake:"{regex:[0-9]{2}}" validate:"required,max=2"`

	Delivery `json:"delivery" `
	Payment  `json:"payment" `
	Items    []Item `json:"items" fakesize:"2"`
}
