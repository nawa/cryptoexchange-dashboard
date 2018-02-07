package handler

import "github.com/kataras/iris"

type BaseHandler struct {
}

func NewBaseHandler() *BaseHandler {
	return &BaseHandler{}
}

func (h *BaseHandler) Ping(ctx iris.Context) {
	_, err := ctx.WriteString("pong")
	if err != nil {
		panic(err)
	}
}
