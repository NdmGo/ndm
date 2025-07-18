package db

import (
	// "fmt"
	// "strings"

	"github.com/pkg/errors"
	"ndm/internal/model"
)

func CreateTasks(task *model.Tasks) error {
	return errors.WithStack(db.Create(task).Error)
}

func DeleteTasksById(id int64) error {
	return errors.WithStack(db.Delete(&model.Tasks{}, id).Error)
}

func UpdateTasks(u *model.Tasks) error {
	return errors.WithStack(db.Save(u).Error)
}

func UpdateTasksCronById(id int64, cron string) error {
	return db.Model(&model.Tasks{ID: id}).Update("cron", cron).Error
}

func GetTasksById(id int64) (*model.Tasks, error) {
	var task model.Tasks
	task.ID = id
	if err := db.First(&task).Error; err != nil {
		return nil, errors.WithStack(err)
	}
	return &task, nil
}

func GetTasksByMpId(mp_id int64) (*model.Tasks, error) {
	var task model.Tasks
	if err := db.First(&task, "mp_id = ?", mp_id).Error; err != nil {
		return nil, errors.WithStack(err)
	}
	return &task, nil
}

func GetTasksList(page, size int) ([]model.Tasks, int64, error) {
	task := db.Model(&model.Tasks{})
	var count int64
	if err := task.Count(&count).Error; err != nil {
		return nil, 0, errors.Wrapf(err, "failed get tasks count")
	}

	var tasks []model.Tasks
	if err := db.Order(columnName("id")).Offset((page - 1) * size).Limit(size).Find(&tasks).Error; err != nil {
		return nil, 0, errors.WithStack(err)
	}
	return tasks, count, nil
}
