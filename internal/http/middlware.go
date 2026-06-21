package http

import (
	"net/http"
	"time"

	"github.com/SayukiDev/VRCLotterySystem/log"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func BodyLimit(maxBytes int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.ContentLength > maxBytes {
			c.AbortWithStatusJSON(http.StatusRequestEntityTooLarge, CommonResp{
				Code: http.StatusRequestEntityTooLarge,
				Msg:  "request entity too large",
			})
			return
		}
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxBytes)
		c.Next()
	}
}

func TokenAuth(token string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Header.Get("Authorization") != token {
			c.AbortWithStatusJSON(http.StatusUnauthorized, CommonResp{
				Code: http.StatusUnauthorized,
				Msg:  "unauthorized",
			})
			return
		}
	}
}

func Logger() gin.HandlerFunc {
	l := log.SubLogger("http")
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		if raw := c.Request.URL.RawQuery; raw != "" {
			path = path + "?" + raw
		}

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()

		fields := []zap.Field{
			zap.Int("status", status),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("ip", c.ClientIP()),
			zap.Duration("latency", latency),
			zap.Int("size", c.Writer.Size()),
		}

		if len(c.Errors) > 0 {
			fields = append(fields, zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()))
		}

		msg := "request"
		switch {
		case status >= 500:
			l.Error(msg, fields...)
		case status >= 400:
			l.Warn(msg, fields...)
		default:
			l.Info(msg, fields...)
		}
	}
}
