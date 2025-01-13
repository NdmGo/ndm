package conf

import (
	// "fmt"
	// "log"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"gopkg.in/ini.v1"

	"ndm/internal/utils"
	"ndm/public"
)

// File is the configuration object.
var File *ini.File

func ReadConf() (*ini.File, error) {
	cfg := ini.Empty()
	data, err := public.Conf.ReadFile("conf/app.conf")
	if err != nil {
		return cfg, errors.Wrap(err, "read file 'conf/app.conf'")
	}

	File, err := ini.LoadSources(ini.LoadOptions{
		IgnoreInlineComment: true,
	}, data)

	File.NameMapper = ini.TitleUnderscore
	if err != nil {
		return cfg, errors.Wrap(err, "parse 'conf/app.conf'")
	}

	return cfg, nil
}

func InstallConf(data map[string]string) error {
	File, err := ReadConf()
	if err != nil {
		return err
	}

	err = renderSection(File)
	if err != nil {
		return err
	}

	customConf := filepath.Join(CustomDir(), "conf", "app.conf")

	if !utils.IsExist(filepath.Dir(customConf)) {
		os.MkdirAll(filepath.Dir(customConf), os.ModePerm)
	}

	File.Section("").Key("app_name").SetValue(App.Name)
	File.Section("").Key("brand_name").SetValue(App.BrandName)
	File.Section("").Key("run_user").SetValue(App.RunUser)
	File.Section("").Key("run_mode").SetValue("prod")

	File.Section("log").Key("format").SetValue(Log.Format)
	File.Section("log").Key("root_path").SetValue(Log.RootPath)

	File.Section("http").Key("port").SetValue("5868")
	File.Section("http").Key("save_path").SetValue(Http.SafePath)
	File.Section("http").Key("debug").SetValue("false")

	if strings.EqualFold(data["type"], "mysql") {
		File.Section("database").Key("type").SetValue("mysql")
		File.Section("database").Key("hostname").SetValue(data["hostname"])
		File.Section("database").Key("hostport").SetValue(data["hostport"])
		File.Section("database").Key("name").SetValue(data["dbname"])
		File.Section("database").Key("user").SetValue(data["username"])
		File.Section("database").Key("password").SetValue(data["password"])
		File.Section("database").Key("table_prefix").SetValue(data["table_prefix"])
	} else if strings.EqualFold(data["type"], "sqlite3") {
		File.Section("database").Key("type").SetValue("sqlite3")
		File.Section("database").Key("path").SetValue(data["dbpath"])
	} else {

	}

	File.Section("security").Key("install_lock").SetValue("true")
	File.Section("security").Key("secret_key").SetValue(Security.SecretKey)

	if err := File.SaveTo(customConf); err != nil {
		return err
	}

	return nil
}

func InitConf(customConf string) error {

	data, err := public.Conf.ReadFile("conf/app.conf")
	if err != nil {
		return errors.Wrap(err, "read file 'conf/app.conf'")
	}

	File, err := ini.LoadSources(ini.LoadOptions{
		IgnoreInlineComment: true,
	}, data)

	File.NameMapper = ini.TitleUnderscore

	if err != nil {
		return errors.Wrap(err, "parse 'conf/app.conf'")
	}

	if customConf == "" {
		customConf = filepath.Join(CustomDir(), "conf", "app.conf")
	} else {
		customConf, err = filepath.Abs(customConf)
		if err != nil {
			return errors.Wrap(err, "get absolute path")
		}
	}

	if utils.IsFile(customConf) {
		if err = File.Append(customConf); err != nil {
			return errors.Wrapf(err, "append %q", customConf)
		}
	}
	// else {
	// 	info := fmt.Sprintf("custom config %s not found. Ignore this warning if you're running for the first time", customConf)
	// 	log.Println(info)
	// }

	err = renderSection(File)
	if err != nil {
		return err
	}
	return nil
}

func renderSection(File *ini.File) error {
	if err := File.Section(ini.DefaultSection).MapTo(&App); err != nil {
		return errors.Wrap(err, "mapping default section")
	}

	// ***************************
	// ----- Log settings -----
	// ***************************
	if err := File.Section("log").MapTo(&Log); err != nil {
		return errors.Wrap(err, "mapping [log] section")
	}

	// ***************************
	// ----- Security settings -----
	// ***************************
	if err := File.Section("database").MapTo(&Database); err != nil {
		return errors.Wrap(err, "mapping [database] section")
	}

	// ***************************
	// ----- Http settings -----
	// ***************************
	if err := File.Section("http").MapTo(&Http); err != nil {
		return errors.Wrap(err, "mapping [http] section")
	}

	// ***************************
	// ----- Security settings -----
	// ***************************
	if err := File.Section("security").MapTo(&Security); err != nil {
		return errors.Wrap(err, "mapping [security] section")
	}
	return nil
}
