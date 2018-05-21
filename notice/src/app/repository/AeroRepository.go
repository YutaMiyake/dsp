package repository

import (
	"github.com/aerospike/aerospike-client-go"
	"app/common"
	"app/models"
	"strconv"
)

type AeroRepository struct {
	db *aerospike.Client
}

func NewAeroRepository(db *aerospike.Client ) AeroRepository {
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

func (r *AeroRepository) SpentBudget(adId string, cost int64) (error) {
	common.Logger.Infof("Spent budget %s %d", adId, cost)
	key, err := aerospike.NewKey("test", "advertisers", adId)
	if err != nil {
		common.Logger.Error(err)
		return err
	}
	err = r.db.AddBins(nil, key, aerospike.NewBin("Spent", cost))
	if err != nil {
		common.Logger.Error(err)
		return err
	}
	return nil
}