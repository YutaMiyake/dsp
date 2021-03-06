package repository

import (
	"app/common"
	"app/models"
	"github.com/aerospike/aerospike-client-go"
	"strconv"
)

type AeroRepository struct {
	db *aerospike.Client
}

func NewAeroRepository(db *aerospike.Client) AeroRepository {
	return AeroRepository{db: db}
}

func (r *AeroRepository) GetAd(id string) (*models.Advertiser, error) {
	common.Logger.Infof("GetAd %s", id)
	key, err := aerospike.NewKey("test", "advertisers", id)
	if err != nil {
		common.Logger.Error(err)
		return nil, err
	}

	ad := models.Advertiser{}
	err = r.db.GetObject(nil, key, &ad)
	if err != nil {
		common.Logger.Error(err)
		return nil, err
	}

	return &ad, nil
}

func (r *AeroRepository) PutAd(ad *models.Advertiser) error {
	common.Logger.Infof("PutAt %s", ad)
	key, err := aerospike.NewKey("test", "advertisers", ad.Id)
	if err != nil {
		common.Logger.Error(err)
		return err
	}
	policy := aerospike.NewWritePolicy(0, 0)
	err = r.db.PutObject(policy, key, ad)
	if err != nil {
		common.Logger.Error(err)
		return err
	}
	return nil
}

func (r *AeroRepository) GetAllAdInfo() (map[string]*models.Advertiser, error) {
	common.Logger.Infof("GetAllAdInfo")
	adInfo := make(map[string]*models.Advertiser,0)

	for i := 1; i <= 20; i++ {
		var adId string
		ad := new(models.Advertiser)
		if i < 10{
			adId = "adv_0"+strconv.Itoa(i)
		}else{
			adId = "adv_"+strconv.Itoa(i)
		}
		ad,err := r.GetAd(adId)
		if err != nil {
			common.Logger.Error(err)
			return nil,err
		}
		adInfo[adId] = ad
	}
	return adInfo, nil
}

func (r *AeroRepository) GetSpent(id string) (int64, error) {
	//common.Logger.Infof("GetSpent %s", id)
	key, err := aerospike.NewKey("test", "advertisers", id)
	if err != nil {
		common.Logger.Error(err)
		return -1, err
	}

	rec, err := r.db.Get(nil, key, "Spent")
	if err != nil {
		common.Logger.Error(err)
		return -1, err
	}
	spent := int64(rec.Bins["Spent"].(int))

	return spent, nil
}
