package infra

import (
	"app/common"
	aerospike"github.com/aerospike/aerospike-client-go"
)

func NewAero() *aerospike.Client {
	hosts := []*aerospike.Host {
		aerospike.NewHost("10.146.0.7", 3000),
		aerospike.NewHost("10.146.0.8", 3000),
	}
	client, err := aerospike.NewClientWithPolicyAndHost(nil, hosts...)

	//client, err := aerospike.NewClient("localhost",3000)

	if err != nil {
		common.Logger.Errorf("error during connecting to aerospike")
	} else {
		common.Logger.Infof("connected to aerospike")
	}

	return client
}
