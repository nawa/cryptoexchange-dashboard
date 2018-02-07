package handler

import (
	"github.com/kataras/iris"
	"github.com/nawa/cryptoexchange-wallet-info/http/dto"
)

func WriteCustomError(ctx iris.Context, status int, message string) {
	ctx.StatusCode(status)
	error := dto.Error{
		Status:  status,
		Message: message,
	}
	_, err := ctx.JSON(error)
	if err != nil {
		panic(err)
	}
}

func WriteBadRequest(ctx iris.Context, message string) {
	WriteCustomError(ctx, iris.StatusBadRequest, message)
}

func WriteInternalServerError(ctx iris.Context, message string) {
	WriteCustomError(ctx, iris.StatusInternalServerError, message)
}

func WriteNotFound(ctx iris.Context, message string) {
	WriteCustomError(ctx, iris.StatusNotFound, message)
}
