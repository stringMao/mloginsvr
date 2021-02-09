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
