package conf

const (
	UNKNOWN = iota
	FOLDER
	// OFFICE
	VIDEO
	AUDIO
	TEXT
	IMAGE
)

const (
	TypeString = "string"
	TypeSelect = "select"
	TypeBool   = "bool"
	TypeText   = "text"
	TypeNumber = "number"
)

const (
	// site
	VERSION      = "version"
	SiteTitle    = "site_title"
	Announcement = "announcement"
	AllowIndexed = "allow_indexed"
	AllowMounted = "allow_mounted"
	RobotsTxt    = "robots_txt"

	Logo      = "logo"
	Favicon   = "favicon"
	MainColor = "main_color"

	// preview
	TextTypes          = "text_types"
	AudioTypes         = "audio_types"
	VideoTypes         = "video_types"
	ImageTypes         = "image_types"
	ProxyTypes         = "proxy_types"
	ProxyIgnoreHeaders = "proxy_ignore_headers"
	AudioAutoplay      = "audio_autoplay"
	VideoAutoplay      = "video_autoplay"

	// global
	HideFiles               = "hide_files"
	CustomizeHead           = "customize_head"
	CustomizeBody           = "customize_body"
	LinkExpiration          = "link_expiration"
	SignAll                 = "sign_all"
	PrivacyRegs             = "privacy_regs"
	OcrApi                  = "ocr_api"
	FilenameCharMapping     = "filename_char_mapping"
	ForwardDirectLinkParams = "forward_direct_link_params"
	IgnoreDirectLinkParams  = "ignore_direct_link_params"
	WebauthnLoginEnabled    = "webauthn_login_enabled"

	// index
	SearchIndex     = "search_index"
	AutoUpdateIndex = "auto_update_index"
	IgnorePaths     = "ignore_paths"
	MaxIndexDepth   = "max_index_depth"

	// aria2
	Aria2Uri    = "aria2_uri"
	Aria2Secret = "aria2_secret"

	// transmission
	TransmissionUri      = "transmission_uri"
	TransmissionSeedtime = "transmission_seedtime"

	// single
	Token         = "token"
	IndexProgress = "index_progress"

	// SSO
	SSOClientId          = "sso_client_id"
	SSOClientSecret      = "sso_client_secret"
	SSOLoginEnabled      = "sso_login_enabled"
	SSOLoginPlatform     = "sso_login_platform"
	SSOOIDCUsernameKey   = "sso_oidc_username_key"
	SSOOrganizationName  = "sso_organization_name"
	SSOApplicationName   = "sso_application_name"
	SSOEndpointName      = "sso_endpoint_name"
	SSOJwtPublicKey      = "sso_jwt_public_key"
	SSOExtraScopes       = "sso_extra_scopes"
	SSOAutoRegister      = "sso_auto_register"
	SSODefaultDir        = "sso_default_dir"
	SSODefaultPermission = "sso_default_permission"
	SSOCompatibilityMode = "sso_compatibility_mode"

	// ldap
	LdapLoginEnabled      = "ldap_login_enabled"
	LdapServer            = "ldap_server"
	LdapManagerDN         = "ldap_manager_dn"
	LdapManagerPassword   = "ldap_manager_password"
	LdapUserSearchBase    = "ldap_user_search_base"
	LdapUserSearchFilter  = "ldap_user_search_filter"
	LdapDefaultPermission = "ldap_default_permission"
	LdapDefaultDir        = "ldap_default_dir"
	LdapLoginTips         = "ldap_login_tips"

	// s3
	S3Buckets         = "s3_buckets"
	S3AccessKeyId     = "s3_access_key_id"
	S3SecretAccessKey = "s3_secret_access_key"

	// qbittorrent
	QbittorrentUrl      = "qbittorrent_url"
	QbittorrentSeedtime = "qbittorrent_seedtime"

	// ftp
	FTPPublicHost        = "ftp_public_host"
	FTPPasvPortMap       = "ftp_pasv_port_map"
	FTPProxyUserAgent    = "ftp_proxy_user_agent"
	FTPMandatoryTLS      = "ftp_mandatory_tls"
	FTPImplicitTLS       = "ftp_implicit_tls"
	FTPTLSPrivateKeyPath = "ftp_tls_private_key_path"
	FTPTLSPublicCertPath = "ftp_tls_public_cert_path"
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
		Port           int64
		SafePath       string
		ApiPath        string
		Template       string
		Debug          bool
		TokenExpiresIn int64
		TempDir        string
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
