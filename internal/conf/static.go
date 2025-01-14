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
		Port     int64
		SafePath string
		ApiPath  string
		Template string
		Debug    bool
	}

	// log
	Log struct {
		Format   string
		RootPath string
	}

	// database
	Database struct {
		Type        string `json:"type" env:"TYPE"`
		Path        string `json:"path" env:"PATH"`
		DSN         string `json:"dsn" env:"DSN"`
		TablePrefix string `json:"table_prefix" env:"TABLE_PREFIX"`
		Hostname    string `json:"hostname" env:"HOST"`
		Hostport    int64  `json:"hostport" env:"PORT"`
		Name        string `json:"name" env:"NAME"`
		User        string `json:"user" env:"USER"`
		Password    string `json:"password" env:"PASS"`
		SSLMode     string `json:"ssl_mode" env:"SSL_MODE"`
	}

	// Security settings
	Security struct {
		InstallLock bool
		SecretKey   string
	}
)
