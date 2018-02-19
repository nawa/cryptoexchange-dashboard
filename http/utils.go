package http

import (
	"github.com/kataras/iris"
	"github.com/nawa/cryptoexchange-dashboard/http/dto"
)

func WriteCustomError(ctx iris.Context, status int, message string) {
	ctx.StatusCode(status)
	dtoError := dto.Error{
		Status:  status,
		Message: message,
	}
	_, err := ctx.JSON(dtoError)
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
