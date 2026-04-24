package models

import (
	"time"

	"gorm.io/datatypes"
)

// UserPermission 用户权限模型
// 对应数据库表: user_permissions
type UserPermission struct {
	UserID             string `gorm:"primaryKey;type:varchar(50)" json:"user_id"` // 用户ID (主键)
	Name               string `gorm:"type:varchar(100);not null" json:"name"`     // 姓名
	FaceFeature        []byte `gorm:"type:bytea" json:"face_feature"`             // 人脸特征值 (二进制)
	FingerprintFeature []byte `gorm:"type:bytea" json:"fingerprint_feature"`      // 指纹特征值 (二进制)

	// 使用 datatypes.JSON 映射 JSONB 字段
	// 它底层是 string，但 GORM 会自动处理 JSON 的编解码
	AllowedDevices datatypes.JSON `gorm:"type:jsonb;default:'[]'" json:"allowed_devices"` // 允许通行的设备列表

	ValidStart time.Time `gorm:"type:timestamptz;not null" json:"valid_start"`     // 有效期开始
	ValidEnd   time.Time `gorm:"type:timestamptz;not null" json:"valid_end"`       // 有效期结束
	IsActive   bool      `gorm:"default:true" json:"is_active"`                    // 账号是否启用
	UpdatedAt  time.Time `gorm:"type:timestamptz;default:now()" json:"updated_at"` // 更新时间
}

// TableName 指定表名
func (UserPermission) TableName() string {
	return "iotplus.user_permissions"
}
