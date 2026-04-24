package repository

import (
	"EmqxBackEnd/models"
	"context"
	"time"

	"gorm.io/gorm"
)

type AccessLogRepository struct {
	db *gorm.DB
}

func NewAccessLogRepository(db *gorm.DB) *AccessLogRepository {
	return &AccessLogRepository{db: db}
}

// Create 新增日志
// 注意：在高并发场景下，如果追求极致性能，可以考虑将日志先写入 Channel，再由协程批量写入 DB
func (r *AccessLogRepository) Create(ctx context.Context, log *models.AccessLog) error {
	// WithContext 确保数据库操作受 ctx 超时控制
	return r.db.WithContext(ctx).Create(log).Error
}

// CreateBatch 批量新增日志
// 适用于设备离线后批量上传记录，或者服务层积攒一批后写入
func (r *AccessLogRepository) CreateBatch(ctx context.Context, logs []*models.AccessLog) error {
	if len(logs) == 0 {
		return nil
	}
	// GORM 的 Create 接收切片时会自动生成批量插入 SQL
	return r.db.WithContext(ctx).Create(logs).Error
}

// List 分页查询日志
// 核心难点：动态构建查询条件
func (r *AccessLogRepository) List(ctx context.Context, query *models.LogQuery, page, pageSize int) ([]*models.AccessLog, int64, error) {
	var logs []*models.AccessLog
	var total int64

	// 1. 构建基础查询对象
	db := r.db.WithContext(ctx).Model(&models.AccessLog{})

	// 2. 动态添加筛选条件
	if query.UserID != "" {
		db = db.Where("user_id = ?", query.UserID)
	}
	if query.DeviceID != "" {
		db = db.Where("device_id = ?", query.DeviceID)
	}
	if query.Result != "" {
		db = db.Where("result = ?", query.Result)
	}
	if query.AuthMethod != "" {
		db = db.Where("auth_method = ?", query.AuthMethod)
	}
	// 时间范围查询：Start Time
	if !query.StartTime.IsZero() {
		db = db.Where("access_time >= ?", query.StartTime)
	}
	// 时间范围查询：End Time
	if !query.EndTime.IsZero() {
		db = db.Where("access_time <= ?", query.EndTime)
	}

	// 3. 查询总数 (用于分页)
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 4. 执行分页查询，按时间倒序排列（最新的在前面）
	offset := (page - 1) * pageSize
	err := db.Order("access_time DESC").Limit(pageSize).Offset(offset).Find(&logs).Error
	if err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

// DeleteBefore 清理旧数据
// 这是一个危险操作，通常由定时任务调用
func (r *AccessLogRepository) DeleteBefore(ctx context.Context, t time.Time) error {
	// 建议使用 Delete 而不是 Unscoped Delete，这样会触发软删除（如果有 gorm.DeletedAt 字段）
	// 如果是物理删除，确保这是你的意图
	return r.db.WithContext(ctx).Where("access_time < ?", t).Delete(&models.AccessLog{}).Error
}
