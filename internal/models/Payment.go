package models

type Payment struct {
	Transaction  string `json:"transaction" fake:"" validate:"required"`
	RequestId    string `json:"request_id" fake:"" validate:"required"`
	Currency     string `json:"currency" fake:"" validate:"required"`
	Provider     string `json:"provider" fake:"" validate:"required"`
	Amount       int    `json:"amount" fake:"{number:1, 2389212}" validate:"required"`
	PaymentDt    int    `json:"payment_dt" fake:"{number:1, 2389212}" validate:"required"`
	Bank         string `json:"bank" fake:"" validate:"required"`
	DeliveryCost int    `json:"delivery_cost" fake:"{number:1, 2389212}" validate:"required"`
	GoodsTotal   int    `json:"goods_total"fake:"{number:1, 2389212}" validate:"required"`
	CustomFee    int    `json:"custom_fee" fake:"{number:1, 2389212}" validate:"required"`
}
