package main

import (
	"EmqxBackEnd/database"
	"EmqxBackEnd/mqtt"
	"EmqxBackEnd/router"
	"context"
	"database/sql"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	mqttBroker := "mqtt://localhost:1883"
	mqttUser := ""
	mqttPass := ""
	if err := mqtt.InitClient(mqttBroker, "cron_task_client", mqttUser, mqttPass); err != nil {
		log.Fatalf("MQTT初始化失败: %v", err)
	}
	defer mqtt.Close()

	db, err := database.Init()
	if err != nil {
		log.Fatal("Failed to connect to DB", err)
		return
	}
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)

	r := router.Setup()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-quit
		log.Println("\n收到退出信号，正在关闭服务...")

		// 5秒超时
		_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// 关闭数据库连接
		_ = db.Close()

		// 断开MQTT连接
		mqtt.Close()

		log.Println("所有资源已释放，服务已停止")
		os.Exit(0)
	}()

	err = r.Run(":8080")
	if err != nil {
		return
	}
}
