package database

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DBGorm *gorm.DB

// InitDBGorm 初始化数据库连接
func InitDBGorm() (*gorm.DB, error) {
	// DSN: 用户名:密码@tcp(地址:端口)/库名?参数
	// 注意：这里使用你之前设置的 15432 端口
	dsn := "root:root@tcp(localhost:15432)/iot_business?charset=utf8mb4&parseTime=True&loc=Local"

	var err error
	DBGorm, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("❌ 数据库连接失败: %v", err)
		return nil, err
	}

	fmt.Println("✅ 数据库连接成功")

	// 获取底层 *sql.DB 对象以配置连接池
	sqlDB, err := DBGorm.DB()
	if err != nil {
		log.Fatal(err)
	}

	// 设置连接池参数（生产环境重要）
	sqlDB.SetMaxIdleConns(10)           // 最大空闲连接数
	sqlDB.SetMaxOpenConns(100)          // 最大打开连接数
	sqlDB.SetConnMaxLifetime(time.Hour) // 连接最大存活时间

	return DBGorm, nil
}
