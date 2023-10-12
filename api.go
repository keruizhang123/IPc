package ipconf

import (
	"context"

	"github.com/KRZ/ipconf/domain"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

type Response struct {
	Message string      `json:"message"`
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
}

// GetIpInfoList API Adaptive application layer
func GetIpInfoList(c context.Context, ctx *app.RequestContext) {
	defer func() {
		if err := recover(); err != nil {
			ctx.JSON(consts.StatusBadRequest, utils.H{"err": err})
		}
	}()
	// Step0: Build customer request information
	ipConfCtx := domain.BuildIpConfContext(&c, ctx)
	// Step1: Perform ip scheduling
	eds := domain.Dispatch(ipConfCtx)
	// Step2: Top5 is returned according to the score
	ipConfCtx.AppCtx.JSON(consts.StatusOK, packRes(top5Endports(eds)))
}
