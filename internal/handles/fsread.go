package handles

import (
	// "errors"
	"fmt"
	// "net/http"
	// "strconv"
	// "strings"
	"time"

	"ndm/internal/common"
	// "ndm/internal/db"
	"ndm/internal/model"
	// "ndm/internal/op"
	"ndm/internal/fs"
	"ndm/pkg/utils"

	"github.com/gin-gonic/gin"
)

type ListReq struct {
	model.PageReq
	Path     string `json:"path" form:"path"`
	Password string `json:"password" form:"password"`
	Refresh  bool   `json:"refresh"`
}

type DirReq struct {
	Path      string `json:"path" form:"path"`
	Password  string `json:"password" form:"password"`
	ForceRoot bool   `json:"force_root" form:"force_root"`
}

type ObjResp struct {
	Id          string                     `json:"id"`
	Path        string                     `json:"path"`
	Name        string                     `json:"name"`
	Size        int64                      `json:"size"`
	IsDir       bool                       `json:"is_dir"`
	Modified    time.Time                  `json:"modified"`
	Created     time.Time                  `json:"created"`
	Sign        string                     `json:"sign"`
	Thumb       string                     `json:"thumb"`
	Type        int                        `json:"type"`
	HashInfoStr string                     `json:"hashinfo"`
	HashInfo    map[*utils.HashType]string `json:"hash_info"`
}

type FsListResp struct {
	Content  []ObjResp `json:"content"`
	Total    int64     `json:"total"`
	Readme   string    `json:"readme"`
	Header   string    `json:"header"`
	Write    bool      `json:"write"`
	Provider string    `json:"provider"`
}

func FsList(c *gin.Context) {
	var req ListReq
	if err := c.ShouldBind(&req); err != nil {
		common.ErrorResp(c, err, 400)
		return
	}
	req.Validate()
	user := c.MustGet("user").(*model.User)
	reqPath, err := user.JoinPath(req.Path)
	if err != nil {
		common.ErrorResp(c, err, 403)
		return
	}
	// reqPath = "/"
	objs, err := fs.List(c, reqPath, &fs.ListArgs{Refresh: req.Refresh})
	if err != nil {
		common.ErrorResp(c, err, 500)
		return
	}

	fmt.Println(common.ToJson(objs), err)
	fmt.Println(reqPath)

	total, objs := pagination(objs, &req.PageReq)
	provider := "unknown"

	storage, err := fs.GetStorage(reqPath, &fs.GetStoragesArgs{})
	if err == nil {
		provider = storage.GetStorage().Driver
	}

	common.SuccessResp(c, FsListResp{
		Content: toObjsResp(objs, reqPath),
		Total:   int64(total),
		// Readme: getReadme(meta, reqPath),
		// Header: getHeader(meta, reqPath),
		// Write:    user.CanWrite() || common.CanWrite(meta, reqPath),
		Provider: provider,
	})

}

func pagination(objs []model.Obj, req *model.PageReq) (int, []model.Obj) {
	pageIndex, pageSize := req.Page, req.Size
	total := len(objs)
	start := (pageIndex - 1) * pageSize
	if start > total {
		return total, []model.Obj{}
	}
	end := start + pageSize
	if end > total {
		end = total
	}
	return total, objs[start:end]
}

func toObjsResp(objs []model.Obj, parent string) []ObjResp {
	var resp []ObjResp
	for _, obj := range objs {
		// thumb, _ := model.GetThumb(obj)
		resp = append(resp, ObjResp{
			Id:          obj.GetID(),
			Path:        obj.GetPath(),
			Name:        obj.GetName(),
			Size:        obj.GetSize(),
			IsDir:       obj.IsDir(),
			Modified:    obj.ModTime(),
			Created:     obj.CreateTime(),
			HashInfoStr: obj.GetHash().String(),
			HashInfo:    obj.GetHash().Export(),
			// Sign:        common.Sign(obj, parent, encrypt),
			// Thumb: thumb,
			Type: utils.GetObjType(obj.GetName(), obj.IsDir()),
		})
	}
	return resp
}
