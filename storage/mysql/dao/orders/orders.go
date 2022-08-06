package orders

import (
	"time"

	"github.com/hromov/jevelina/storage/mysql"
	"github.com/hromov/jevelina/storage/mysql/dao/models"
)

func CreateTask(c models.CreateLeadReq, lead models.Lead) error {
	task := new(models.Task)
	if c.Description != "" {
		task.Description = c.Description
	} else {
		task.Description = "Call me!"
	}
	t := time.Now()
	task.DeadLine = &t
	task.ParentID = lead.ID
	task.ResponsibleID = lead.ResponsibleID
	if err := mysql.Misc().Create(task).Error; err != nil {
		return err
	}
	return nil
}
