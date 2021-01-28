package middle

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mloginsvr/common/log"
	"mloginsvr/global"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type proto struct {
	ServerToken string `json:"servertoken"`
}

//Authentication 服务器身份验证
func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		var bodyBytes []byte
		if c.Request.Body != nil {
			bodyBytes, _ = ioutil.ReadAll(c.Request.Body)
			// 读取后写回
			c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		var data proto
		if err := json.Unmarshal(bodyBytes, &data); err == nil && data.ServerToken == global.HallToken {
			log.Logger.Debug("servertoken 身份验证成功")
			return
		}
		c.Abort()
		c.JSON(http.StatusOK, global.GetResultData(global.CodeSvrIdeErr, "svr 身份验证失败", nil))
		return
	}
}

//AuthCheckSign 自动验签
func AuthCheckSign(key string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if sign := c.GetHeader("sign"); sign != "" {
			var bodyBytes []byte
			if c.Request.Body != nil {
				bodyBytes, _ = ioutil.ReadAll(c.Request.Body)
				// 读取后写回
				c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
			}
			//拼接url最后一项
			urls := strings.Split(c.Request.URL.RequestURI(), "/")
			if urls != nil && len(urls) > 0 {
				bodyBytes = append(bodyBytes, []byte("&"+urls[len(urls)-1])...)
			}

			//拼接加密的key
			bodyBytes = append(bodyBytes, []byte("&"+key)...)

			//32位小写
			mysign := fmt.Sprintf("%x", md5.Sum(bodyBytes))

			if strings.EqualFold(mysign, sign) || sign == global.TestSign {
				//log.Logger.Debug("sign 验证成功")
				return
			}
			log.Logger.Debugf("sign 验签失败，my[%s]==get[%s]", mysign, sign)
			c.Abort()
			c.JSON(http.StatusOK, global.GetResultData(global.CodeSignErr, "sign is err", nil))
			return
		}

		c.Abort()
		c.JSON(http.StatusOK, global.GetResultData(global.CodeSignErr, "no sign", nil))
		return
	}
}
