package global

//RespDataResult .
type RespDataResult struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

//GetResultData ..
func GetResultData(code int, msg string, obj interface{}) RespDataResult {
	r := RespDataResult{code, msg, obj}
	return r
}

//GetResultSucData ..
func GetResultSucData(obj interface{}) RespDataResult {
	r := RespDataResult{CodeSuccess, "", obj}
	return r
}
