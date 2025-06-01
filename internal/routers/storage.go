package routers

import (
	"context"
	// "fmt"

	"ndm/internal/db"
	"ndm/internal/model"
	"ndm/internal/op"
	"ndm/pkg/utils"
	// log "github.com/sirupsen/logrus"
)

func LoadStorages() {
	storages, err := db.GetEnabledStorages()
	if err != nil {
		utils.Log.Fatalf("failed get enabled storages: %+v", err)
	}

	go func(storages []model.Storage) {
		for i := range storages {
			err := op.LoadStorage(context.Background(), storages[i])
			if err != nil {
				utils.Log.Errorf("failed get enabled storages: %+v", err)
			} else {
				utils.Log.Infof("success load storage: [%s], driver: [%s], order: [%d]",
					storages[i].MountPath, storages[i].Driver, storages[i].Order)
			}
		}
		// conf.StoragesLoaded = true
	}(storages)
}
