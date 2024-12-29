package conf

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/pkg/errors"
	"gopkg.in/ini.v1"

	"ndm/internal/utils"
	"ndm/public"
)

// File is the configuration object.
var File *ini.File

func InitConf(customConf string) error {

	data, _ := public.Conf.ReadFile("conf/app.conf")

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

	fmt.Println(customConf)

	if utils.IsFile(customConf) {
		if err = File.Append(customConf); err != nil {
			return errors.Wrapf(err, "append %q", customConf)
		}
	} else {
		info := fmt.Sprintf("custom config %s not found. Ignore this warning if you're running for the first time", customConf)
		log.Println(info)
	}

	if err = File.Section(ini.DefaultSection).MapTo(&App); err != nil {
		return errors.Wrap(err, "mapping default section")
	}

	// ***************************
	// ----- Log settings -----
	// ***************************
	if err = File.Section("log").MapTo(&Log); err != nil {
		return errors.Wrap(err, "mapping [log] section")
	}

	fmt.Println(data, err)

	return nil
}
