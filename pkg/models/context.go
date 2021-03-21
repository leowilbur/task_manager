package models

import (
	"github.com/robfig/cron/v3"
)

type Context struct {
	Cron    *cron.Cron
	JobInfo map[cron.EntryID]Task
}
