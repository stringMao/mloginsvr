package alibaba

//阿里云短信服务
import (
	"encoding/json"

	"mloginsvr/common/log"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
)

// todo: 替换为自己账号
const (
	SignName          = "--------------------------"
	REGION_ID         = "cn-shenzhen"
	ACCESS_KEY_ID     = "------------------------------------"
	ACCESS_KEY_SECRET = "------------------------------"
)

//SendSms ..
func SendSms(phone string, code string) (err error) {

	client, err := sdk.NewClientWithAccessKey(REGION_ID, ACCESS_KEY_ID, ACCESS_KEY_SECRET)
	if err != nil {
		log.Logger.Error("ali ecs client failed, err:%s", err.Error())
		return

	}

	request := requests.NewCommonRequest()                           // 构造一个公共请求
	request.Method = "POST"                                          // 设置请求方式
	request.Product = "Ecs"                                          // 指定产品
	request.Scheme = "https"                                         // https | http
	request.Domain = "dysmsapi.aliyuncs.com"                         // 指定域名则不会寻址，如认证方式为 Bearer Token 的服务则需要指定
	request.Version = "20201231"                                     // 指定产品版本
	request.ApiName = "SendSms"                                      // 指定接口名
	request.QueryParams["RegionId"] = "cn-hangzhou"                  // 地区
	request.QueryParams["PhoneNumbers"] = phone                      //手机号
	request.QueryParams["SignName"] = SignName                       //阿里云验证过的项目名 自己设置
	request.QueryParams["TemplateCode"] = "SMS_175543553"            //阿里云的短信模板号 自己设置
	request.QueryParams["TemplateParam"] = "{\"code\":" + code + "}" //短信模板中的验证码内容 自己生成

	response, err := client.ProcessCommonRequest(request)
	if err != nil {
		log.Logger.Error("ali ecs client failed, err:%s", err.Error())
		return
	}
	log.Logger.Info(response.String())
	var message Message //阿里云返回的json信息对应的类
	//记得判断错误信息
	json.Unmarshal(response.GetHttpContentBytes(), &message)
	if message.Message != "OK" {
		//错误处理
		return
	}
	return nil
}

//Message json数据解析
type Message struct {
	Message   string
	RequestId string
	BizId     string
	Code      string
}
