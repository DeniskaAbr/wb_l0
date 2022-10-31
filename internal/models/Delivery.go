package models

type Delivery struct {
	Name    string `json:"name" fake:"{firstname} {lastname}" validate:"required,max=30"`
	Phone   string `json:"phone" fake:"{phone}" validate:"required,max=20"`
	Zip     string `json:"zip" fake:"{zip}" validate:"required,max=20"`
	City    string `json:"city" fake:"{city}" validate:"required,max=30"`
	Address string `json:"address" fake:"{streetname}" validate:"required,max=20"`
	Region  string `json:"region" fake:"{state}" validate:"required,max=20"`
	Email   string `json:"email" fake:"{email}" validate:"required,max=30"`
}
