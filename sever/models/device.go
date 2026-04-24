package models

import (
	"time"

	"net"
)

// Device 设备模型
// 对应数据库表: devices
type Device struct {
	DeviceID        string     `gorm:"primaryKey;type:varchar(50)" json:"device_id"`     // 设备ID (主键)
	Location        string     `gorm:"type:varchar(150);not null;index" json:"location"` // 安装位置 (带索引)
	IPAddress       net.IP     `gorm:"type:inet" json:"ip_address"`                      // IP地址 (INET 类型)
	Status          string     `gorm:"type:varchar(20);default:'offline'" json:"status"` // 状态
	FirmwareVersion string     `gorm:"type:varchar(20)" json:"firmware_version"`         // 固件版本
	LastHeartbeat   *time.Time `gorm:"type:timestamptz" json:"last_heartbeat"`           // 最后心跳 (可为空)
	CreatedAt       time.Time  `gorm:"type:timestamptz;default:now()" json:"created_at"` // 录入时间
}

func (Device) TableName() string {
	return "iotplus.devices"
}
