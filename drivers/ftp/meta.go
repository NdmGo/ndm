package ftp

import (
	"ndm/internal/driver"
	"ndm/internal/op"

	"github.com/axgle/mahonia"
)

func encode(str string, encoding string) string {
	if encoding == "" {
		return str
	}
	encoder := mahonia.NewEncoder(encoding)
	return encoder.ConvertString(str)
}

func decode(str string, encoding string) string {
	if encoding == "" {
		return str
	}
	decoder := mahonia.NewDecoder(encoding)
	return decoder.ConvertString(str)
}

type Addition struct {
	driver.RootPath

	Address  string `json:"address" required:"true"`
	Encoding string `json:"encoding" required:"true"`
	Username string `json:"username" required:"true"`
	Password string `json:"password" required:"true"`

	// backup
	EnableBackup bool   `json:"enable_backup" type:"bool" default:"false" required:"false"`
	BackupDir    string `json:"backup_dir" default:"" required:"false"`
}

var config = driver.Config{
	Name:        "ftp",
	LocalSort:   true,
	OnlyLocal:   true,
	DefaultRoot: "/",
}

func init() {
	op.RegisterDriver(func() driver.Driver {
		return &FTP{}
	})
}
