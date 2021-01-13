package middle

import (
	"mloginsvr/common/log"
	"mloginsvr/global"
	"net/http"

	"github.com/gin-gonic/gin"
)

//Authentication 服务器身份验证
func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		json := make(map[string]interface{}) //
		c.BindJSON(&json)

		if token, ok := json["servertoken"]; ok && token == global.HallToken {
			log.Logger.Debug("servertoken 验证成功")
			return
		}
		c.Abort()
		c.JSON(http.StatusOK, global.GetResultData(global.CodeSvrIdeErr, "svr 身份验证失败", nil))
		return
	}
}
