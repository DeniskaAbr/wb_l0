package v1

func (h *orderHandlers) MapRoutes() {
	h.group.GET("/:order_uid", h.GetByUID())
}
