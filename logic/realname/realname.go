package realname

import (
	"mloginsvr/common/log"
	"mloginsvr/global"
	"mloginsvr/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type askRealNameVerify struct {
	Userid   int64  `json:"userid" binding:"required"`
	Name     string `json:"name" binding:"required,min=2,max=12"`
	Identity string `json:"identity" binding:"required"`
	Addr     string `json:"addr" `
}
type replyRealNameVerify struct {
	Gender           int  //0男 1女
	CompleteRealname bool //是否完成实名认证
	Age              int  //年龄
}

//RealNameVerify 实名认证
func RealNameVerify(c *gin.Context) {
	var data askRealNameVerify
	if err := c.ShouldBindJSON(&data); err != nil {
		log.Logger.Errorln("RealNameVerify read data err:", err)
		c.JSON(http.StatusOK, global.GetResultData(global.CodeProtoErr, "协议解析错误", nil))
		return
	}
	log.Logger.Debugf("%+v", data)

	identynum := []byte(data.Identity)
	if !global.IsValidCitizenNo(&identynum) {
		c.JSON(http.StatusOK, global.GetResultData(global.CodeIdentityNumErr, "身份证号错误", nil))
		return
	}
	var realname models.UserRealInfo
	realname.Userid = data.Userid
	realname.Name = data.Name
	realname.Identity = data.Identity
	realname.Gender = global.GetCitizenGender([]byte(data.Identity), false)
	realname.Addr = data.Addr

	if !realname.UpdateOrInsert() {
		c.JSON(http.StatusOK, global.GetResultData(global.CodeDBExecErr, "实名信息db写入失败", nil))
		return
	}

	reply := replyRealNameVerify{CompleteRealname: true, Gender: realname.Gender, Age: global.GetCitizenAge([]byte(realname.Identity), false)}
	c.JSON(http.StatusOK, global.GetResultSucData(reply))
	return
}
