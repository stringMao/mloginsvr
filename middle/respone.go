package middle

import (
	"bytes"
	"mloginsvr/common/logmanager"

	"github.com/gin-gonic/gin"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	buffer *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.buffer.Write(b)
	return w.ResponseWriter.Write(b)
}

// Respone middleware 输出回复信息日志
func Respone() gin.HandlerFunc {
	return func(c *gin.Context) {
		blw := &bodyLogWriter{c.Writer, &bytes.Buffer{}}
		c.Writer = blw

		c.Next()
		logmanager.Logger.Debug("Response data:", blw.buffer.String())
		blw.Flush()
	}
}
