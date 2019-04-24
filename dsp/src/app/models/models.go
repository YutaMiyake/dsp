package models

type SSPRequest struct {
	Id         string  `json:"id"`
	FloorPrice float64 `json:"floorPrice"`
	Site       string  `json:"site"`
	Date       string  `json:"date"`
	User       string  `json:"user"`
	Test       string  `json:"test"`
}

type SSPResponse struct {
	Id           string  `json:"id"`
	BidPrice     float64 `json:"bidPrice"`
	AdvertiserId string  `json:"advertiserId"`
	Nurl         string  `json:"nurl"`
}

type WinNoticeRequest struct {
	Id      string  `json:"id"`
	Price   float64 `json:"price"`
	IsClick int     `json:"isClick"`
}

type Advertiser struct {
	Id     string `json:"id"`
	Budget int64  `json:"budget"`
	Spent  int64  `json:"spent"`
	Cpc    int64  `json:"cpc"`
}

type PredictRequest struct {
	FloorPrice float64 `json:"floorPrice"`
	AdvertiserId int  `json:"advertiserId"`
}
//Site       int  `json:"site"`
//User       string  `json:"user"`

type PredictResp struct {
	Ctr string `json:"ctr"`
}

// for adtech compe ---------------------------------
type BidRequest struct {
	Id             string  `json:"id"`
	FloorPrice     float64 `json:"floorPrice"`
	DeviceId       string  `json:"deviceId"`
	MediaId        string  `json:"mediaId"`
	Timestamp      int64   `json:"timestamp"`
	OsType         string  `json:"osType"`
	BannerSize     int64   `json:"bannerSize"`
	BannerPosition int64   `json:"bannerPosition"`
	DeviceType     int64   `json:"deviceToken"`
}

type BidResponse struct {
	Id           string  `json:"id"`
	BidPrice     float64 `json:"bid_price"`
	AdvertiserId string  `json:"advertiser_id"`
	Nurl         string  `json:"nurl"`
}