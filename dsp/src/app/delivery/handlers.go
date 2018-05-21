package delivery

import (
	"app/common"
	"app/models"
	"app/repository"
	"app/service"
	"encoding/json"
	"fmt"
	"github.com/valyala/fasthttp"
	"io/ioutil"
	"time"
)

var (
	strContentType     = []byte("Content-Type")
	strApplicationJSON = []byte("application/json")
)

type Handlers struct {
	logicSvc service.LogicService
	aeroRepo repository.AeroRepository
}

func NewHandlers(logicSvc service.LogicService, aeroRepo repository.AeroRepository) *Handlers {
	return &Handlers{logicSvc, aeroRepo}
}

func doJSONWrite(ctx *fasthttp.RequestCtx, code int, obj interface{}) {
	ctx.Response.Header.SetCanonical(strContentType, strApplicationJSON)
	ctx.Response.SetStatusCode(code)
	start := time.Now()
	if err := json.NewEncoder(ctx).Encode(obj); err != nil {
		elapsed := time.Since(start)
		common.Logger.Error("", elapsed, err.Error(), obj)
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
	}
}

func (h *Handlers) XgBoost(ctx *fasthttp.RequestCtx) {
	//common.Logger.Info("bid")

	var req models.SSPRequest
	err := json.Unmarshal(ctx.PostBody(), &req)
	if err != nil {
		common.Logger.Error(err)
		ctx.Response.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	buy, res, err := h.logicSvc.XgBoost(&req)

	if err != nil {
		common.Logger.Error(err)
		ctx.Response.SetStatusCode(fasthttp.StatusNoContent)
		return
	}

	if buy {
		doJSONWrite(ctx, fasthttp.StatusOK, res)
		return
	}

	ctx.Response.SetStatusCode(fasthttp.StatusNoContent)
	return
}

func (h *Handlers) XgBoostPlusCRate(ctx *fasthttp.RequestCtx) {
	//common.Logger.Info("bid")

	var req models.SSPRequest
	err := json.Unmarshal(ctx.PostBody(), &req)
	if err != nil {
		common.Logger.Error(err)
		ctx.Response.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	buy, res, err := h.logicSvc.XgBoostPlusCRate(&req)

	if err != nil {
		common.Logger.Error(err)
		ctx.Response.SetStatusCode(fasthttp.StatusNoContent)
		return
	}

	if buy {
		doJSONWrite(ctx, fasthttp.StatusOK, res)
		return
	}

	ctx.Response.SetStatusCode(fasthttp.StatusNoContent)
	return
}

func (h *Handlers) FloorPricePlus3(ctx *fasthttp.RequestCtx) {
	//common.Logger.Info("bid")

	var req models.SSPRequest
	err := json.Unmarshal(ctx.PostBody(), &req)
	if err != nil {
		common.Logger.Error(err)
		ctx.Response.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	buy, res, err := h.logicSvc.FloorPricePlus3(&req)

	if err != nil {
		common.Logger.Error(err)
		ctx.Response.SetStatusCode(fasthttp.StatusNoContent)
		return
	}

	if buy {
		doJSONWrite(ctx, fasthttp.StatusOK, res)
		return
	}

	ctx.Response.SetStatusCode(fasthttp.StatusNoContent)
	return
}

func (h *Handlers) Bakugai5(ctx *fasthttp.RequestCtx) {
	//common.Logger.Info("bid")

	var req models.SSPRequest
	err := json.Unmarshal(ctx.PostBody(), &req)
	if err != nil {
		common.Logger.Error(err)
		ctx.Response.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	buy, res, err := h.logicSvc.Bakugai5(&req)

	if err != nil {
		common.Logger.Error(err)
		ctx.Response.SetStatusCode(fasthttp.StatusNoContent)
		return
	}

	if buy {
		doJSONWrite(ctx, fasthttp.StatusOK, res)
		return
	}

	ctx.Response.SetStatusCode(fasthttp.StatusNoContent)
	return
}

func (h *Handlers) AeroTest(ctx *fasthttp.RequestCtx) {
	// load budgets from file
	raw, err := ioutil.ReadFile("resources/budgets.json")
	if err != nil {
		common.Logger.Error(err)
		fmt.Println(err.Error())
		return
	}
	result := make([]models.Advertiser, 0)
	err = json.Unmarshal(raw, &result)
	if err != nil {
		common.Logger.Error(err)
		fmt.Println(err.Error())
		return
	}

	for _, ad := range result {
		common.Logger.Info(ad)

		// put
		err = h.aeroRepo.PutAd(&ad)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		// get
		ret, err := h.aeroRepo.GetAd(ad.Id)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		fmt.Printf("%#v\n", *ret)
	}
}
