package db

import (
	"encoding/base64"
	"sync"
	"time"

	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/pkg/errors"

	"ndm/internal/model"
	"ndm/internal/utils"
)

// 性能优化：用户缓存
var (
	userCache     *utils.MemoryCache
	userCacheOnce sync.Once
)

// getUserCache 获取用户缓存实例
func getUserCache() *utils.MemoryCache {
	userCacheOnce.Do(func() {
		userCache = utils.NewMemoryCache(1000, 5*time.Minute)
	})
	return userCache
}

func GetAdmin() (*model.User, error) {
	var adminUser *model.User
	if adminUser == nil {
		user, err := GetUserByRole(model.ADMIN)
		if err != nil {
			return nil, err
		}
		adminUser = user
	}
	return adminUser, nil
}

func GetGuest() (*model.User, error) {
	var guestUser *model.User
	if guestUser == nil {
		user, err := GetUserByRole(model.GUEST)
		if err != nil {
			return nil, err
		}
		guestUser = user
	}
	return guestUser, nil
}

func GetUserByRole(role int64) (*model.User, error) {
	user := model.User{Role: role}
	if err := db.Where(user).Take(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUserByName(username string) (*model.User, error) {
	// 性能优化：先检查缓存
	cache := getUserCache()
	if cachedUser, exists := cache.Get("user:" + username); exists {
		if user, ok := cachedUser.(*model.User); ok {
			return user, nil
		}
	}

	user := model.User{Username: username}
	if err := db.Where(user).First(&user).Error; err != nil {
		return nil, errors.Wrapf(err, "failed find user")
	}

	// 性能优化：更新缓存
	cache.Set("user:"+username, &user)

	return &user, nil
}

func GetUserBySSOID(ssoID string) (*model.User, error) {
	user := model.User{SsoID: ssoID}
	if err := db.Where(user).First(&user).Error; err != nil {
		return nil, errors.Wrapf(err, "The single sign on platform is not bound to any users")
	}
	return &user, nil
}

func GetUserById(id int64) (*model.User, error) {
	var u model.User
	if err := db.First(&u, id).Error; err != nil {
		return nil, errors.Wrapf(err, "failed get old user")
	}
	return &u, nil
}

func CreateUser(u *model.User) error {
	return errors.WithStack(db.Create(u).Error)
}

func UpdateUser(u *model.User) error {
	return errors.WithStack(db.Save(u).Error)
}

func GetUsers(pageIndex, pageSize int) (users []model.User, count int64, err error) {
	userDB := db.Model(&model.User{})
	if err := userDB.Count(&count).Error; err != nil {
		return nil, 0, errors.Wrapf(err, "failed get users count")
	}
	if err := userDB.Order(columnName("id")).Offset((pageIndex - 1) * pageSize).Limit(pageSize).Find(&users).Error; err != nil {
		return nil, 0, errors.Wrapf(err, "failed get find users")
	}
	return users, count, nil
}

func DeleteUserById(id int64) error {
	return errors.WithStack(db.Delete(&model.User{}, id).Error)
}

func UpdateAuthn(userID int64, authn string) error {
	return db.Model(&model.User{ID: userID}).Update("authn", authn).Error
}

func RegisterAuthn(u *model.User, credential *webauthn.Credential) error {
	if u == nil {
		return errors.New("user is nil")
	}
	exists := u.WebAuthnCredentials()
	if credential != nil {
		exists = append(exists, *credential)
	}
	res, err := utils.Json.Marshal(exists)
	if err != nil {
		return err
	}
	return UpdateAuthn(u.ID, string(res))
}

func RemoveAuthn(u *model.User, id string) error {
	exists := u.WebAuthnCredentials()
	for i := 0; i < len(exists); i++ {
		idEncoded := base64.StdEncoding.EncodeToString(exists[i].ID)
		if idEncoded == id {
			exists[len(exists)-1], exists[i] = exists[i], exists[len(exists)-1]
			exists = exists[:len(exists)-1]
			break
		}
	}

	res, err := utils.Json.Marshal(exists)
	if err != nil {
		return err
	}
	return UpdateAuthn(u.ID, string(res))
}
