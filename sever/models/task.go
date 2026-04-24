package models

import "context"

// TaskFunc 任务函数类型
type TaskFunc func(ctx context.Context, params map[string]interface{}) error

// TaskConfig 任务配置结构
type TaskConfig struct {
	Name        string                 `json:"task_name"`
	CronExpr    string                 `json:"cron_expr"`
	Description string                 `json:"description"`
	Status      bool                   `json:"status"`
	Params      map[string]interface{} `json:"params"`
}

type TaskModel struct {
	ID          int    `db:"id"`
	TaskName    string `db:"task_name"`
	CronExpr    string `db:"cron_expr"`
	Description string `db:"description"`
	Status      bool   `db:"status"`
	Params      string `db:"params"`
}
