package handlers

import (
	"EmqxBackEnd/jobs"
	"EmqxBackEnd/models"
	"EmqxBackEnd/repository"
	"EmqxBackEnd/service"
	"EmqxBackEnd/state"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

var globalEmqxMsg models.EMQXMessagePublish

func ReceiveEmpx(c *gin.Context) {
	var emqxMsg models.EMQXMessagePublish
	if err := c.ShouldBindJSON(&emqxMsg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// log.Println(emqxMsg)
	var rawMsg models.RawEmpxMessage
	if err := json.Unmarshal([]byte(emqxMsg.Payload), &rawMsg); err != nil {
		log.Println("Error unmarshalling payload:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "payload 解析失败"})
		return
	}
	// log.Println(rawMsg)

	// 字符处理
	realType, err := strconv.Atoi(rawMsg.Type)
	if err != nil {
		log.Println("转换失败:", err)
		return
	}

	realValue := strings.TrimLeft(rawMsg.Value, "0")
	// trimmedValue = "48"

	// 处理全零的特殊情况
	if realValue == "" {
		realValue = "0"
	}

	var msg models.EmpxMessage
	msg.Value = realValue
	msg.NodeID = rawMsg.NodeID
	msg.Type = realType
	msg.TS = time.Now()

	if err := service.ProcessEmpxMessage(&msg); err != nil {
		log.Println("Error saving message:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save"})
		return
	}

	if msg.Type == 4 {
		ppm, err := strconv.Atoi(msg.Value)
		if err != nil {
			log.Println("Error converting ppm to int:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "ppm 解析失败"})
			return
		}
		if ppm >= 2100 {
			// 进入危险值
			// 打开蜂鸣器
			state.SetCache("ppm", 3) // 打开蜂鸣器
			nodeId := msg.NodeID
			messageType := state.GetCache("ppm")
			message := fmt.Sprintf("{\n  \"nodeId\": \"%d\",\n  \"type\": \"%d\"\n}", nodeId, messageType)
			singleParams := map[string]interface{}{
				"topic":    emqxMsg.Topic,
				"message":  message,
				"qos":      emqxMsg.QoS,
				"retained": true,
			}
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if err := jobs.MqttPublishTask(ctx, singleParams); err != nil {
				log.Printf("发布失败[%d]: %v", nodeId, err)
			}
		} else {
			state.SetCache("ppm", 4) // 关闭蜂鸣器
		}
	}
	globalEmqxMsg = emqxMsg
	c.JSON(http.StatusOK, gin.H{"status": "saved"})
}

func GetMessages(c *gin.Context) {
	token := c.GetHeader("Authorization")
	messageType := c.Param("type")
	userId, err := repository.GetUserIdByToken(token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get token"})
		return
	}
	messageTypeId, err := strconv.ParseInt(messageType, 10, 32)
	if err != nil {
		// Handle conversion error
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid type parameter"})
		return
	}

	startTime := "2000-10-11"
	endTime := "2199-01-10"

	if c.Query("startTime") != "" && c.Query("endTime") != "" {
		startTime = c.Query("startTime")
		endTime = c.Query("endTime")
	}

	var messages []models.EmpxMessage

	if messageType == "3" || messageType == "4" {
		var messages3 []models.EmpxMessage
		messages3, err = repository.GetMessagesByDaily(3, userId, startTime, endTime)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get messages"})
			return
		}
		var messages4 []models.EmpxMessage
		messages4, err = repository.GetMessagesByDaily(4, userId, startTime, endTime)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get messages"})
			return
		}
		messages = append(messages3, messages4...)
	} else {
		messages, err = repository.GetMessagesByDaily(int(messageTypeId), userId, startTime, endTime)
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get messages"})
		return
	}
	if len(messages) == 0 {
		c.JSON(http.StatusOK, gin.H{"messages": []models.EmpxMessage{}})
		return
	}
	c.JSON(http.StatusOK, messages)
}

func OpenTheDoor(c *gin.Context) {
	nodeIdS := c.Param("nodeId")
	nodeId, err := strconv.Atoi(nodeIdS)
	if err != nil {
		fmt.Println("转换失败:", err)
		return
	}

	message := fmt.Sprintf("{\n  \"nodeId\": \"%d\",\n  \"type\": \"5\"\n}", nodeId)
	singleParams := map[string]interface{}{
		"topic":    "cmd/esp32",
		"message":  message,
		"qos":      globalEmqxMsg.QoS,
		"retained": true,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := jobs.MqttPublishTask(ctx, singleParams); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get messages"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "开门"})
}

func CloseTheDoor(c *gin.Context) {
	nodeIdS := c.Param("nodeId")
	nodeId, err := strconv.Atoi(nodeIdS)
	if err != nil {
		fmt.Println("转换失败:", err)
		return
	}
	message := fmt.Sprintf("{\n  \"nodeId\": \"%d\",\n  \"type\": \"6\"\n}", nodeId)
	singleParams := map[string]interface{}{
		"topic":    "cmd/esp32",
		"message":  message,
		"qos":      globalEmqxMsg.QoS,
		"retained": true,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := jobs.MqttPublishTask(ctx, singleParams); err != nil {
		log.Printf("发布失败[%d]: %v", nodeId, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get messages"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "关门"})
}
