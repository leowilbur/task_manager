package models

// Task encapsulates a work item that should go in a work
type Task struct {
	ID   int              `json:"id"`
	Name string           `json:"name"`
	Cron string           `json:"cron"`
	Data chan interface{} `json:"-"`
}
