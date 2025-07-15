package crontab

import (
	"fmt"

	"ndm/internal/db"
)

func Load() {

	tasks, total, err := db.GetTasksList(1, 10)
	if err != nil {
		return
	}

	for _, task := range tasks {
		fmt.Println(task, total)
	}

	fmt.Println("crontab")

	if total < 10 {
		return
	}
}
