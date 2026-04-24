package task

import (
	"EmqxBackEnd/models"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
)

type Manager struct {
	cron      *cron.Cron
	Db        *sql.DB
	tasks     map[string]cron.EntryID
	taskCfgs  map[string]models.TaskConfig
	taskFuncs map[string]models.TaskFunc
	mutex     sync.RWMutex
}

// NewManager 创建任务管理器
func NewManager(db *sql.DB) *Manager {
	return &Manager{
		cron:      cron.New(),
		Db:        db,
		tasks:     make(map[string]cron.EntryID),
		taskCfgs:  make(map[string]models.TaskConfig),
		taskFuncs: make(map[string]models.TaskFunc),
	}
}

// RegisterTask 注册任务函数
func (tm *Manager) RegisterTask(desc string, fn models.TaskFunc) {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	// 使用函数名作为键（简单实现）
	var name string
	switch {
	case desc == "温度传感器数据":
		name = "temp_sensor"
	case desc == "获取气体ppm值":
		name = "get_gas_ppm"
	case desc == "获取空气湿度":
		name = "get_gas_moisture"
	case desc == "获取红外传感器数据":
		name = "get_infrared_sensor"
	default:
		name = "unknown_task"
	}

	tm.taskFuncs[name] = fn
	log.Printf("✅ 注册任务函数: %s - %s", name, desc)
}

// LoadTasksFromDB 从PostgreSQL加载任务配置
func (tm *Manager) LoadTasksFromDB() error {
	query := `
		SELECT task_name, cron_expr, description, status, params 
		FROM cron_tasks 
		WHERE status = true
		ORDER BY id
	`

	rows, err := tm.Db.QueryContext(context.Background(), query)
	if err != nil {
		return fmt.Errorf("查询失败: %w", err)
	}
	defer rows.Close()

	loaded := 0
	for rows.Next() {
		var (
			name, cronExpr, desc string
			status               bool
			paramsJSON           []byte
			params               map[string]interface{}
		)

		err := rows.Scan(&name, &cronExpr, &desc, &status, &paramsJSON)
		if err != nil {
			log.Printf("扫描失败: %v", err)
			continue
		}

		// 解析JSON参数
		if len(paramsJSON) > 0 {
			json.Unmarshal(paramsJSON, &params)
		}

		cfg := models.TaskConfig{
			Name:        name,
			CronExpr:    cronExpr,
			Description: desc,
			Status:      status,
			Params:      params,
		}

		if err := tm.AddTask(cfg); err != nil {
			log.Printf("添加任务失败 %s: %v", name, err)
		} else {
			loaded++
		}
	}

	log.Printf("📦 从数据库加载了 %d 个任务", loaded)
	return rows.Err()
}

// AddTask 添加并启动任务
func (tm *Manager) AddTask(cfg models.TaskConfig) error {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	// 检查任务函数是否已注册
	taskFunc, ok := tm.taskFuncs[cfg.Name]
	if !ok {
		return fmt.Errorf("未注册任务函数: %s", cfg.Name)
	}

	// 删除旧任务
	if id, exists := tm.tasks[cfg.Name]; exists {
		tm.cron.Remove(id)
		delete(tm.tasks, cfg.Name)
	}

	// 包装任务函数
	wrappedFunc := func() {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		log.Printf("📋 执行任务: %s (表达式: %s)", cfg.Name, cfg.CronExpr)
		if err := taskFunc(ctx, cfg.Params); err != nil {
			log.Printf("❌ 失败 %s: %v", cfg.Name, err)
		} else {
			log.Printf("✅ 成功 %s", cfg.Name)
		}
	}

	// 添加到cron
	id, err := tm.cron.AddFunc(cfg.CronExpr, wrappedFunc)
	if err != nil {
		return fmt.Errorf("Cron表达式无效: %w", err)
	}

	tm.tasks[cfg.Name] = id
	tm.taskCfgs[cfg.Name] = cfg

	// 根据状态启用/禁用
	if cfg.Status {
		// Entry is already added; nothing more needed here since cron handles execution
		log.Printf("🚀 Scheduled task '%s' with cron expression: %s", cfg.Name, cfg.CronExpr)
	} else {
		// If not active, do not schedule it — remove from cron and internal maps
		tm.cron.Remove(id)
		delete(tm.tasks, cfg.Name)
		log.Printf("⏸️ Task '%s' is inactive and was not scheduled", cfg.Name)
	}

	return nil
}

// UpdateTaskCron 更新任务的Cron表达式
func (tm *Manager) UpdateTaskCron(name, newCron string) error {
	tm.mutex.Lock()

	// 更新数据库
	_, err := tm.Db.ExecContext(context.Background(),
		"UPDATE cron_tasks SET cron_expr = $1 WHERE task_name = $2",
		newCron, name,
	)
	if err != nil {
		return fmt.Errorf("数据库更新失败: %w", err)
	}

	// 从数据库重新加载任务配置
	var cfg models.TaskConfig
	query := "SELECT task_name, cron_expr, description, status, params FROM cron_tasks WHERE task_name = $1"
	var paramsJSON []byte

	err = tm.Db.QueryRowContext(context.Background(), query, name).Scan(
		&cfg.Name, &cfg.CronExpr, &cfg.Description, &cfg.Status, &paramsJSON,
	)
	if err != nil {
		return fmt.Errorf("加载新配置失败: %w", err)
	}
	json.Unmarshal(paramsJSON, &cfg.Params)
	tm.mutex.Unlock()
	return tm.AddTask(cfg)
}

// StopTask 停止任务
func (tm *Manager) StopTask(name string) error {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	// 更新数据库状态
	_, err := tm.Db.ExecContext(context.Background(),
		"UPDATE cron_tasks SET status = false WHERE task_name = $1", name,
	)
	if err != nil {
		return err
	}

	// 停止内存中的任务
	if id, ok := tm.tasks[name]; ok {
		tm.cron.Remove(id) // ✅ 正确方式：使用 Remove 方法移除任务
		delete(tm.tasks, name)
	}

	return nil
}

// StartCron 启动Cron调度器
func (tm *Manager) StartCron() {
	tm.cron.Start()
	log.Println("⏰ Cron调度器已启动")
}

// StopCron 停止Cron调度器
func (tm *Manager) StopCron() {
	tm.cron.Stop()
	log.Println("⏰ Cron调度器已停止")
}
