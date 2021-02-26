package global

import (
	"strconv"
)

//Params ..
type Params map[string]string

//map本来已经是引用类型了，所以不需要 *Params

//SetString ..
func (p Params) SetString(k, s string) Params {
	p[k] = s
	return p
}

//GetString ..
func (p Params) GetString(k string) string {
	s, _ := p[k]
	return s
}

//SetInt64 ..
func (p Params) SetInt64(k string, i int64) Params {
	p[k] = strconv.FormatInt(i, 10)
	return p
}

//GetInt64 ..
func (p Params) GetInt64(k string) int64 {
	i, _ := strconv.ParseInt(p.GetString(k), 10, 64)
	return i
}

//GetInt ..
func (p Params) GetInt(k string) int {
	i, _ := strconv.Atoi(p.GetString(k))
	return i
}

//ContainsKey 判断key是否存在
func (p Params) ContainsKey(key string) bool {
	_, ok := p[key]
	return ok
}
