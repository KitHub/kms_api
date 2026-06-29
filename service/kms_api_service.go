package service

import (
	"context"
	"sync"

	"github.com/KitHub/kms_api/logic"
	"github.com/KitHub/protocols/kms_api"
)

var (
	kmsapiServiceInstance *KMSAPIService
	kmsapiServiceOnce     sync.Once
)

type KMSAPIService struct {
	kms_api.UnimplementedKMSAPIServer
	demoLogic *logic.DemoLogic
}

// Load implements [kms_api.KMSAPIServer].
func (d *KMSAPIService) Load(context.Context, *kms_api.LoadRequest) (*kms_api.LoadResponse, error) {
	panic("unimplemented")
}

// Store implements [kms_api.KMSAPIServer].
func (d *KMSAPIService) Store(context.Context, *kms_api.StoreRequest) (*kms_api.StoreResponse, error) {
	panic("unimplemented")
}

func NewKMSAPIService(demoLogic *logic.DemoLogic) *KMSAPIService {
	kmsapiServiceOnce.Do(func() {
		kmsapiServiceInstance = &KMSAPIService{
			demoLogic: demoLogic,
		}
	})
	return kmsapiServiceInstance
}
