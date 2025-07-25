package op

import (
	stdpath "path"
	"strings"

	"ndm/internal/driver"
	"ndm/internal/errs"
	"ndm/pkg/utils"

	log "github.com/sirupsen/logrus"
)

// GetStorageAndActualPath Get the corresponding storage and actual path
// for path: remove the mount path prefix and join the actual root folder if exists
func GetStorageAndActualPath(rawPath string) (storage driver.Driver, actualPath string, err error) {
	rawPath = utils.FixAndCleanPath(rawPath)
	storage = GetBalancedStorage(rawPath)
	if storage == nil {
		if rawPath == "/" {
			err = errs.NewErr(errs.StorageNotFound, "please add a storage first")
			return
		}
		err = errs.NewErr(errs.StorageNotFound, "rawPath: %s", rawPath)
		return
	}
	log.Debugln("use storage: ", storage.GetStorage().MountPath)
	mountPath := utils.GetActualMountPath(storage.GetStorage().MountPath)
	actualPath = utils.FixAndCleanPath(strings.TrimPrefix(rawPath, mountPath))
	return
}

// urlTreeSplitLineFormPath split real path and UrlTree definition string from path
func urlTreeSplitLineFormPath(path string) (pp string, file string) {
	// url.PathUnescape will remove //, manually add it back
	path = strings.Replace(path, "https:/", "https://", 1)
	path = strings.Replace(path, "http:/", "http://", 1)
	if strings.Contains(path, ":https:/") || strings.Contains(path, ":http:/") {
		// URL-Tree mode /url_tree_drivr/file_name[:size[:time]]:https://example.com/file
		fPath := strings.SplitN(path, ":", 2)[0]
		pp, _ = stdpath.Split(fPath)
		file = path[len(pp):]
	} else if strings.Contains(path, "/https:/") || strings.Contains(path, "/http:/") {
		// URL-Tree mode /url_tree_drivr/https://example.com/file
		index := strings.Index(path, "/http://")
		if index == -1 {
			index = strings.Index(path, "/https://")
		}
		pp = path[:index]
		file = path[index+1:]
	} else {
		pp, file = stdpath.Split(path)
	}
	if pp == "" {
		pp = "/"
	}
	return
}
