package handles

import (
	// "errors"
	"fmt"
	// "net/http"
	// "strconv"
	// stdpath "path"
	"strings"
	"time"

	"ndm/internal/common"
	// "ndm/internal/db"
	"ndm/internal/fs"
	"ndm/internal/model"
	// "ndm/internal/op"
	"ndm/internal/sign"
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

type FsGetReq struct {
	Path     string `json:"path" form:"path"`
	Password string `json:"password" form:"password"`
}

type FsGetResp struct {
	ObjResp
	RawURL   string    `json:"raw_url"`
	Readme   string    `json:"readme"`
	Header   string    `json:"header"`
	Provider string    `json:"provider"`
	Related  []ObjResp `json:"related"`
}

func FsGet(c *gin.Context) {
	var req FsGetReq
	if err := c.ShouldBind(&req); err != nil {
		common.ErrorResp(c, err, 400)
		return
	}

	user := c.MustGet("user").(*model.User)
	reqPath, err := user.JoinPath(req.Path)
	if err != nil {
		common.ErrorResp(c, err, 403)
		return
	}

	obj, err := fs.Get(c, reqPath, &fs.GetArgs{})
	if err != nil {
		common.ErrorResp(c, err, 500)
		return
	}
	var rawURL string

	storage, err := fs.GetStorage(reqPath, &fs.GetStoragesArgs{})
	provider := "unknown"
	if err == nil {
		provider = storage.Config().Name
	}

	fmt.Println("FsGet:obj.IsDir():", obj.IsDir())
	if !obj.IsDir() {
		if err != nil {
			common.ErrorResp(c, err, 500)
			return
		}
		fmt.Println("FsGet:", storage.Config().MustProxy(), storage.GetStorage().WebProxy)
		if storage.Config().MustProxy() || storage.GetStorage().WebProxy {
			query := ""
			// if isEncrypt(meta, reqPath) || setting.GetBool(conf.SignAll) {
			query = "?sign=" + sign.Sign(reqPath)
			// }
			if storage.GetStorage().DownProxyUrl != "" {
				rawURL = fmt.Sprintf("%s%s?sign=%s",
					strings.Split(storage.GetStorage().DownProxyUrl, "\n")[0],
					utils.EncodePath(reqPath, true),
					sign.Sign(reqPath))
			} else {
				rawURL = fmt.Sprintf("%s/p%s%s",
					common.GetApiUrl(c.Request),
					utils.EncodePath(reqPath, true),
					query)
			}
			fmt.Println("....rawURL:", rawURL)
		} else {
			// file have raw url
			if url, ok := model.GetUrl(obj); ok {
				rawURL = url

			} else {
				// if storage is not proxy, use raw url by fs.Link
				link, _, err := fs.Link(c, reqPath, model.LinkArgs{
					IP:       c.ClientIP(),
					Header:   c.Request.Header,
					HttpReq:  c.Request,
					Redirect: true,
				})
				fmt.Println("link:", link)
				if err != nil {
					common.ErrorResp(c, err, 500)
					return
				}
				rawURL = link.URL
			}

			fmt.Println("rawURL:", rawURL)
		}
	}

	// var related []model.Obj
	// parentPath := stdpath.Dir(reqPath)
	// sameLevelFiles, err := fs.List(c, parentPath, &fs.ListArgs{})
	// if err == nil {
	// 	related = filterRelated(sameLevelFiles, obj)
	// }
	// parentMeta, _ := op.GetNearestMeta(parentPath)
	common.SuccessResp(c, FsGetResp{
		ObjResp: ObjResp{
			Id:          obj.GetID(),
			Path:        obj.GetPath(),
			Name:        obj.GetName(),
			Size:        obj.GetSize(),
			IsDir:       obj.IsDir(),
			Modified:    obj.ModTime(),
			Created:     obj.CreateTime(),
			HashInfoStr: obj.GetHash().String(),
			HashInfo:    obj.GetHash().Export(),
			// Sign:        common.Sign(obj, parentPath, isEncrypt(meta, reqPath)),
			Type: utils.GetFileType(obj.GetName()),
		},
		RawURL: rawURL,
		// Readme:   getReadme(meta, reqPath),
		// Header:   getHeader(meta, reqPath),
		Provider: provider,
		// Related:  toObjsResp(related, parentPath, isEncrypt(parentMeta, parentPath)),
	})
}
