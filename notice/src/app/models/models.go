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
	BidPrice     float64 `json:"bidprice"`
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
