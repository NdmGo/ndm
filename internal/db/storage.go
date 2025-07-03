package db

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"ndm/internal/model"
)

func CreateStorage(storage *model.Storage) error {
	return errors.WithStack(db.Create(storage).Error)
}

// UpdateStorage just update storage in database
func UpdateStorage(storage *model.Storage) error {
	return errors.WithStack(db.Save(storage).Error)
}

// DeleteStorageById just delete storage from database by id
func DeleteStorageById(id int64) error {
	return errors.WithStack(db.Delete(&model.Storage{}, id).Error)
}

// TriggerDisabledStorageById just disabled storage from database by id
func TriggerDisabledStorageById(id int64) error {
	s, err := GetStorageById(id)
	if err != nil {
		return err
	}

	disabled := 1
	if s.Disabled {
		disabled = 0
	} else {
		disabled = 1
	}
	return db.Model(&model.Storage{ID: id}).Update("disabled", disabled).Error
}

// GetStorages Get all storages from database order by index
func GetStoragesDriver(page, size int, driverName string) ([]model.Storage, int64, error) {
	storageDB := db.Model(&model.Storage{})
	var count int64
	if err := storageDB.Count(&count).Error; err != nil {
		return nil, 0, errors.Wrapf(err, "failed get storages count")
	}

	if !strings.EqualFold(driverName, "") {
		storageDB = storageDB.Where("driver = ?", driverName)
	}

	var storages []model.Storage
	if err := addStorageOrder(storageDB).Order(columnName("order")).Offset((page - 1) * size).Limit(size).Find(&storages).Error; err != nil {
		return nil, 0, errors.WithStack(err)
	}
	return storages, count, nil
}

// GetNetStorages get network storages from database order by index
func GetNetStorages() ([]model.Storage, error) {
	storageDB := db.Model(&model.Storage{})
	storageDB = storageDB.Where("driver != ?", "local")

	var storages []model.Storage
	if err := addStorageOrder(storageDB).Order(columnName("order")).Find(&storages).Error; err != nil {
		return nil, errors.WithStack(err)
	}
	return storages, nil
}

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

// GetStorageById Get Storage by id, used to update storage usually
func GetStorageById(id int64) (*model.Storage, error) {
	var storage model.Storage
	storage.ID = id
	if err := db.First(&storage).Error; err != nil {
		return nil, errors.WithStack(err)
	}
	return &storage, nil
}

// GetStorageByMountPath Get Storage by mountPath, used to update storage usually
func GetStorageByMountPath(mountPath string) (*model.Storage, error) {
	var storage model.Storage
	if err := db.Where("mount_path = ?", mountPath).First(&storage).Error; err != nil {
		return nil, errors.WithStack(err)
	}
	return &storage, nil
}

func GetEnabledStorages() ([]model.Storage, error) {
	var storages []model.Storage
	err := addStorageOrder(db).Where(fmt.Sprintf("%s = ?", columnName("disabled")), false).Find(&storages).Error
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return storages, nil
}
