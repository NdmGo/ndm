package conf

var (
	BuildTime   string
	BuildCommit string
)

var (
	App struct {
		Version string `ini:"-"`

		Name      string
		BrandName string
		RunUser   string
		RunMode   string
	}

	Http struct {
		Port int64
	}

	// log
	Log struct {
		Format   string
		RootPath string
	}

	// Security settings
	Security struct {
		InstallLock bool
		SecretKey   string
	}
)
