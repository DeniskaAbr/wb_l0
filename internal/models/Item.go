package models

type Item struct {
	ChrtId      int    `json:"chrt_id" fake:"{number:1, 5000}" validate:"required"`
	TrackNumber string `json:"track_number" fake:"{regex:[A-Z]{14}}" validate:"required,min=14,max=14"`
	Price       int    `json:"price" fake:"{number:1, 5000000}" validate:"required"`
	Rid         string `json:"rid" fake:"{regex:[0-9a-z]{21}}" validate:"required,max=21"`
	Name        string `json:"name" fake:"{minecraftarmorpart}" validate:"required"`
	Sale        int    `json:"sale" fake:"{number:1, 100}" validate:"required"`
	Size        string `json:"size" fake:"{regex:[0-1]{3}}" validate:"required"`
	TotalPrice  int    `json:"total_price" fake:"{number:1, 5000000}" validate:"required"`
	NmId        int    `json:"nm_id" fake:"{number:1, 2389212}" validate:"required"`
	Brand       string `json:"brand" fake:"{firstname} {lastname} " validate:"required"`
	Status      int    `json:"status" fake:"{regex:[0-1]{3}}" validate:"required,max=999"`
}
