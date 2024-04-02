package helper

import (
	"bytes"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"encoding/gob"
	"fmt"
	"io"
	"math"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

const (
	DEFALUT_LIMIT = 10
	DEFALUT_PAGE  = 1
)

var uniqueInt64 uint64
var uniqueInt64Chan chan int
var TIMELOCAL *time.Location

func init() {
	uniqueInt64 = 0
	uniqueInt64Chan = make(chan int, 1)

	local, _ := time.LoadLocation("Asia/Chongqing") //服务器设置的时区
	TIMELOCAL = local
}

// 类型转化 string  to int
func StrToInt(str string) int {
	i, _ := strconv.Atoi(str)
	return i
}

// 类型转化 string  to uint64
func StrToUint64(str string) uint64 {
	i, _ := strconv.ParseUint(str, 0, 64)
	return i
}

// 类型转化 string  to float64
func StrToFloat64(str string) float64 {
	f, _ := strconv.ParseFloat(str, 64)
	return f
}

// 类型转化 string  to float32
func StrToFloat32(str string) float32 {
	f, _ := strconv.ParseFloat(str, 64)
	return float32(f)
}

// 类型转化 int to string
func IntToString(i int) string {
	return fmt.Sprintf("%d", i)
}

// 类型转化 int64 to string
func Int64ToString(i int64) string {
	return fmt.Sprintf("%d", i)
}

// 类型转化 uint64 to string
func Uint64ToString(i uint64) string {
	return fmt.Sprintf("%d", i)
}

// 类型转化 uint32 to string
func Uint32ToString(i uint32) string {
	return fmt.Sprintf("%d", i)
}

// 类型转换inerface to string
func InterfaceToString(data interface{}) string {
	return fmt.Sprintf("%s", data)
}

// string数组转化为int数组
func ConvertStringSliceToIntSlice(strSlice []string) ([]int, error) {
	intSlice := make([]int, len(strSlice))
	for i, str := range strSlice {
		num, err := strconv.Atoi(str)
		if err != nil {
			return nil, err
		}
		intSlice[i] = num
	}
	return intSlice, nil
}

// md5
func Md5(str string) string {
	md5Str := fmt.Sprintf("%x", md5.Sum([]byte(str)))
	return md5Str
}

// md5(16位的)
func Md516(str string) string {
	data := md5.Sum([]byte(str))
	md5Str := string(data[0:16])
	return md5Str
}

// get now datatime(Y-m-d H:i:s)
func GetNowDateTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

// get now datatime2(YmdHis)
func GetNowDateTime2() string {
	return time.Now().Format("20060102150405")
}

// get now datatime3(Y-m-d)
func GetNowDateTime3() string {
	return time.Now().Format("2006-01-02")
}

// get now datatime4(H:i:s)
func GetNowDateTime4() string {
	return time.Now().Format("15:04:05")
}

// get yestoday(Y-m-d)
func GetYestoday() string {
	return time.Now().Add(-time.Minute * (time.Duration(24*60) - 1)).Format("2006-01-02")
}

// 获取当前的时间字符串
func GetNowDateTimeDefault() string {
	return time.Now().String()
}

// 获取当前几分钟前的时间(Y-m-d H:i:s)
func GetDateTimeBeforeMinute(num int) string {
	return time.Now().Add(-time.Minute * time.Duration(num)).Format("2006-01-02 15:04:05")
}

// 获取当前几秒钟前的时间(Y-m-d H:i:s)
func GetDateTimeBeforeSecond(num int) string {
	return time.Now().Add(-time.Second * time.Duration(num)).Format("2006-01-02 15:04:05")
}

// 获取当前几分钟后的时间(Y-m-d H:i:s)
func GetDateTimeAfterMinute(num int) string {
	return time.Now().Add(time.Minute * time.Duration(num)).Format("2006-01-02 15:04:05")
}

// 把一个时间字符串转为unix时间戳
func StrToTimeStamp(timeStr string) int64 {
	//	time = "2015-09-14 16:33:00"
	loc, _ := time.LoadLocation("Asia/Shanghai")
	t, _ := time.ParseInLocation("2006-01-02 15:04:05", timeStr, loc)
	return t.Unix()
}

// 把一个unix时间戳转为Y-m-d H:i:s格式的日期
func TimeStampToStr(timeStamp int64) string {
	timeObj := time.Unix(timeStamp-int64(8*60*60), int64(0))
	return timeObj.Format("2006-01-02 15:04:05")
}

// 把一个unix时间戳转为Y-m-d格式的日期
func TimeStampToStr2(timeStamp int64) string {
	timeObj := time.Unix(timeStamp-int64(8*60*60), int64(0))
	return timeObj.Format("2006-01-02")
}

// 切分一个字符串为字符串数组
func Split(str string, flag string) []string {
	return strings.Split(str, flag)
}

// 将给定的日期增加一天并返回
func AddOneDay(date string) (string, error) {
	t, err := time.Parse("2006-01-02", date)
	if err != nil {
		return "", err
	}

	// 将日期增加一天
	t = t.AddDate(0, 0, 1)

	// 返回增加一天后的日期字符串
	return t.Format("2006-01-02"), nil
}

// 合并字符串数组
func JoinString(list []string, flag string) string {
	result := ""
	if len(list) > 0 {
		for _, v := range list {
			result += v + flag
		}
		result = strings.Trim(result, flag)
	}
	return result
}

// 合并int数组
func JoinInt(list []int, flag string) string {
	result := ""
	if len(list) > 0 {
		for _, v := range list {
			result += IntToString(v) + flag
		}
		result = strings.Trim(result, flag)
	}
	return result
}

// 检查一个字符串是否在字符串数组里面
func StringInArray(value string, list []string) bool {
	result := false
	for _, item := range list {
		if value == item {
			result = true
			break
		}
	}
	return result
}

// 检查一个int值是否在int数组里面
func IntInArray(value int, list []int) bool {
	result := false
	for _, item := range list {
		if value == item {
			result = true
			break
		}
	}
	return result
}

// 检查文件或目录是否存在
// 如果由 filename 指定的文件或目录存在则返回 true，否则返回 false
func Exist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

// 检查是否为文件且存在
// 如果由 filename 指定的文件存在则返回 true，否则返回 false
func Exist2(filename string) (exists bool) {
	exists = false
	fileInfo, err := os.Stat(filename)
	if err == nil || os.IsExist(err) {
		if !fileInfo.IsDir() {
			exists = true
		}
	}
	return
}

// sha1
func Sha1(str string) string {
	return fmt.Sprintf("%x", sha1.Sum([]byte(str)))
}

// 获取interface的类型
func GetInterfaceType(i interface{}) string {
	typeObj := reflect.TypeOf(i)
	return typeObj.Kind().String()
}

// 检查interface数据是否为string类型
func CheckInterfaceIsString(i interface{}) bool {
	if i != nil {
		if GetInterfaceType(i) == "string" {
			return true
		}
	}
	return false
}

// 生成Guid字串(32位)
func GetGuid() string {
	b := make([]byte, 48)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return Md5(base64.URLEncoding.EncodeToString(b))
}

// 创建一个20位的唯一数字字符串(flag为两位数的字符数字字符串)
func GetUqunieNumString20(flag string) string {
	nowTime := time.Now().Format("20060102150405")
	uniqueInt64Chan <- 1
	atomic.CompareAndSwapUint64(&uniqueInt64, uint64(9999), uint64(0))
	atomic.AddUint64(&uniqueInt64, 1)
	tmp := Uint64ToString(uniqueInt64)
	<-uniqueInt64Chan
	if len(tmp) == 1 {
		tmp = "000" + tmp
	} else if len(tmp) == 2 {
		tmp = "00" + tmp
	} else if len(tmp) == 3 {
		tmp = "0" + tmp
	}
	id := flag + nowTime + tmp
	return id
}

// copy file
func CopyFile(dstName, srcName string) (written int64, err error) {
	src, err := os.Open(srcName)
	if err != nil {
		return
	}
	defer src.Close()
	dst, err := os.OpenFile(dstName, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return
	}
	defer dst.Close()
	return io.Copy(dst, src)
}

// 深度复制一个对象
func DeepCopy(dst, src interface{}) error {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(src); err != nil {
		return err
	}
	return gob.NewDecoder(bytes.NewBuffer(buf.Bytes())).Decode(dst)
}

// 去重一个int数组
func RmDuplicateInt(list *[]int) []int {
	var x []int = []int{}
	for _, i := range *list {
		if len(x) == 0 {
			x = append(x, i)
		} else {
			for k, v := range x {
				if i == v {
					break
				}
				if k == len(x)-1 {
					x = append(x, i)
				}
			}
		}
	}
	return x
}

// 检查某个汉字有效性
func CheckHanziValid(hanzi string) bool {
	matched, err := regexp.MatchString("^[\u4e00-\u9fa5]{1}$", hanzi)
	if err == nil && matched {
		return true
	}
	return false
}

// get now time
func GetNow() time.Time {
	return time.Now().In(TIMELOCAL)
}

// 字符slice转化成 interface slice
func SliceStrToI(t []string) []interface{} {
	s := make([]interface{}, len(t))
	for i, v := range t {
		s[i] = v
	}
	return s
}

// 字符slice转化成 map slice
func SliceStrToMap(t []string) map[string]struct{} {
	s := make(map[string]struct{}, len(t))
	for _, v := range t {
		s[v] = struct{}{}
	}
	return s
}

// map 转化成 slice
func MapStrToSliceStr(t map[string]struct{}) (s []string) {
	for k, _ := range t {
		s = append(s, k)
	}
	return
}

// map[int] 转化成 int slice
func MapIntToSliceInt(t map[int]struct{}) (s []int) {
	for k, _ := range t {
		s = append(s, k)
	}
	return
}

// 将索引转换为Excel列号（A、B、C等）
func intToExcelColumn(index int) string {
	if index >= 0 && index < 26 {
		return string('A' + index)
	}
	return ""
}

// 清除ID的前缀(学期前缀)
func ClearIdPre(id string) (newId string) {
	newId = id
	if len(id) > 0 {
		matched, err := regexp.MatchString("^[0-9]-[0-9]{4}-[0-9]-", id)
		if err == nil && matched {
			newId = id[9:]
		}
	}
	return
}

func InitPage(paramsLimit int, paramsPage int) (limit int, page int) {
	limit = DEFALUT_LIMIT
	if paramsLimit > 0 {
		limit = paramsLimit
	}
	page = DEFALUT_PAGE
	if paramsPage > 0 {
		page = paramsPage
	}

	return limit, page
}

func PageOffset(limit int, page int) (offset int) {
	if limit > 0 && page > 0 {
		offset = (page - 1) * limit
	}
	return
}

func PageTotal(limit int, page int, count int64) (totalPage int, pageIsEnd int) {
	if count > 0 {
		totalPage = int(math.Ceil(float64(count) / float64(limit)))
	}
	if page >= totalPage {
		pageIsEnd = 1
	}
	return
}

// ExtractStrings 解析包含逗号分隔值的字符串为字符串切片,如:   ["sdasdasd","asdwqewe"]
func ExtractStrings(input string) []string {
	// 使用正则表达式匹配双引号之间的内容
	re := regexp.MustCompile(`"(.*?)"`)
	matches := re.FindAllStringSubmatch(input, -1)

	var extractedStrings []string
	for _, match := range matches {
		extractedStrings = append(extractedStrings, match[1])
	}

	return extractedStrings
}

// MergeSliceToString 将字符串切片合并为一个字符串，并添加方括号
func MergeSliceToString(slice []string) string {
	// 使用逗号连接切片中的值，并添加方括号
	var quotedSlice []string
	for _, s := range slice {
		quotedSlice = append(quotedSlice, `"`+s+`"`)
	}

	resultString := "[" + strings.Join(quotedSlice, ",") + "]"

	return resultString
}

// StrSliceDiff 字符串切片取差集 （ps：a 有，b 没有）
func StrSliceDiff(a, b []string) []string {
	ret := make([]string, 0)

	bSet := StrSliceToSet(b)

	for _, v := range a {
		if _, exist := bSet[v]; !exist {
			ret = append(ret, v)
		}
	}

	return ret
}

// StrSliceToSet 字符串切片转哈希
func StrSliceToSet(a []string) map[string]struct{} {
	ret := make(map[string]struct{})

	for _, v := range a {
		ret[v] = struct{}{}
	}

	return ret
}

// RemoveDuplicateStrArray 去重一个string数组
func RemoveDuplicateStrArray(a []string) []string {
	m := make(map[string]struct{})

	var ret []string
	for _, item := range a {
		if _, exist := m[item]; exist {
			continue

		}

		ret = append(ret, item)
		m[item] = struct{}{}
	}

	return ret
}

// ConvertChineseNumber 将汉字数字转换为对应的数字
func ConvertChineseNumber(s string) (float32, error) {
	chineseNum := map[rune]float32{
		'零': 0, '一': 1, '二': 2, '三': 3, '四': 4,
		'五': 5, '六': 6, '七': 7, '八': 8, '九': 9,
	}

	var result float32
	for _, c := range s {
		if v, ok := chineseNum[c]; ok {
			result = result*10 + v
		} else {
			return 0, fmt.Errorf("无法解析汉字数字：%s", s)
		}
	}

	return result, nil
}
