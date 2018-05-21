package delivery

import (
	"github.com/valyala/fasthttp"
	"time"
	"encoding/json"
	"app/models"
	"app/common"
	"app/repository"
	"io/ioutil"
	"fmt"
)


var (
	strContentType     = []byte("Content-Type")
	strApplicationJSON = []byte("application/json")
)

type Handlers struct {
	aeroRepo repository.AeroRepository
	AdInfo map[string]*models.Advertiser
}

func NewHandlers(aeroRepo repository.AeroRepository, adInfo map[string]*models.Advertiser) *Handlers {
	return &Handlers{aeroRepo,adInfo}
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

func (h *Handlers) NoticeHandler(ctx *fasthttp.RequestCtx) {
	adId := ctx.UserValue("adId")

	var req models.WinNoticeRequest
	err := json.Unmarshal(ctx.PostBody(), &req)
	if err != nil {
		common.Logger.Error(err)
	}

	adIdStr := adId.(string)

	if req.IsClick == 1 {
		h.aeroRepo.SpentBudget(adIdStr,h.AdInfo[adIdStr].Cpc)
	}

	ctx.Response.SetStatusCode(fasthttp.StatusNoContent)
}

func (h *Handlers) PutAll(ctx *fasthttp.RequestCtx) {
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

func (h *Handlers) PutDummies(ctx *fasthttp.RequestCtx) {
	// load budgets from file
	raw, err := ioutil.ReadFile("resources/test-budgets.json")
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

func (h *Handlers) GetAll(ctx *fasthttp.RequestCtx) {
	result,err := h.aeroRepo.GetAllAdInfo()
	if err != nil {

	}
	doJSONWrite(ctx,fasthttp.StatusOK,result)
}