package crontab

import (
	"fmt"

	"ndm/internal/db"

	"github.com/robfig/cron"
)

var ndmCron = cron.New()

func Load() {

	tasks, total, err := db.GetTasksList(1, 10)
	if err != nil {
		return
	}

	for _, task := range tasks {
		fmt.Println(task, total)

		ndmCron.AddFunc(task.Cron, func() { fmt.Println("Every hour on the half hour") })

	}

	ndmCron.Start()

	fmt.Println("crontab")

	if total < 10 {
		return
	}
}
