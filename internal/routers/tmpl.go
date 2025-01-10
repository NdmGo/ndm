package routers

import (
	"fmt"
	"html/template"
	"strings"
	"sync"
	"time"

	"ndm/internal/conf"
)

var (
	funcMap     template.FuncMap
	funcMapOnce sync.Once
)

// FuncMap returns a list of user-defined template functions.
func FuncMap() template.FuncMap {
	// funcMapOnce.Do(func() {

	funcMap = template.FuncMap{
		"title": "测试",
		"AppName": func() string {
			return conf.App.Name
		},
		"AppVer": func() string {
			return conf.App.Version
		},
		"LoadTimes": func(startTime time.Time) string {
			return fmt.Sprint(time.Since(startTime).Nanoseconds()/1e6) + "ms"
		},
		"Join": strings.Join,
		"DateFmtLong": func(t time.Time) string {
			return t.Format(time.RFC1123Z)
		},
		"DateFmtShort": func(t time.Time) string {
			return t.Format("Jan 02, 2006")
		},
	}
	// ---
	// })
	return funcMap
}
