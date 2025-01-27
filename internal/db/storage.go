package db

import (
	// "fmt"

	"github.com/pkg/errors"

	"ndm/internal/model"
)

// GetStorages Get all storages from database order by index
func GetStorages(page, size int) ([]model.Storage, int64, error) {
	storageDB := db.Model(&model.Storage{})
	var count int64
	if err := storageDB.Count(&count).Error; err != nil {
		return nil, 0, errors.Wrapf(err, "failed get storages count")
	}
	var storages []model.Storage
	if err := addStorageOrder(storageDB).Order(columnName("order")).Offset((page - 1) * size).Limit(size).Find(&storages).Error; err != nil {
		return nil, 0, errors.WithStack(err)
	}
	return storages, count, nil
}
