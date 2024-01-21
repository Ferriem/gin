package utils

import (
	"context"
	"encoding/json"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

type RedisHook struct {
	client *redis.Client
	key    string
}

func NewRedisHook(client *redis.Client, key string) *RedisHook {
	return &RedisHook{
		client: client,
		key:    key,
	}
}

func (hook *RedisHook) Fire(entry *logrus.Entry) error {
	b := entry.Data

	logData, err := json.Marshal(b)
	if err != nil {
		return err
	}

	ctx := context.Background()
	err = hook.client.RPush(ctx, hook.key, logData).Err()
	if err != nil {
		return err
	}
	return nil
}

func (hook *RedisHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func LoggerToRedis(rdb *redis.Client) gin.HandlerFunc {
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	logger.SetFormatter(&logrus.JSONFormatter{})

	redisHook := NewRedisHook(rdb, "logrus")

	logger.AddHook(redisHook)

	return func(c *gin.Context) {
		startTime := time.Now()
		c.Next()
		endTime := time.Now()
		latencyTime := endTime.Sub(startTime)
		reqMethod := c.Request.Method
		reqUri := c.Request.RequestURI

		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		errorMessage := c.Errors.ByType(gin.ErrorTypePrivate).String()

		logger.WithFields(logrus.Fields{
			"status_code":   statusCode,
			"latency_time":  latencyTime,
			"client_ip":     clientIP,
			"error_message": errorMessage,
			"req_method":    reqMethod,
			"req_uri":       reqUri,
		}).Info()
	}
}
