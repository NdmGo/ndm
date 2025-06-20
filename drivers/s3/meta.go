package s3

import (
	// "fmt"

	"ndm/internal/driver"
	"ndm/internal/op"
)

type Addition struct {
	driver.RootPath
	Bucket                   string `json:"bucket" required:"true"`
	Endpoint                 string `json:"endpoint" required:"true"`
	Region                   string `json:"region"`
	AccessKeyID              string `json:"access_key_id" required:"true"`
	SecretAccessKey          string `json:"secret_access_key" required:"true"`
	SessionToken             string `json:"session_token"`
	CustomHost               string `json:"custom_host"`
	EnableCustomHostPresign  bool   `json:"enable_custom_host_presign"`
	SignURLExpire            int    `json:"sign_url_expire" type:"number" default:"4"`
	Placeholder              string `json:"placeholder"`
	ForcePathStyle           bool   `json:"force_path_style"`
	ListObjectVersion        string `json:"list_object_version" type:"select" options:"v1,v2" default:"v1"`
	RemoveBucket             bool   `json:"remove_bucket" help:"Remove bucket name from path when using custom host."`
	AddFilenameToDisposition bool   `json:"add_filename_to_disposition" help:"Add filename to Content-Disposition header."`

	// backup
	EnableBackup bool   `json:"enable_backup" type:"bool" default:"false" required:"false"`
	BackupDir    string `json:"backup_dir" default:"" required:"false"`
}

func init() {
	op.RegisterDriver(func() driver.Driver {
		return &S3{
			config: driver.Config{
				Name:        "s3",
				DefaultRoot: "/",
				LocalSort:   true,
				CheckStatus: true,
			},
		}
	})
	op.RegisterDriver(func() driver.Driver {
		return &S3{
			config: driver.Config{
				Name:        "doge",
				DefaultRoot: "/",
				LocalSort:   true,
				CheckStatus: true,
			},
		}
	})
}
