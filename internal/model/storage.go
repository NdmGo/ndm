package model

import (
	// "fmt"
	"time"

	// "ndm/internal/utils"
)

type Storage struct {
	ID              int64  `json:"id" gorm:"primaryKey"`     // unique key
	MountPath       string `json:"mount_path" gorm:"unique"` // must be standardized
	Order           int    `json:"order"`                    // use to sort
	Driver          string `json:"driver"`                   // driver used
	CacheExpiration int    `json:"cache_expiration"`         // cache expire time
	Status          string `json:"status"`
	Addition        string `json:"addition" gorm:"type:text"` // Additional information, defined in the corresponding driver
	Remark          string `json:"remark"`
	Disabled        bool   `json:"disabled"` // if disabled
	EnableSign      bool   `json:"enable_sign"`
	Sort
	Proxy
	Modified time.Time `json:"modified"`
}

type Sort struct {
	OrderBy        string `json:"order_by"`
	OrderDirection string `json:"order_direction"`
	ExtractFolder  string `json:"extract_folder"`
}

type Proxy struct {
	WebProxy     bool   `json:"web_proxy"`
	WebdavPolicy string `json:"webdav_policy"`
	ProxyRange   bool   `json:"proxy_range"`
	DownProxyUrl string `json:"down_proxy_url"`
}

func (s *Storage) GetStorage() *Storage {
	return s
}

func (s *Storage) SetStorage(storage Storage) {
	*s = storage
}

func (s *Storage) SetStatus(status string) {
	s.Status = status
}

func (s *Storage) GetAdditionMkdirPerm() string {

	// c := s.Addition
	// fmt.Println("GetAdditionMkdirPerm", c)

	// err = utils.Json.UnmarshalFromString(s.Addition, s.GetAddition())
	// fmt.Println(err)
	return ""
}

func (p Proxy) Webdav302() bool {
	return p.WebdavPolicy == "302_redirect"
}

func (p Proxy) WebdavProxy() bool {
	return p.WebdavPolicy == "use_proxy_url"
}

func (p Proxy) WebdavNative() bool {
	return !p.Webdav302() && !p.WebdavProxy()
}
