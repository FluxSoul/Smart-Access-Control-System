package handlers

import (
	"EmqxBackEnd/models"
	"EmqxBackEnd/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	var msg models.EmpxAdmin
	if err := c.ShouldBindJSON(&msg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id, isRight := service.CheckLogin(msg.Username, msg.Password)
	if id == -1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户不存在或已被禁止", "message": "登录失败"})
		return
	}
	if isRight {
		token, err := service.GenerateToken(msg.Username)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": "登录失败"})
			return
		}
		if err := service.SaveToken(token, id); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": "登录失败"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "登录成功", "user": gin.H{
			"token": token,
		}})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户名或密码错误", "message": "登录失败"})
	}
}

func Register(c *gin.Context) {
	var msg models.EmpxAdmin
	if err := c.ShouldBindJSON(&msg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if msg.Username == "" || msg.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户名或密码不能为空"})
		return
	}
	createdTime, err := service.CreateAdmin(msg.Username, msg.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": "注册失败"})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "注册成功", "user": gin.H{
			"created_time": createdTime,
		}})
	}
}

func GetAdminByAuth(c *gin.Context) {
	id, exists := c.Get("adminId")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户未登陆"})
		return
	}
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户不存在"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": gin.H{
		"id":       id,
		"username": username,
	}})
}

func ChangeUserStatus(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if !service.IsAdmin(token) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户权限不足"})
		return
	}
	var msg models.EmpxAdmin
	if err := c.ShouldBindJSON(&msg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	if msg.Status != 0 && msg.Status != 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "状态码错误"})
		return
	}
	if err := service.ChangeUserStatus(msg.ID, msg.Status); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "修改成功"})
	}
}

func GetAllUsers(c *gin.Context) {
	// 先对用户进行身份验证
	// 只有status==2时，才可以获取其他用户并修改
	token := c.GetHeader("Authorization")
	if !service.IsAdmin(token) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户权限不足"})
		return
	}
	users, err := service.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"users": users})
	}
}
