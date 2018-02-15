package handler

import (
	"github.com/kataras/iris"
	"github.com/nawa/cryptoexchange-dashboard/http/dto"
	"github.com/nawa/cryptoexchange-dashboard/usecase"
)

type BalanceHandler struct {
	balanceUsecase usecase.BalanceUsecase
}

func NewBalanceHandler(balanceUsecase usecase.BalanceUsecase) *BalanceHandler {
	return &BalanceHandler{
		balanceUsecase: balanceUsecase,
	}
}

func (h *BalanceHandler) Hourly(ctx iris.Context) {
	hours, err := ctx.Params().GetInt("hours")
	if err != nil {
		WriteBadRequest(ctx, "':hours' is wrong")
		return
	}
	if hours <= 0 {
		WriteBadRequest(ctx, "':hours' is <= 0")
		return
	}

	currency := ctx.URLParam("currency")
	if currency == "" {
		WriteBadRequest(ctx, "'currency' is empty")
		return
	}

	mBalances, err := h.balanceUsecase.FetchHourly(currency, hours)
	if err != nil {
		WriteInternalServerError(ctx, "internal error")
		return
	}

	curBalancesDTO := []dto.CurrencyBalanceDTO{}
	for _, b := range mBalances {
		curBalancesDTO = append(curBalancesDTO, *dto.NewBalanceDTO(b))
	}

	balanceDTO := dto.BalanceDTO{}
	balanceDTO.Add(currency, curBalancesDTO...)

	_, err = ctx.JSON(balanceDTO)
	if err != nil {
		panic(err)
	}
}

func (h *BalanceHandler) Weekly(ctx iris.Context) {
	currency := ctx.URLParam("currency")
	if currency == "" {
		WriteBadRequest(ctx, "'currency' is empty")
		return
	}

	mBalances, err := h.balanceUsecase.FetchWeekly(currency)
	if err != nil {
		WriteInternalServerError(ctx, "internal error")
		return
	}

	curBalancesDTO := []dto.CurrencyBalanceDTO{}
	for _, b := range mBalances {
		curBalancesDTO = append(curBalancesDTO, *dto.NewBalanceDTO(b))
	}

	balanceDTO := dto.BalanceDTO{}
	balanceDTO.Add(currency, curBalancesDTO...)

	_, err = ctx.JSON(balanceDTO)
	if err != nil {
		panic(err)
	}
}

func (h *BalanceHandler) Monthly(ctx iris.Context) {
	currency := ctx.URLParam("currency")
	if currency == "" {
		WriteBadRequest(ctx, "'currency' is empty")
		return
	}

	mBalances, err := h.balanceUsecase.FetchMonthly(currency)
	if err != nil {
		WriteInternalServerError(ctx, "internal error")
		return
	}

	curBalancesDTO := []dto.CurrencyBalanceDTO{}
	for _, b := range mBalances {
		curBalancesDTO = append(curBalancesDTO, *dto.NewBalanceDTO(b))
	}

	balanceDTO := dto.BalanceDTO{}
	balanceDTO.Add(currency, curBalancesDTO...)

	_, err = ctx.JSON(balanceDTO)
	if err != nil {
		panic(err)
	}
}

func (h *BalanceHandler) All(ctx iris.Context) {
	currency := ctx.URLParam("currency")
	if currency == "" {
		WriteBadRequest(ctx, "'currency' is empty")
		return
	}

	mBalances, err := h.balanceUsecase.FetchAll(currency)
	if err != nil {
		WriteInternalServerError(ctx, "internal error")
		return
	}

	curBalancesDTO := []dto.CurrencyBalanceDTO{}
	for _, b := range mBalances {
		curBalancesDTO = append(curBalancesDTO, *dto.NewBalanceDTO(b))
	}

	balanceDTO := dto.BalanceDTO{}
	balanceDTO.Add(currency, curBalancesDTO...)

	_, err = ctx.JSON(balanceDTO)
	if err != nil {
		panic(err)
	}
}

func (h *BalanceHandler) ActiveCurrencies(ctx iris.Context) {
	mBalances, err := h.balanceUsecase.GetActiveCurrencies()
	if err != nil {
		WriteInternalServerError(ctx, "internal error")
		return
	}

	balanceDTO := dto.BalanceDTO{}
	for _, b := range mBalances {
		balanceDTO.Add(b.Currency, *dto.NewBalanceDTO(b))
	}
	_, err = ctx.JSON(balanceDTO)
	if err != nil {
		panic(err)
	}
}
