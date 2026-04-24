package handlers

import (
	"EmqxBackEnd/models"
	"EmqxBackEnd/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SaveNode(c *gin.Context) {
	var nodemsg models.Node
	if err := c.ShouldBindJSON(&nodemsg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := service.SaveNode(&nodemsg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "保存成功", "nodeId": nodemsg.ID})
}

func GetAllNodeByUserId(c *gin.Context) {
	token := c.GetHeader("Authorization")
	userId, err := service.GetUserIdByToken(token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	nodes, err := service.GetAllNodeByUserId(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"nodes": nodes})
	}
}
