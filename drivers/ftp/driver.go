package ftp

import (
	"context"
	"fmt"
	"os"
	stdpath "path"
	"path/filepath"
	"strings"

	"ndm/internal/driver"
	"ndm/internal/errs"
	"ndm/internal/model"
	"ndm/internal/utils"

	"github.com/jlaffaye/ftp"
)

type FTP struct {
	model.Storage
	Addition
	conn *ftp.ServerConn
}

func (d *FTP) Config() driver.Config {
	return config
}

func (d *FTP) GetAddition() driver.Additional {
	return &d.Addition
}

func (d *FTP) Init(ctx context.Context) error {
	return d.login()
}

func (d *FTP) Drop(ctx context.Context) error {
	if d.conn != nil {
		_ = d.conn.Logout()
	}
	return nil
}

func (d *FTP) List(ctx context.Context, dir model.Obj, args model.ListArgs) ([]model.Obj, error) {
	if err := d.login(); err != nil {
		return nil, err
	}
	entries, err := d.conn.List(encode(dir.GetPath(), d.Encoding))
	if err != nil {
		return nil, err
	}
	res := make([]model.Obj, 0)
	for _, entry := range entries {
		if entry.Name == "." || entry.Name == ".." {
			continue
		}
		f := model.Object{
			Name:     decode(entry.Name, d.Encoding),
			Size:     int64(entry.Size),
			Modified: entry.Time,
			IsFolder: entry.Type == ftp.EntryTypeFolder,
		}
		res = append(res, &f)
	}
	return res, nil
}

func (d *FTP) Link(ctx context.Context, file model.Obj, args model.LinkArgs) (*model.Link, error) {
	if err := d.login(); err != nil {
		return nil, err
	}

	r := NewFileReader(d.conn, encode(file.GetPath(), d.Encoding), file.GetSize())
	link := &model.Link{
		MFile: r,
	}
	return link, nil
}

func (d *FTP) MakeDir(ctx context.Context, parentDir model.Obj, dirName string) error {
	if err := d.login(); err != nil {
		return err
	}
	return d.conn.MakeDir(encode(stdpath.Join(parentDir.GetPath(), dirName), d.Encoding))
}

func (d *FTP) Move(ctx context.Context, srcObj, dstDir model.Obj) error {
	if err := d.login(); err != nil {
		return err
	}
	return d.conn.Rename(
		encode(srcObj.GetPath(), d.Encoding),
		encode(stdpath.Join(dstDir.GetPath(), srcObj.GetName()), d.Encoding),
	)
}

func (d *FTP) Rename(ctx context.Context, srcObj model.Obj, newName string) error {
	if err := d.login(); err != nil {
		return err
	}
	return d.conn.Rename(
		encode(srcObj.GetPath(), d.Encoding),
		encode(stdpath.Join(stdpath.Dir(srcObj.GetPath()), newName), d.Encoding),
	)
}

func (d *FTP) Copy(ctx context.Context, srcObj, dstDir model.Obj) error {
	return errs.NotSupport
}

func (d *FTP) Remove(ctx context.Context, obj model.Obj) error {
	if err := d.login(); err != nil {
		return err
	}
	path := encode(obj.GetPath(), d.Encoding)
	if obj.IsDir() {
		return d.conn.RemoveDirRecur(path)
	} else {
		return d.conn.Delete(path)
	}
}

func (d *FTP) Put(ctx context.Context, dstDir model.Obj, s model.FileStreamer, up driver.UpdateProgress) error {
	if err := d.login(); err != nil {
		return err
	}
	path := stdpath.Join(dstDir.GetPath(), s.GetName())
	return d.conn.Stor(encode(path, d.Encoding), driver.NewLimitedUploadStream(ctx, &driver.ReaderUpdatingProgress{
		Reader:         s,
		UpdateProgress: up,
	}))
}

func (d *FTP) downloadFile(ctx context.Context, key, localfile string) error {
	resp, err := d.conn.Retr(key)
	if err != nil {
		return fmt.Errorf("failed to retrieve remote files: %v", err)
	}
	defer resp.Close()

	dstconn, err := os.Create(localfile)
	if err != nil {
		return err
	}

	if _, err := dstconn.ReadFrom(resp); err != nil {
		return fmt.Errorf("download file failed: %v", err)
	}

	return nil
}

func (d *FTP) BackupFile(ctx context.Context, obj model.Obj, mount_path string) error {
	if !d.EnableBackup {
		return errs.NotEnbleBackup
	}

	if obj.IsDir() {
		return errs.DirNotSupportBackup
	}

	if !utils.IsExist(d.BackupDir) {
		return errs.BackupDirNotExist
	}

	key := obj.GetPath()

	bkdir := strings.TrimRight(d.BackupDir, "/")
	localfile := fmt.Sprintf("%s%s/%s", bkdir, mount_path, key)

	dir := filepath.Dir(localfile)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	if utils.IsExist(localfile) {
		remote_entry, err := d.conn.GetEntry(key)
		if err != nil {
			return fmt.Errorf("ftp failed to obtain remote file information: %v", err)
		}

		localFile, err := os.Stat(localfile)
		if err != nil {
			return fmt.Errorf("local file error: %v", err)
		}

		// compare file sizes
		if remote_entry.Size != uint64(localFile.Size()) {
			return d.downloadFile(ctx, obj.GetPath(), localfile)
		}

		// compare modification time
		remoteModTime := remote_entry.Time
		localModTime := localFile.ModTime()
		if remoteModTime.After(localModTime) {
			return d.downloadFile(ctx, obj.GetPath(), localfile)
		}

		return nil
	}

	return d.downloadFile(ctx, key, localfile)
}

var _ driver.Driver = (*FTP)(nil)
