package sign

import (
	"sync"
	"time"

	// "ndm/internal/conf"
	// "ndm/internal/setting"
	"ndm/pkg/sign"
)

var once sync.Once
var instance sign.Sign

func Sign(data string) string {
	// expire := setting.GetInt(10, 0)
	expire := 0
	if expire == 0 {
		return NotExpired(data)
	} else {
		return WithDuration(data, time.Duration(expire)*time.Hour)
	}
}

func WithDuration(data string, d time.Duration) string {
	once.Do(Instance)
	return instance.Sign(data, time.Now().Add(d).Unix())
}

func NotExpired(data string) string {
	once.Do(Instance)
	return instance.Sign(data, 0)
}

func Verify(data string, sign string) error {
	once.Do(Instance)
	return instance.Verify(data, sign)
}

func Instance() {
	// setting.GetStr("1231")
	instance = sign.NewHMACSign([]byte("12312"))
}
