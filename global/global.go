package global

import (
	"encoding/json"
	"fmt"
	"food-service/global/resp"
	"strings"

	"github.com/astaxie/beego/config"
	_ "github.com/go-sql-driver/mysql"
)

type Mconfig struct {
	NoticeLog                  bool
	Debuglog                   bool
	Errlog                     bool
	PrepareLessonDomain        string
	TeachingAdminApiDomain     string
	TeachingAdminApiToken      string
	ZsResourceDomain           string
	ApiToken                   string
	ConfigMyResponse           map[int]string
	Subject                    map[int]string
	LessonHistoryLimitDefault  int
	TeacherOnlineMQOpen        int
	TeacherOnlineMQChannNum    int
	ReStartTeacherOnlineMQTime int
	TeacherOnlineMQ            string
	TeacherOnlineMqUser        string
	TeacherOnlineMqPwd         string
	TeacherOnlineMqVhost       string
	TeacherOnlineMqDomain      string
	DoCoursewareTypeStr        string
	DoCoursewareTypeList       []string
	ClassDomain                string
	ClassDomainApiToken        string
	LoginServerDomain          string
	SendStatistic              int
	ClientDefaultTimeout       int
	TeacherOnlineBeatHeartTime int
	TsInterationDomain         string
	TsInterationDomainApiToken string
	TransLoop                  bool
	QiniuAccessKey             string
	QiniuSecretKey             string
	QiniuBucket                string
	QiniuZiyuanDomain          string
	EbagIMDomain               string
	EbagIMDomainApiToken       string
	NoticeApiDomain            string
	NoticeApiToken             string
	MsgUrl                     string
	MsgToken                   string
	LiveServerDomain           string
	TrtcLivePullDomain         string
	LiveTrtcAppId              string
	// 定制化功能
	CustomConfigData map[string]CustomConfigItem

	AccountCenterApiDomain string
	WrongBookApiDomain     string
	WrongBookApiToken      string
}

type CustomConfigItem struct {
	F_module string   `json:"module"`
	F_level  int      `json:"level"`
	F_idList []string `json:"idList"`
}

var (
	MyConfig Mconfig
)

func init() {
	appConf, _ := config.NewConfig("ini", "conf/app.conf")
	MyConfig = Mconfig{}

	// 加载定制化功能
	configData := appConf.String("CustomConfig::configdata")
	if len(configData) > 0 {
		configData = strings.Replace(configData, "\\\"", "\"", -1)
		dec := json.NewDecoder(strings.NewReader(configData))
		dec.UseNumber()
		if err := dec.Decode(&MyConfig.CustomConfigData); err != nil {
			fmt.Println("加载定制化功能, 失败:" + err.Error())
		}
	}
	//rep no init
	MyConfig.ConfigMyResponse = resp.GetResponseConfig()
}
