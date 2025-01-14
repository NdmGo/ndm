package data

import (
	// "fmt"
	"time"

	"github.com/pkg/errors"
	"gorm.io/gorm"

	"ndm/internal/db"
	"ndm/internal/logs"
	"ndm/internal/model"
	"ndm/internal/utils"
)

func InitAdmin(adminUser string, adminPassword string) {

	admin, err := db.GetAdmin()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			salt := utils.RandString(16)
			admin = &model.User{
				Username:   adminUser,
				Salt:       salt,
				PwdHash:    model.TwoHashPwd(adminPassword, salt),
				Role:       model.ADMIN,
				BasePath:   "/",
				Authn:      "[]",
				Permission: 0xFF, // 0(can see hidden) - 7(can remove)
			}
			admin.PwdTS = time.Now().Unix()
			admin.BasePath = utils.FixAndCleanPath(admin.BasePath)
			if err := db.CreateUser(admin); err != nil {
				logs.Infof("Created the admin user error: %v", err)
			} else {
				logs.Infof("Successfully created the admin user and the initial password is: %s", adminPassword)
			}
		} else {
			logs.Infof("[init user] Failed to get admin user: %v", err)
		}
	}

	guest, err := db.GetGuest()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			salt := utils.RandString(16)
			guest = &model.User{
				Username:   "guest",
				PwdHash:    model.TwoHashPwd("guest", salt),
				Salt:       salt,
				Role:       model.GUEST,
				BasePath:   "/",
				Permission: 0,
				Disabled:   true,
				Authn:      "[]",
			}
			guest.PwdTS = time.Now().Unix()
			if err := db.CreateUser(guest); err != nil {
				logs.Errorf("[init user] Failed to create guest user: %v", err)
			}
		} else {
			logs.Errorf("[init user] Failed to get guest user: %v", err)
		}
	}
}
