package mysql

import (
	"github.com/astaxie/beego/config"
)

type mconfig struct {
	DBlog      bool
	DBName     string
	DBHost     string
	DBUsername string
	DBPassword string
	dBMaxIdle  int
	dBMaxConn  int

	ExaminationDBName     string
	ExaminationDBHost     string
	ExaminationDBUsername string
	ExaminationDBPassword string
}

var (
	ModuleConfig mconfig
)

func init() {
	appConf, _ := config.NewConfig("ini", "conf/app.conf")
	appName := appConf.String("appname")
	ModuleConfig = mconfig{}
	ModuleConfig.DBlog, _ = appConf.Bool(appName + "::dBlog")
	ModuleConfig.DBName = appConf.String(appName + "::dBName")
	ModuleConfig.DBHost = appConf.String(appName + "::dBHost")
	ModuleConfig.DBUsername = appConf.String(appName + "::dBUsername")
	ModuleConfig.DBPassword = appConf.String(appName + "::dBPassword")
	ModuleConfig.dBMaxIdle, _ = appConf.Int(appName + "::dBMaxIdle")
	ModuleConfig.dBMaxConn, _ = appConf.Int(appName + "::dBMaxConn")

}
