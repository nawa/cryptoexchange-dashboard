package handler

import (
	"github.com/nawa/cryptoexchange-wallet-info/http/dto"
	"github.com/nawa/cryptoexchange-wallet-info/usecase"

	"github.com/kataras/iris"
)

type OrderHandler struct {
	orderUsecase usecase.OrderUsecase
}

func NewOrderHandler(orderUsecase usecase.OrderUsecase) *OrderHandler {
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

	orderDTO := []dto.OrderDTO{}
	for _, o := range mOrders {
		orderDTO = append(orderDTO, *dto.NewOrderDTO(o))
	}
	_, err = ctx.JSON(orderDTO)
	if err != nil {
		panic(err)
	}
}
