package global

//错误码
const (
	//CodeSuccess 成功
	CodeSuccess = 0
	//CodeProtoErr 协议解析错误
	CodeProtoErr = 1
	//CodeDBExecErr 数据库执行出错
	CodeDBExecErr = 2
	//CodeSignErr 签名错误
	CodeSignErr = 3
	//CodeAskFail 请求失败，一般与客户端无关
	CodeAskFail = 4
	//CodeSvrIdeErr 服务器身份认证失败
	CodeSvrIdeErr = 100

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
	//CodeNicknameErr 修改的昵称有问题
	CodeNicknameErr = 515

	//实名认证相关==========================
	//CodeIdentityNumErr 身份证号格式错误
	CodeIdentityNumErr = 530

	//==第三方登入============================
	//==微信===
	//CodeWechatLoginCodeErr  由于code错误导致获取access_token 失败
	CodeWechatLoginCodeErr = 550
	//CodeWechatAccessTokenErr 用户的access_token已经无效，请重新code登入
	CodeWechatAccessTokenErr = 551
)
