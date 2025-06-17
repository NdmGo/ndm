package fs

import (
	"context"
	// "fmt"

	"ndm/internal/model"
	"ndm/internal/op"
	// "ndm/pkg/utils"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// List files
func list(ctx context.Context, path string, args *ListArgs) ([]model.Obj, error) {
	user, _ := ctx.Value("user").(*model.User)
	virtualFiles := op.GetStorageVirtualFilesByPath(path)
	storage, actualPath, err := op.GetStorageAndActualPath(path)

	if err != nil && len(virtualFiles) == 0 {
		return nil, errors.WithMessage(err, "failed get storage")
	}

	var _objs []model.Obj
	if storage != nil {
		_objs, err = op.List(ctx, storage, actualPath, model.ListArgs{
			ReqPath: path,
			Refresh: args.Refresh,
		})

		if err != nil {
			if !args.NoLog {
				log.Errorf("fs/list: %+v", err)
			}
			if len(virtualFiles) == 0 {
				return nil, errors.WithMessage(err, "failed get objs")
			}
		}
	}

	om := model.NewObjMerge()
	if whetherHide(user, path) {
		om.InitHideReg("")
	}
	objs := om.Merge(_objs, virtualFiles...)
	return objs, nil
}

func whetherHide(user *model.User, path string) bool {
	// if is admin, don't hide
	if user == nil || user.CanSeeHides() {
		return false
	}
	// if is guest, hide
	return true
}
