package local

import (
	// "fmt"

	"ndm/internal/driver"
	"ndm/internal/op"
)

type Addition struct {
	driver.RootPath
	driver.MpId

	ShowHidden     bool   `json:"show_hidden" default:"true" required:"false" help:"show hidden directories and files"`
	MkdirPerm      string `json:"mkdir_perm" default:"777"`
	RecycleBinPath string `json:"recycle_bin_path" default:"delete permanently" help:"path to recycle bin, delete permanently if empty or keep 'delete permanently'"`

	EnableSync bool `json:"enable_sync" default:"false" required:"false" help:"show hidden directories and files"`
	// SyncMpId   int64 `json:"sync_mp_id" default:"0" required:"false"`
}

var config = driver.Config{
	Name:        "local",
	OnlyLocal:   true,
	LocalSort:   true,
	NoCache:     true,
	DefaultRoot: "/",
}

func init() {
	// fmt.Println("init local driver")
	op.RegisterDriver(func() driver.Driver {
		return &Local{}
	})
}
