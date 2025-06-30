package errs

import (
	"errors"

	pkgerr "github.com/pkg/errors"
)

var (
	ObjectNotFound = errors.New("object not found")
	NotFolder      = errors.New("not a folder")
	NotFile        = errors.New("not a file")

	FailCreateDir = errors.New("failed to create directory")

	//backup
	NotEnbleBackup      = errors.New("backup not supported")
	DirNotSupportBackup = errors.New("directory does not support backup")
	BackupDirNotExist   = errors.New("backup dir does not exist")
	BackupTaskIsRun     = errors.New("backup task is currently in progress")
)

func IsObjectNotFound(err error) bool {
	return errors.Is(pkgerr.Cause(err), ObjectNotFound)
}
