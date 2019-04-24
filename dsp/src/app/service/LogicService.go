package service

import (
	"app/common"
	"app/models"
	"app/repository"
	"math/rand"
	"strconv"
	"io/ioutil"
	"fmt"
	"encoding/json"
	"net/http"
	"bytes"
)

type LogicService interface {
	XgBoost(req *models.SSPRequest) (bool, *models.SSPResponse, error)
	XgBoostPlusCRate(req *models.SSPRequest) (bool, *models.SSPResponse, error)
	Bakugai5(req *models.SSPRequest) (bool, *models.SSPResponse, error)

	AdtechCompe(req *models.BidRequest) (bool, *models.BidResponse, error)
}

type logicService struct{
	aeroRepo repository.AeroRepository
	AdList []string
	AdInfo map[string]*models.Advertiser
}

func NewLogicService(aeroRepo repository.AeroRepository, adlist []string, adinfo map[string]*models.Advertiser) LogicService {
	return &logicService{aeroRepo:aeroRepo,AdList:adlist,AdInfo:adinfo}
}

func (w *logicService) XgBoost(req *models.SSPRequest) (bool, *models.SSPResponse, error) {
	//common.Logger.Info("BidOrSkip %s", req)

	num := rand.Int()%20+1
	var advId string
	if num < 10{
		advId = "adv_0"+strconv.Itoa(num)
	}else{
		advId = "adv_"+strconv.Itoa(num)
	}

	// calc consumption rate
	spent, err := w.aeroRepo.GetSpent(advId)
	rate := float64(spent)/float64(w.AdInfo[advId].Budget)
	if err != nil || rate > 0.9 {
		common.Logger.Error(err)
		return false, nil, nil
	}

	// request
	ctr, err := predict(req, num)
	if err != nil {
		return false, nil, nil
	}

	f64, _ := strconv.ParseFloat(ctr.Ctr,64)
	cpm := float64(w.AdInfo[advId].Cpc) * f64 * 1000

	if req.FloorPrice > cpm {
		common.Logger.Infof("CPM %f < %f", cpm, req.FloorPrice)
		return false, nil, nil
	}

	res := models.SSPResponse{
		req.Id,
		cpm,
		advId,
		"http://35.194.99.70/notice/"+advId,
	}
	common.Logger.Infof("Bid Req FP %f Res %s", req.FloorPrice, res)

	return true, &res, nil
}

func (w *logicService) XgBoostPlusCRate(req *models.SSPRequest) (bool, *models.SSPResponse, error) {
	//common.Logger.Info("BidOrSkip %s", req)

	num := rand.Int()%20+1
	var advId string
	if num < 10{
		advId = "adv_0"+strconv.Itoa(num)
	}else{
		advId = "adv_"+strconv.Itoa(num)
	}

	// calc consumption rate
	spent, err := w.aeroRepo.GetSpent(advId)
	rate := float64(spent)/float64(w.AdInfo[advId].Budget)
	if err != nil || rate > 0.9 {
		common.Logger.Error(err)
		return false, nil, nil
	}

	// request
	ctr, err := predict(req, num)
	if err != nil {
		return false, nil, nil
	}

	f64, _ := strconv.ParseFloat(ctr.Ctr,64)
	cpm := float64(w.AdInfo[advId].Cpc) * f64 * 1000 * (1+1-rate)

	if req.FloorPrice > cpm {
		common.Logger.Infof("CPM %f < %f", cpm, req.FloorPrice)
		return false, nil, nil
	}

	res := models.SSPResponse{
		req.Id,
		cpm,
		advId,
		"http://35.194.99.70/notice/"+advId,
	}
	common.Logger.Infof("Bid Req FP %f Res %s", req.FloorPrice, res)

	return true, &res, nil
}

func (w *logicService) AdtechCompe(req *models.BidRequest) (bool, *models.BidResponse, error) {
	//common.Logger.Info("BidOrSkip %s", req)

	num := rand.Int()%20+1
	var advId string
	if num < 10{
		advId = "adv_0"+strconv.Itoa(num)
	}else{
		advId = "adv_"+strconv.Itoa(num)
	}

	// calc consumption rate
	spent, err := w.aeroRepo.GetSpent(advId)
	rate := float64(spent)/float64(w.AdInfo[advId].Budget)
	if err != nil || rate > 0.9 {
		common.Logger.Error(err)
		return false, nil, nil
	}

	res := models.BidResponse{
		req.Id,
		req.FloorPrice+3,
		advId,
		"http://localhost:8083/notice/"+advId,
	}
	common.Logger.Infof("[Bid Request] FP %f | Response %s", req.FloorPrice, res)

	return true, &res, nil
}

func (w *logicService) Bakugai5(req *models.SSPRequest) (bool, *models.SSPResponse, error) {

	advId := "adv_05"
	num := 5

	// calc consumption rate
	spent, err := w.aeroRepo.GetSpent(advId)
	rate := float64(spent)/float64(w.AdInfo[advId].Budget)
	if err != nil || rate > 0.9 {
		common.Logger.Error(err)
		return false, nil, nil
	}

	// request
	ctr, err := predict(req, num)
	if err != nil {
		return false, nil, nil
	}

	f64, _ := strconv.ParseFloat(ctr.Ctr,64)
	cpm := float64(w.AdInfo[advId].Cpc) * f64 * 1000 * 30

	if req.FloorPrice > cpm {
		common.Logger.Infof("CPM %f < %f", cpm, req.FloorPrice)
		return false, nil, nil
	}

	res := models.SSPResponse{
		req.Id,
		cpm,
		advId,
		"http://35.194.99.70/notice/"+advId,
	}
	common.Logger.Infof("Bid Req FP %f Res %s", req.FloorPrice, res)

	return true, &res, nil
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

func predict(sReq *models.SSPRequest, advId int) (*models.PredictResp, error) {
	//site,err := strconv.Atoi(sReq.Site)
	data := models.PredictRequest{
		sReq.FloorPrice,
		advId,
	}
	p,err := json.Marshal(data)
	if err != nil {
		common.Logger.Error(err)
		return nil, err
	}

	pReq, err := http.NewRequest(
		"POST",
		"http://35.189.130.108/predict",
		bytes.NewBuffer(p),
	)
	if err != nil {
		common.Logger.Errorf("request creation failed: %s",err.Error())
		return nil, err
	}

	pReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(pReq)
	if err != nil {
		common.Logger.Errorf("request failed: %s",err.Error())
		return nil, err
	}
	defer resp.Body.Close()

	result := models.PredictResp{}
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&result)
	if err != nil {
		common.Logger.Errorf("decoding ctr error %s",err.Error())
		return nil, err
	}

	return &result, nil
}
