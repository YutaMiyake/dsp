package main

import (
	"app/common"
	"app/delivery"
	"app/infra"
	"app/repository"
	"app/service"
	"fmt"
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
	"github.com/aerospike/aerospike-client-go"
	"app/models"
	"io/ioutil"
	"encoding/json"
)

func main() {
	server := New()
	defer server.Close()

	server.Init()
	server.Run()
}

type Server struct {
	Engine   fasthttprouter.Router
	AeroConn *aerospike.Client
	AdList []string
	AdInfo map[string]*models.Advertiser
}

func New() *Server {
	r := fasthttprouter.New()
	return &Server{Engine: *r}
}

func (s *Server) Run() {
	s.Route()
	common.Logger.Fatal(fasthttp.ListenAndServe(common.Conf.Dsp.Addr, s.Engine.Handler))
}

func (s *Server) Init() {
	common.SetupLogger()
	common.SetupEnv()
}

func (s *Server) Close() {
	s.AeroConn.Close()
}

func (s *Server) Route() {

	s.AdList, s.AdInfo = loadAeroInfo()

    s.AeroConn = infra.NewAero()
	aeroRepo := repository.NewAeroRepository(s.AeroConn)
	logicSvc := service.NewLogicService(aeroRepo,s.AdList,s.AdInfo)
	handlers := delivery.NewHandlers(logicSvc, aeroRepo)

	s.Engine.GET("/health-check", func(ctx *fasthttp.RequestCtx) {
		common.Logger.Info("health-check")
		fmt.Fprintf(ctx, "ok")
	})
	s.Engine.POST("/bid", handlers.XgBoost)
	s.Engine.POST("/bidPlus", handlers.XgBoostPlusCRate)
	s.Engine.POST("/floor", handlers.FloorPricePlus3)
	s.Engine.POST("/bakugai5", handlers.Bakugai5)
	s.Engine.GET("/aerotest", handlers.AeroTest)
}

func loadAeroInfo() ([]string, map[string]*models.Advertiser) {
	common.Logger.Info("Loading advertiser info")

	// if there data exist in aerospike, read it

	// otherwise, from file

	raw, err := ioutil.ReadFile("resources/budgets.json")
	if err != nil {
		common.Logger.Error(err)
		fmt.Println(err.Error())
		return nil, nil
	}
	result := make([]models.Advertiser, 0)
	err = json.Unmarshal(raw, &result)
	if err != nil {
		common.Logger.Error(err)
		fmt.Println(err.Error())
		return nil, nil
	}

	adInfo := make(map[string]*models.Advertiser,0)
	adList := make([]string,0)
	for idx, ad := range result {
		adInfo[ad.Id] = &result[idx]
		adList = append(adList, ad.Id)
	}
	return adList, adInfo
}