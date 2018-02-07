package handler

import (
	"github.com/kataras/iris"
	"github.com/nawa/cryptoexchange-wallet-info/http/dto"
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

func (h *BalanceHandler) Daily(ctx iris.Context) {
	currency := ctx.URLParam("currency")
	if currency == "" {
		WriteBadRequest(ctx, "'currency' is empty")
		return
	}

	mBalances, err := h.balanceUsecase.FetchDaily(currency)
	if err != nil {
		WriteInternalServerError(ctx, "internal error")
		return
	}

	var balances []dto.BalanceDTO
	for _, b := range mBalances {
		balances = append(balances, *dto.NewBalanceDTO(b))
	}

	if len(balances) == 0 {
		//without this transformation json will equal to "null", but "[]" is expected
		balances = make([]dto.BalanceDTO, 0)
	}

	_, err = ctx.JSON(balances)
	if err != nil {
		panic(err)
	}
}
