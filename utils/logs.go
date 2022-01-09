package utils

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"math"
	"os"
	"path"
	"strings"
	"time"
)

type ConsoleFormatter struct {
	logrus.TextFormatter
}

func (c *ConsoleFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	logStr := fmt.Sprintf("%s %s %s:%d %v\n",
		entry.Time.Format("2006/01/02 15:03:04"),
		strings.ToUpper(entry.Level.String()),
		path.Base(entry.Caller.File),
		entry.Caller.Line,
		entry.Message,
	)
	return []byte(logStr), nil
}

// gin access log
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 开始时间
		startTime := time.Now()
		// 处理请求
		c.Next()
		// 捕抓异常
		defer func() {
			if err := recover(); err != nil {
				logrus.Error(err)
			}
		}()
		// 主机名
		hostname, err := os.Hostname()
		if err != nil {
			hostname = "unknow"
		}
		// 结束时间
		endTime := time.Since(startTime)
		// 执行时间
		latencyTime := int(math.Ceil(float64(endTime.Nanoseconds()) / 1000000.0))
		// 请求方式
		reqMethod := c.Request.Method
		// 请求路由
		reqUri := c.Request.RequestURI
		// 状态码
		statusCode := c.Writer.Status()
		// client User Agent
		clientUserAgent := c.Request.UserAgent()
		// referer
		referer := c.Request.Referer()
		// dataLength
		dataLength := c.Writer.Size()
		// 请求IP
		clientIP := c.ClientIP()
		// 请求协议
		proto := c.Request.Proto
		// 日志格式
		fields := logrus.Fields{
			"hostname":     hostname,
			"client_ip":    clientIP,
			"status_code":  statusCode,
			"latency_time": latencyTime,
			"req_method":   reqMethod,
			"req_uri":      reqUri,
			"referer":      referer,
			"data_length":  dataLength,
			"user_agent":   clientUserAgent,
			"req_proto":    proto,
		}
		logrus.Println(Map2StrSliceE(fields))
	}
}
