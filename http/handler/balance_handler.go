package handler

import (
	"github.com/kataras/iris"
	"github.com/nawa/cryptoexchange-wallet-info/usecase"
)

type BalanceHandler struct {
	balanceUsecase usecase.BalanceUsecase
}

func NewBalanceHandler(balanceUsecase usecase.BalanceUsecase) *BalanceHandler {
	return &BalanceHandler{
		balanceUsecase: balanceUsecase,
	}
}

func (h *BalanceHandler) Get(ctx iris.Context) {
	_, err := ctx.WriteString("some value")
	if err != nil {
		panic(err)
	}
}
