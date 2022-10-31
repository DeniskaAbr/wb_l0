package v1

import (
	"log"
	"net/http"

	"wb_l0/internal/order"

	"github.com/gin-gonic/gin"
)

type orderHandlers struct {
	group   *gin.RouterGroup
	orderUC order.UseCase
	log     *log.Logger
}

func NewOrderHandlers(group *gin.RouterGroup, orderUC order.UseCase, log *log.Logger) *orderHandlers {
	return &orderHandlers{group: group, orderUC: orderUC, log: log}
}

func (h *orderHandlers) GetByUID() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		orderUID := ctx.Param("order_uid")
		if orderUID == "" {
			h.log.Print("order_uid is empty")
			ctx.JSON(http.StatusBadRequest, gin.H{"msg": "empty value not accepted"})
		} else {
			m, err := h.orderUC.GetByUID(ctx, orderUID)
			if err != nil {
				ctx.JSON(http.StatusNotFound, gin.H{"msg": "order with UID value not found"})
			} else {
				ctx.JSON(http.StatusOK, gin.H{"status": 200, "data": m})
			}
		}
	}
}
