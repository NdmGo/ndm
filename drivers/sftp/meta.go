package sftp

import (
	"ndm/internal/driver"
	"ndm/internal/op"
)

type Addition struct {
	Address    string `json:"address" required:"true"`
	Username   string `json:"username" required:"true"`
	PrivateKey string `json:"private_key" type:"text"`
	Password   string `json:"password"`
	Passphrase string `json:"passphrase"`
	driver.RootPath
	IgnoreSymlinkError bool `json:"ignore_symlink_error" default:"false" info:"Ignore symlink error"`

	// backup
	EnableBackup bool   `json:"enable_backup" type:"bool" default:"false" required:"true"`
	BackupDir    string `json:"backup_dir" default:"" required:"false"`
}

var config = driver.Config{
	Name:        "sftp",
	LocalSort:   true,
	OnlyLocal:   true,
	DefaultRoot: "/",
	CheckStatus: true,
}

func init() {
	op.RegisterDriver(func() driver.Driver {
		return &SFTP{}
	})
}
