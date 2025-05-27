package s3

import (
	"context"

	"ndm/internal/driver"
	"ndm/internal/model"
)

type Addition struct {
}

type Tpl struct {
	model.Storage
	config driver.Config

	Addition
}

func (d *Tpl) Config() driver.Config {
	return d.config
}

func (d *Tpl) GetAddition() driver.Additional {
	return &d.Addition
}

func (d *Tpl) List(ctx context.Context, dir model.Obj, args model.ListArgs) ([]model.Obj, error) {
	var files []model.Obj
	return files, nil
}

func (d *Tpl) Init(ctx context.Context) error {
	return nil
}

func (d *Tpl) Drop(ctx context.Context) error {
	return nil
}

func (d *Tpl) Link(ctx context.Context, file model.Obj, args model.LinkArgs) (*model.Link, error) {
	var link model.Link
	return &link, nil
}

func (d *Tpl) MakeDir(ctx context.Context, parentDir model.Obj, dirName string) error {
	return nil
}

func (d *Tpl) Move(ctx context.Context, srcObj, dstDir model.Obj) error {
	return nil
}

func (d *Tpl) Rename(ctx context.Context, srcObj model.Obj, newName string) error {
	return nil
}

func (d *Tpl) Copy(ctx context.Context, srcObj, dstDir model.Obj) error {
	return nil
}

func (d *Tpl) Remove(ctx context.Context, obj model.Obj) error {
	return nil
}

var _ driver.Driver = (*Tpl)(nil)
