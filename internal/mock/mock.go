package mock

import (
	"fmt"
	"wb_l0/internal/models"

	"github.com/brianvoe/gofakeit/v6"
)

func CreateRandomOrder() *models.Order {

	var o models.Order
	gofakeit.Seed(0)
	gofakeit.Struct(&o)
	fmt.Println(o)

	return &o
}
