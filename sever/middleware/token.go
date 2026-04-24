package middleware

import (
	"EmqxBackEnd/database"
	"EmqxBackEnd/models"

	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type tokenCache struct {
	data map[string]*cacheItem
	mu   sync.RWMutex
}

type cacheItem struct {
	adminID   int
	username  string
	expiresAt time.Time
}

var cache = &tokenCache{
	data: make(map[string]*cacheItem),
}

func init() {
	go func() {
		ticker := time.NewTicker(time.Minute * 5)
		for range ticker.C {
			cache.cleanup()
		}
	}()
}

func (c *tokenCache) cleanup() {
	c.mu.Lock()
	defer c.mu.Unlock()
	now := time.Now()
	for token, item := range c.data {
		if now.After(item.expiresAt) {
			delete(c.data, token)
		}
	}
}

func AuthMiddlewareWithCache() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			token = c.Query("token")
		}

		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token缺失"})
			c.Abort()
			return
		}

		// 查缓存
		cache.mu.RLock()
		item, exist := cache.data[token]
		cache.mu.RUnlock()
		if exist && time.Now().Before(item.expiresAt) {
			// 缓存有效
			c.Set("adminID", item.adminID)
			c.Set("username", item.username)
			c.Next()
			return
		}

		// 查数据库
		admin := models.EmpxAdmin{}

		query := `SELECT id, username, status, token_expires_at 
		          FROM admin WHERE token = $1 AND token_expires_at > CURRENT_TIMESTAMP`

		err := database.DB.QueryRow(query, token).Scan(&admin.ID, &admin.Username, &admin.Status, &admin.TokenExpiresAt)
		if err != nil || admin.Status != 1 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token无效"})
			c.Abort()
			return
		}
		// 写缓存
		cache.mu.Lock()
		cache.data[token] = &cacheItem{
			adminID:   admin.ID,
			username:  admin.Username,
			expiresAt: admin.TokenExpiresAt,
		}
		cache.mu.Unlock()

		c.Set("adminId", admin.ID)
		c.Set("username", admin.Username)
		c.Next()
	}
}
