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
)
