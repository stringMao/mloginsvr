package global

import (
	"encoding/json"
	"mloginsvr/common/log"
	"net/http"
)

type Client struct {
	httpConnectTimeoutMs int // 连接超时时间
	httpReadTimeoutMs    int // 读取超时时间
}

func NewClient() *Client {
	return &Client{
		httpConnectTimeoutMs: 2000,
		httpReadTimeoutMs:    1000,
	}
}

//SetHttpConnectTimeoutMs 设置 连接超时时间
func (c *Client) SetHttpConnectTimeoutMs(ms int) {
	c.httpConnectTimeoutMs = ms
}

//SetHttpReadTimeoutMs  设置 读取超时时间
func (c *Client) SetHttpReadTimeoutMs(ms int) {
	c.httpReadTimeoutMs = ms
}

//GetWithoutCert http GET请求
func (c *Client) GetWithoutCert(url string, params Params, result interface{}) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Logger.Errorf("GetWithoutCert 1 is err url[%s],err[%s] :", url, err.Error())
		return err
	}
	//添加参数
	q := req.URL.Query()
	for k, v := range params {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()

	log.Logger.Debugln(req.URL.String())
	//发送请求
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		log.Logger.Errorf("GetWithoutCert 2 is err url[%s],err[%s] :", url, err.Error())
		return err
	}
	defer response.Body.Close()

	//============

	// body, err := ioutil.ReadAll(response.Body)
	// if err != nil {
	// 	log.Logger.Errorf("GetWithoutCert 3 is err url[%s],err[%s] :", url, err.Error())
	// 	return nil, err
	// }
	// log.Logger.Errorf(string(body))

	//================

	err = json.NewDecoder(response.Body).Decode(result)
	if err != nil {
		log.Logger.Errorf("GetWithoutCert 3 is err url[%s],err[%s] :", url, err.Error())
		return err
	}
	return nil

	// body, err := ioutil.ReadAll(response.Body)
	// if err != nil {
	// 	log.Logger.Errorf("GetWithoutCert 3 is err url[%s],err[%s] :", url, err.Error())
	// 	return nil, err
	// }
	// return XmlToMap(string(body)), nil
}
