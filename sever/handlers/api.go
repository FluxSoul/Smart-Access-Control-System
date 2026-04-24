package handlers

import (
	"EmqxBackEnd/models"
	"encoding/json"
	"net/http"
	"time"

	"EmqxBackEnd/task"

	"github.com/gin-gonic/gin"
)

var taskMgr *task.Manager

// SetTaskManager 在main中注入任务管理器实例
func SetTaskManager(tm *task.Manager) {
	taskMgr = tm
}

// GetTasksHandler 获取所有定时任务列表
func GetTasksHandler(c *gin.Context) {
	if taskMgr == nil || taskMgr.Db == nil {
		c.JSON(500, gin.H{"error": "Database not initialized"})
		return
	}
	rows, err := taskMgr.Db.QueryContext(c,
		"SELECT id, task_name, cron_expr, description, status, params, created_at, updated_at FROM cron_tasks ORDER BY id")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var tasks []gin.H
	for rows.Next() {
		var (
			id                   int
			name, cronExpr, desc string
			status               bool
			paramsJSON           []byte
			createdAt, updatedAt time.Time
		)

		err := rows.Scan(&id, &name, &cronExpr, &desc, &status, &paramsJSON, &createdAt, &updatedAt)
		if err != nil {
			continue
		}

		params := make(map[string]interface{})
		json.Unmarshal(paramsJSON, &params)

		tasks = append(tasks, gin.H{
			"id":          id,
			"task_name":   name,
			"cron_expr":   cronExpr,
			"description": desc,
			"status":      status,
			"params":      params,
			"created_at":  createdAt,
			"updated_at":  updatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
		"data": tasks,
	})
}

// UpdateTaskCronHandler 更新任务的Cron表达式
func UpdateTaskCronHandler(c *gin.Context) {
	var req struct {
		CronExpr string `json:"cronExpr" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	name := c.Param("name")
	if err := taskMgr.UpdateTaskCron(name, req.CronExpr); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "任务时间更新成功",
		"task":    name,
		"cron":    req.CronExpr,
	})
}

// UpdateTaskStatusHandler 启用/禁用任务
func UpdateTaskStatusHandler(c *gin.Context) {
	var req struct {
		Status bool `json:"status"` // true=启用, false=禁用
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	name := c.Param("name")

	if req.Status {
		// 启用任务：从数据库重新加载配置
		var cfg struct {
			Name, CronExpr, Desc string
			Status               bool
			ParamsJSON           []byte
		}

		query := "SELECT task_name, cron_expr, description, status, params FROM cron_tasks WHERE task_name = $1"
		err := taskMgr.Db.QueryRowContext(c, query, name).Scan(
			&cfg.Name, &cfg.CronExpr, &cfg.Desc, &cfg.Status, &cfg.ParamsJSON,
		)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "任务不存在"})
			return
		}

		params := make(map[string]interface{})
		json.Unmarshal(cfg.ParamsJSON, &params)

		taskCfg := models.TaskConfig{
			Name:        cfg.Name,
			CronExpr:    cfg.CronExpr,
			Description: cfg.Desc,
			Status:      cfg.Status,
			Params:      params,
		}

		if err := taskMgr.AddTask(taskCfg); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	} else {
		// 禁用任务
		if err := taskMgr.StopTask(name); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "状态更新成功",
	})
}
