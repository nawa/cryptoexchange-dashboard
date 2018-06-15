package http

import (
	"github.com/nawa/cryptoexchange-dashboard/http/dto"
	"github.com/nawa/cryptoexchange-dashboard/usecase"

	"github.com/kataras/iris"
)

type OrderHandler struct {
	orderUsecase usecase.OrderUsecases
}

func NewOrderHandler(orderUsecase usecase.OrderUsecases) *OrderHandler {
	return &OrderHandler{
		orderUsecase: orderUsecase,
	}
}

func (h *OrderHandler) GetActiveOrders(ctx iris.Context) {
	mOrders, err := h.orderUsecase.GetActiveOrders()
	if err != nil {
		WriteInternalServerError(ctx, "internal error")
		return
	}

	orderDTO := make([]dto.OrderDTO, len(mOrders))
	for i, o := range mOrders {
		orderDTO[i] = *dto.NewOrderDTO(o)
	}
	_, err = ctx.JSON(orderDTO)
	if err != nil {
		panic(err)
	}
}
