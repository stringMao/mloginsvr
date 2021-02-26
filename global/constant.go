package global

//HallToken 大厅服务器令牌
const HallToken = "fskgkshgdsfjgshjk"

//ClientSignKey 客户端协议防篡改签名key
const ClientSignKey = "dgshjhjrerykt"

//HallSignKey 大厅服务器协议防篡改签名key
const HallSignKey = "dgshjhjrerykt"

//ManagerSignKey 管理员协议防篡改签名key
const ManagerSignKey = "dgshjhjrerykt"

//TestSign 测试接口用的默认签名
const TestSign = "test"

//LoginTokenKey 生成用户登入token的key
const LoginTokenKey = "default"

//=账号类型======================================

//AccTypeCommon 普通注册账号
const AccTypeCommon = 0

//AccTypeWechat 微信账号
const AccTypeWechat = 1

//TokenActiveTime 登入token有效时间（秒）
const TokenActiveTime = 60 * 60 * 2
