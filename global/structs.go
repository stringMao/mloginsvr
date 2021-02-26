package global

//RedisUserInfo redis储存的用户信息
type RedisUserInfo struct {
	Token            string
	Nickname         string
	Accounttype      int
	CompleteRealname bool //是否完成实名认证
	Gender           int  //0男 1女
	Age              int  //年龄
}

type RedisWeChatAccessToken struct {
	Userid      int64
	AccessToken string
}

//HallSvrverInfo 大厅服务器信息结构
type HallSvrverInfo struct {
	Address string
	Svrname string
}
