package repository

import (
	"EmqxBackEnd/models"
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type UserPermissionRepository struct {
	db *gorm.DB
}

func NewUserPermissionRepository(db *gorm.DB) *UserPermissionRepository {
	return &UserPermissionRepository{db: db}
}

// Create 新增人员权限：录入新用户
func (r *UserPermissionRepository) Create(ctx context.Context, user *models.UserPermission) error {
	return r.db.WithContext(ctx).Create(user).Error
}

// Update 更新人员信息：修改姓名、有效期或权限列表
func (r *UserPermissionRepository) Update(ctx context.Context, user *models.UserPermission) error {
	// 使用 Save 进行全量更新
	// 注意：如果是部分更新（如只更新状态），建议使用 .Updates(map[string]interface{}{...})
	return r.db.WithContext(ctx).Save(user).Error
}

// GetByID 获取人员详情：通过用户ID查询
func (r *UserPermissionRepository) GetByID(ctx context.Context, userID string) (*models.UserPermission, error) {
	var user models.UserPermission
	err := r.db.WithContext(ctx).First(&user, userID).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Delete 删除人员权限：离职或毕业注销
func (r *UserPermissionRepository) Delete(ctx context.Context, userID string) error {
	return r.db.WithContext(ctx).Delete(&models.UserPermission{}, userID).Error
}

// CheckPermission 权限校验（核心）：检查某用户是否有某设备的通行权限
// 逻辑：利用 PostgreSQL 的 JSONB 操作符 @> (包含)
func (r *UserPermissionRepository) CheckPermission(ctx context.Context, userID string, targetDeviceID string) (bool, error) {
	var count int64

	// 构造查询：
	// 1. 匹配 UserID
	// 2. 匹配 is_active = true
	// 3. 匹配有效期 (valid_start <= now AND valid_end >= now)
	// 4. 匹配 JSONB 字段 allowed_devices 包含 targetDeviceID
	err := r.db.WithContext(ctx).
		Model(&models.UserPermission{}).
		Where("user_id = ? AND is_active = ?", userID, true).
		Where("valid_start <= ? AND valid_end >= ?", time.Now(), time.Now()).
		// 关键点：使用 PostgreSQL 的 jsonb 包含操作符 @>
		// 意思是：查询 allowed_devices 字段包含 ["targetDeviceID"] 的记录
		Where("allowed_devices @> ?", fmt.Sprintf(`["%s"]`, targetDeviceID)).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	// 如果找到记录，说明有权限
	return count > 0, nil
}

// ListExpiringUsers 查询过期或即将过期的用户：用于定时任务清理或提醒
func (r *UserPermissionRepository) ListExpiringUsers(ctx context.Context, days int) ([]*models.UserPermission, error) {
	var users []*models.UserPermission

	// 计算 N 天后的日期
	expiryThreshold := time.Now().AddDate(0, 0, days)

	err := r.db.WithContext(ctx).
		Where("valid_end <= ? AND valid_end >= ?", expiryThreshold, time.Now()). // 有效期在未来 N 天内
		Where("is_active = ?", true).                                            // 必须是激活状态
		Find(&users).Error

	return users, err
}
