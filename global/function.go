package global

import (
	"crypto/md5"
	"fmt"
	"mloginsvr/common/log"
	"reflect"
	"regexp"
	"strconv"
	"time"
)

//StructToRedis ..将结构对象转化成Redis的hash参数
func StructToRedis(key string, src interface{}) []interface{} {
	args := []interface{}{key}
	// 获取结构体实例的反射类型对象
	typeOfSrc := reflect.TypeOf(src)
	valueofSrc := reflect.ValueOf(src)
	// 遍历结构体所有成员
	for i := 0; i < typeOfSrc.NumField(); i++ {
		log.Logger.Debugf("StructToRedis name: %v  value: '%v'", typeOfSrc.Field(i).Name, valueofSrc.Field(i))
		args = append(args, typeOfSrc.Field(i).Name, valueofSrc.Field(i))
	}
	return args
}

//VerifyMobileFormat 手机号合法性检测
func VerifyMobileFormat(mobileNum string) bool {
	regular := "^((13[0-9])|(14[5,7])|(15[0-3,5-9])|(17[0,3,5-8])|(18[0-9])|166|198|199|(147))\\d{8}$"
	reg := regexp.MustCompile(regular)
	return reg.MatchString(mobileNum)
}

//GetLoginToken 计算一个新的token
func GetLoginToken(userid int64) string {
	buf := fmt.Sprintf("userid=%d&time=%d&key=%s", userid, time.Now().Unix(), LoginTokenKey)
	return fmt.Sprintf("%x", md5.Sum([]byte(buf)))
}

//==身份证号相关函数========================================================================================
var weight = [17]int{7, 9, 10, 5, 8, 4, 2, 1, 6, 3, 7, 9, 10, 5, 8, 4, 2}
var valid_value = [11]byte{'1', '0', 'X', '9', '8', '7', '6', '5', '4', '3', '2'}
var valid_province = map[string]string{
	"11": "北京市",
	"12": "天津市",
	"13": "河北省",
	"14": "山西省",
	"15": "内蒙古自治区",
	"21": "辽宁省",
	"22": "吉林省",
	"23": "黑龙江省",
	"31": "上海市",
	"32": "江苏省",
	"33": "浙江省",
	"34": "安徽省",
	"35": "福建省",
	"36": "山西省",
	"37": "山东省",
	"41": "河南省",
	"42": "湖北省",
	"43": "湖南省",
	"44": "广东省",
	"45": "广西壮族自治区",
	"46": "海南省",
	"50": "重庆市",
	"51": "四川省",
	"52": "贵州省",
	"53": "云南省",
	"54": "西藏自治区",
	"61": "陕西省",
	"62": "甘肃省",
	"63": "青海省",
	"64": "宁夏回族自治区",
	"65": "新疆维吾尔自治区",
	"71": "台湾省",
	"81": "香港特别行政区",
	"91": "澳门特别行政区",
}

//IsValidCitizenNo18 效验18位身份证
func IsValidCitizenNo18(citizenNo18 *[]byte) bool {
	nLen := len(*citizenNo18)
	if nLen != 18 {
		return false
	}

	nSum := 0
	for i := 0; i < nLen-1; i++ {
		n, _ := strconv.Atoi(string((*citizenNo18)[i]))
		nSum += n * weight[i]
	}
	mod := nSum % 11
	if valid_value[mod] == (*citizenNo18)[17] {
		return true
	}
	return false
}

//CheckProvinceValid 省份号码效验
func CheckProvinceValid(citizenNo []byte) bool {
	provinceCode := make([]byte, 0)
	provinceCode = append(provinceCode, citizenNo[:2]...)
	provinceStr := string(provinceCode)

	for i := range valid_province {
		if provinceStr == i {
			return true
		}
	}

	return false
}

//IsLeapYear 闰年判断
func IsLeapYear(nYear int) bool {
	if nYear <= 0 {
		return false
	}
	if (nYear%4 == 0 && nYear%100 != 0) || nYear%400 == 0 {
		return true
	}
	return false
}

//CheckBirthdayValid 生日日期格式效验
func CheckBirthdayValid(nYear, nMonth, nDay int) bool {
	if nYear < 1900 || nMonth <= 0 || nMonth > 12 || nDay <= 0 || nDay > 31 {
		return false
	}

	curYear, curMonth, curDay := time.Now().Date()
	if nYear == curYear {
		if nMonth > int(curMonth) {
			return false
		} else if nMonth == int(curMonth) && nDay > curDay {
			return false
		}
	}

	if 2 == nMonth {
		if IsLeapYear(nYear) && nDay > 29 {
			return false
		} else if nDay > 28 {
			return false
		}
	} else if 4 == nMonth || 6 == nMonth || 9 == nMonth || 11 == nMonth {
		if nDay > 30 {
			return false
		}
	}

	return true
}

//IsValidCitizenNo 效验有效地身份证号码
func IsValidCitizenNo(citizenNo *[]byte) bool {
	if !IsValidCitizenNo18(citizenNo) {
		return false
	}

	for i, v := range *citizenNo {
		n, _ := strconv.Atoi(string(v))
		if n >= 0 && n <= 9 {
			continue
		}
		if v == 'X' && i == 16 {
			continue
		}
		return false
	}
	if !CheckProvinceValid(*citizenNo) {
		return false
	}
	nYear, _ := strconv.Atoi(string((*citizenNo)[6:10]))
	nMonth, _ := strconv.Atoi(string((*citizenNo)[10:12]))
	nDay, _ := strconv.Atoi(string((*citizenNo)[12:14]))
	if !CheckBirthdayValid(nYear, nMonth, nDay) {
		return false
	}
	return true

}

//GetCitizenNoInfo 得到身份证号码，生日, 性别, 省份地址信息
// func GetCitizenNoInfo(citizenNo []byte) (err error, birthday time.Time, sex string, address string) {
// 	err = nil
// 	if !IsValidCitizenNo(&citizenNo) {
// 		err = errors.New("不合法的身份证号码。")
// 		return
// 	}
// 	birthday, _ = time.Parse("2006-01-02", string(citizenNo[6:10])+"-"+string(citizenNo[10:12])+"-"+string(citizenNo[12:14]))
// 	genderMask, _ := strconv.Atoi(string(citizenNo[16]))
// 	if genderMask%2 == 0 {
// 		sex = "女"
// 	} else {
// 		sex = "男"
// 	}
// 	address = valid_province[string(citizenNo[:2])]
// 	return nil, birthday, sex, address
// }

//GetCitizenAge 获得身份证年龄  isvalid是否需要先验证身份证号格式合法性
func GetCitizenAge(citizenNo []byte, isvalid bool) int {
	if isvalid && !IsValidCitizenNo(&citizenNo) {
		//"不合法的身份证号码。
		return -1
	}
	birthday, err := time.Parse("2006-01-02", string(citizenNo[6:10])+"-"+string(citizenNo[10:12])+"-"+string(citizenNo[12:14]))
	if err != nil {
		return -1
	}
	cur := time.Now()

	age := cur.Year() - birthday.Year()

	if cur.Month() < birthday.Month() {
		age--
	}
	if cur.Month() == birthday.Month() && cur.Day() < birthday.Day() {
		age--
	}

	return age
}

//GetCitizenGender 根据身份证号获得性别 0男 1女
func GetCitizenGender(citizenNo []byte, isvalid bool) int {
	if isvalid && !IsValidCitizenNo(&citizenNo) {
		//"不合法的身份证号码。
		return 0 //默认男性
	}
	genderMask, _ := strconv.Atoi(string(citizenNo[16]))
	if genderMask%2 == 0 {
		//sex = "女"
		return 1
	}
	//sex = "男"
	return 0
}

//================================================================================================
