package router

import (
	"EmqxBackEnd/handlers"
	"EmqxBackEnd/middleware"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Setup() *gin.Engine {
	r := gin.Default()

	corsConfig := cors.Config{
		AllowAllOrigins: false, // 生产环境必须设为false
		AllowOrigins: []string{
			"http://localhost:3000",
			"http://localhost:5173",
			"http://172.20.10.5:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,           // 允许携带Cookie
		MaxAge:           12 * time.Hour, // 预检请求缓存时间
	}

	// 应用中间件
	r.Use(cors.New(corsConfig))
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	// r.POST("/empx", handlers.Empx)
	r.POST("/empx/saveMessage", handlers.ReceiveEmpx)
	r.POST("/admin/login", handlers.Login)
	protected := r.Group("")
	protected.Use(middleware.AuthMiddlewareWithCache())
	{
		protected.GET("/admin/getinfo", handlers.GetAdminByAuth)
		protected.GET("/empx/getMessage/:type", handlers.GetMessages)
		protected.GET("/empx/openTheDoor/:nodeId", handlers.OpenTheDoor)
		protected.GET("/empx/closeTheDoor/:nodeId", handlers.CloseTheDoor)
		protected.POST("/admin/register", handlers.Register)
		protected.POST("/admin/saveNode", handlers.SaveNode)
		protected.POST("/admin/changeUserStatus", handlers.ChangeUserStatus)
		protected.GET("/admin/getAllUser", handlers.GetAllUsers)
		protected.GET("/admin/getAllNode", handlers.GetAllNodeByUserId)
	}
	taskGroup := protected.Group("/task")
	{
		taskGroup.GET("", handlers.GetTasksHandler)                      // 获取任务列表
		taskGroup.PUT("/:name/cron", handlers.UpdateTaskCronHandler)     // 更新Cron表达式
		taskGroup.PUT("/:name/status", handlers.UpdateTaskStatusHandler) // 启用/禁用任务
	}
	return r
}
