package main

import (
	"app/common"
	"app/infra"
	"fmt"
	"github.com/aerospike/aerospike-client-go"
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
	"app/delivery"
	"app/repository"
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
	AdInfo map[string]*models.Advertiser
}

func New() *Server {
	r := fasthttprouter.New()
	return &Server{Engine: *r}
}

func (s *Server) Run() {
	s.Route()
	common.Logger.Fatal(fasthttp.ListenAndServe(common.Conf.Notice.Addr, s.Engine.Handler))
}

func (s *Server) Init() {
	common.SetupLogger()
	common.SetupEnv()
}

func (s *Server) Close() {
}

func (s *Server) Route() {
	AdInfo := loadAeroInfo()

	s.AeroConn = infra.NewAero()
	aeroRepo := repository.NewAeroRepository(s.AeroConn)
	handlers := delivery.NewHandlers(aeroRepo, AdInfo)

	s.Engine.GET("/health-check", func(ctx *fasthttp.RequestCtx) {
		fmt.Fprintf(ctx, "ok")
	})
	s.Engine.POST("/notice/:adId", handlers.NoticeHandler)
	s.Engine.GET("/putall", handlers.PutAll)
	s.Engine.GET("/putdummies", handlers.PutDummies)
	s.Engine.GET("/getall", handlers.GetAll)
}

func loadAeroInfo() map[string]*models.Advertiser {
	common.Logger.Info("Loading advertiser info")

	// if there data exist in aerospike, read it

	// otherwise, from file

	raw, err := ioutil.ReadFile("resources/budgets.json")
	if err != nil {
		common.Logger.Error(err)
		fmt.Println(err.Error())
		return nil
	}
	result := make([]models.Advertiser, 0)
	err = json.Unmarshal(raw, &result)
	if err != nil {
		common.Logger.Error(err)
		fmt.Println(err.Error())
		return nil
	}

	adInfo := make(map[string]*models.Advertiser,0)
	for idx, ad := range result {
		adInfo[ad.Id] = &result[idx]
	}
	return adInfo
}