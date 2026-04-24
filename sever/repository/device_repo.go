package repository

import (
	"EmqxBackEnd/models"
	"context"
	"time"

	"gorm.io/gorm"
)

type DeviceRepository struct {
	db *gorm.DB
}

func NewDeviceRepository(db *gorm.DB) *DeviceRepository {
	return &DeviceRepository{db: db}
}

// Create 创建设备：新设备首次接入时注册
func (r *DeviceRepository) Create(ctx context.Context, device *models.Device) error {
	// 使用 WithContext 确保数据库操作受 ctx 控制（如超时取消）
	return r.db.WithContext(ctx).Create(device).Error
}

// Update 更新设备信息：修改位置、IP 或固件版本
func (r *DeviceRepository) Update(ctx context.Context, device *models.Device) error {
	// 使用 Save 方法，它会根据主键 device_id 更新所有非零字段
	// 如果只想更新特定字段，可以使用 .Updates(device)
	return r.db.WithContext(ctx).Save(device).Error
}

// GetByID 获取设备详情：通过设备ID查询
func (r *DeviceRepository) GetByID(ctx context.Context, deviceID string) (*models.Device, error) {
	var device models.Device
	// First 方法会自动处理主键查询
	err := r.db.WithContext(ctx).First(&device, deviceID).Error
	if err != nil {
		return nil, err // 如果找不到，GORM 会返回 gorm.ErrRecordNotFound
	}
	return &device, nil
}

// UpdateHeartbeat 更新心跳状态：设备上报心跳时调用
func (r *DeviceRepository) UpdateHeartbeat(ctx context.Context, deviceID string, status string) error {
	// 使用 UpdateColumns 可以强制更新字段，即使传入的是零值，且不会触发 Hooks
	// 这里我们显式更新 last_heartbeat 为当前时间
	return r.db.WithContext(ctx).Model(&models.Device{}).
		Where("device_id = ?", deviceID).
		UpdateColumns(map[string]interface{}{
			"last_heartbeat": time.Now(),
			"status":         status,
		}).Error
}

// ListOnlineDevices 查询所有在线设备：用于后台监控大屏
func (r *DeviceRepository) ListOnlineDevices(ctx context.Context) ([]*models.Device, error) {
	var devices []*models.Device
	// 查询 status 为 'online' 的记录
	err := r.db.WithContext(ctx).
		Where("status = ?", "online").
		Find(&devices).Error
	return devices, err
}

// Delete 删除设备：设备报废或移除
func (r *DeviceRepository) Delete(ctx context.Context, deviceID string) error {
	// 注意：如果 Device 模型中嵌入了 gorm.Model，Delete 会执行软删除（更新 deleted_at）
	// 如果需要物理删除，请使用 db.Unscoped().Delete(...)
	return r.db.WithContext(ctx).Delete(&models.Device{}, deviceID).Error
}
