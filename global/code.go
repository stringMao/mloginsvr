package global

//错误码
const (
	//CodeSuccess 成功
	CodeSuccess = 0
	//CodeProtoErr 协议解析错误
	CodeProtoErr = 1

	//=登入相关错误码=====
	//CodeLoginFail 登入失败
	CodeLoginFail = 500

	//CodeCheckTokenFail token验证失败
	CodeCheckTokenFail = 501

	//注册相关==================
	//CodePhoneErr 手机号不正确
	CodePhoneErr = 510
	//CodeSMSOften 申请短息验证码时间间隔未到
	CodeSMSOften = 511
	//CodeSmsVerifyFail 短信验证码错误
	CodeSmsVerifyFail = 512
	//CodeAccOfPhoneRepeat 手机账号重复
	CodeAccOfPhoneRepeat = 513
	//CodeRegiterOfPasswdErr 注册时的密码不合规
	CodeRegiterOfPasswdErr = 514
)
