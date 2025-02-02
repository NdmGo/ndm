package op

import (
	"context"
	"fmt"
	"runtime"
	"time"

	"github.com/pkg/errors"

	log "github.com/sirupsen/logrus"
	"ndm/internal/db"
	"ndm/internal/driver"
	"ndm/internal/model"
	"ndm/pkg/generic_sync"
	"ndm/pkg/utils"
)

var storagesMap generic_sync.MapOf[string, driver.Driver]

func GetAllStorages() []driver.Driver {
	return storagesMap.Values()
}

func HasStorage(mountPath string) bool {
	return storagesMap.Has(utils.FixAndCleanPath(mountPath))
}

func getCurrentGoroutineStack() string {
	buf := make([]byte, 1<<16)
	n := runtime.Stack(buf, false)
	return string(buf[:n])
}

// MustSaveDriverStorage call from specific driver
func MustSaveDriverStorage(driver driver.Driver) {
	err := saveDriverStorage(driver)
	if err != nil {
		log.Errorf("failed save driver storage: %s", err)
	}
}

func saveDriverStorage(driver driver.Driver) error {
	storage := driver.GetStorage()
	addition := driver.GetAddition()
	str, err := utils.Json.MarshalToString(addition)
	if err != nil {
		return errors.Wrap(err, "error while marshal addition")
	}
	storage.Addition = str
	err = db.UpdateStorage(storage)
	if err != nil {
		return errors.WithMessage(err, "failed update storage in database")
	}
	return nil
}

// initStorage initialize the driver and store to storagesMap
func initStorage(ctx context.Context, storage model.Storage, storageDriver driver.Driver) (err error) {
	storageDriver.SetStorage(storage)
	driverStorage := storageDriver.GetStorage()
	defer func() {
		if err := recover(); err != nil {
			errInfo := fmt.Sprintf("[panic] err: %v\nstack: %s\n", err, getCurrentGoroutineStack())
			log.Errorf("panic init storage: %s", errInfo)
			driverStorage.SetStatus(errInfo)
			MustSaveDriverStorage(storageDriver)
			storagesMap.Store(driverStorage.MountPath, storageDriver)
		}
	}()
	// Unmarshal Addition
	err = utils.Json.UnmarshalFromString(driverStorage.Addition, storageDriver.GetAddition())
	if err == nil {
		err = storageDriver.Init(ctx)
	}
	storagesMap.Store(driverStorage.MountPath, storageDriver)
	if err != nil {
		driverStorage.SetStatus(err.Error())
		err = errors.Wrap(err, "failed init storage")
	} else {
		driverStorage.SetStatus(WORK)
	}
	MustSaveDriverStorage(storageDriver)
	return err
}

// CreateStorage Save the storage to database so storage can get an id
// then instantiate corresponding driver and save it in memory
func CreateStorage(ctx context.Context, storage model.Storage) (uint, error) {
	storage.Modified = time.Now()
	storage.MountPath = utils.FixAndCleanPath(storage.MountPath)
	var err error
	// check driver first
	driverName := storage.Driver
	driverNew, err := GetDriver(driverName)
	if err != nil {
		return 0, errors.WithMessage(err, "failed get driver new")
	}
	storageDriver := driverNew()
	// insert storage to database
	err = db.CreateStorage(&storage)
	if err != nil {
		return storage.ID, errors.WithMessage(err, "failed create storage in database")
	}
	// already has an id
	err = initStorage(ctx, storage, storageDriver)
	go callStorageHooks("add", storageDriver)
	if err != nil {
		return storage.ID, errors.Wrap(err, "failed init storage but storage is already created")
	}
	log.Debugf("storage %+v is created", storageDriver)
	return storage.ID, nil
}

// LoadStorage load exist storage in db to memory
func LoadStorage(ctx context.Context, storage model.Storage) error {
	storage.MountPath = utils.FixAndCleanPath(storage.MountPath)
	// check driver first
	driverName := storage.Driver
	driverNew, err := GetDriver(driverName)
	if err != nil {
		return errors.WithMessage(err, "failed get driver new")
	}
	storageDriver := driverNew()

	err = initStorage(ctx, storage, storageDriver)
	go callStorageHooks("add", storageDriver)
	log.Debugf("storage %+v is created", storageDriver)
	return err
}
