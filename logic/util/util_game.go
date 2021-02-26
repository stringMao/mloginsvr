package util

//一些逻辑业务操作整合接口

import (
	"encoding/json"
	"fmt"
	"mloginsvr/common/db"
	"mloginsvr/global"
	"mloginsvr/models"
)

//InsertTokenToRedis 将登入token插入redis
func InsertTokenToRedis(token string, acc models.Account) bool {
	u := global.RedisUserInfo{Gender: 0, CompleteRealname: false, Age: 0}
	u.Token = token
	u.Nickname = acc.Nickname
	u.Accounttype = acc.Accounttype
	//u.Indentity = acc.Identity
	//获得实名信息
	var userinfo models.UserRealInfo
	if userinfo.GetByUserid(acc.Userid) {
		u.CompleteRealname = true
		u.Gender = userinfo.Gender
		u.Age = global.GetCitizenAge([]byte(userinfo.Identity), false)
	}

	var key = fmt.Sprintf("token_%d", acc.Userid)
	if val, err := json.Marshal(u); err == nil {
		db.GetRedis().Do("SET", key, string(val), "EX", global.TokenActiveTime)
		return true
	}
	return false
}

//HallLoadBalanced 大厅渠道分发及负载策略
func HallLoadBalanced(channel int) []global.HallSvrverInfo {
	var result []global.HallSvrverInfo
	//渠道分发
	for _, v := range global.HallList {
		if (channel&v.Channel) > 0 && v.Status == 0 {
			result = append(result, global.HallSvrverInfo{Address: v.Address, Svrname: v.Servername})
		}
	}
	// 负载均衡
	// if len(result) > 1 {

	// }

	return result
}
