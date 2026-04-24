package models

import (
	"time"
)

// AccessLog 通行日志模型
// 对应数据库表: access_logs
type AccessLog struct {
	LogID    int64  `gorm:"primaryKey;autoIncrement;column:log_id" json:"log_id"` // 日志ID (自增主键)
	UserID   string `gorm:"type:varchar(50);index" json:"user_id"`                // 用户ID (可为空，加索引方便查某人的记录)
	DeviceID string `gorm:"type:varchar(50);not null;index" json:"device_id"`     // 设备ID (加索引方便查某设备的记录)

	AccessTime time.Time `gorm:"type:timestamptz;default:now()" json:"access_time"` // 通行时间
	AuthMethod string    `gorm:"type:varchar(20);not null" json:"auth_method"`      // 通行方式
	Result     string    `gorm:"type:varchar(10);not null" json:"result"`           // 结果

	PhotoURL string `gorm:"type:text" json:"photo_url"` // 抓拍照片 URL
	Reason   string `gorm:"type:text" json:"reason"`    // 失败原因

	CreatedAt time.Time `gorm:"type:timestamptz;default:now()" json:"created_at"` // 入库时间
}

// TableName 指定表名
func (AccessLog) TableName() string {
	return "iotplus.access_logs"
}
