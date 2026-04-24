package models

import "time"

// LogQuery 用于封装复杂的查询条件
type LogQuery struct {
	UserID     string    // 用户ID
	DeviceID   string    // 设备ID
	Result     string    // 结果：success, failure
	AuthMethod string    // 认证方式：face, rfid
	StartTime  time.Time // 开始时间（零值表示不限制）
	EndTime    time.Time // 结束时间（零值表示不限制）
}
