package mysql

import (
	"food-service/infrastructure/helper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

var (
	dbList    []*gorm.DB
	dbDefault *gorm.DB
)

var LoggerConfig = logger.New(
	log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
	logger.Config{
		SlowThreshold:             time.Second, // Slow SQL threshold
		LogLevel:                  logger.Info, // Log level
		IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
		ParameterizedQueries:      false,       // Don't include params in the SQL log
		Colorful:                  true,
	})

// mysql 注册信息结构
type DbMysqlRegister struct {
	BaseName   string
	Host       string
	User       string
	Pwd        string
	MaxIdle    int
	MaxConn    int
	SplitTotal int
	Log        bool
}

func init() {
	//初始化db
	dbReigsterInfo := DbMysqlRegister{
		BaseName: ModuleConfig.DBName,
		Host:     ModuleConfig.DBHost,
		User:     ModuleConfig.DBUsername,
		Pwd:      ModuleConfig.DBPassword,
		MaxIdle:  ModuleConfig.dBMaxIdle,
		MaxConn:  ModuleConfig.dBMaxConn,
		Log:      ModuleConfig.DBlog,
	}
	mysqlDbinit(dbReigsterInfo)
}

// 初始化数据库
func mysqlDbinit(reigsterInfo DbMysqlRegister) {
	// 初始化分库的数据库
	dbList = make([]*gorm.DB, reigsterInfo.SplitTotal)
	for i := 0; i < reigsterInfo.SplitTotal; i++ {
		iStr := "_" + helper.IntToString(i)
		db, err := gorm.Open(mysql.Open(reigsterInfo.User+":"+reigsterInfo.Pwd+"@tcp("+reigsterInfo.Host+")/"+reigsterInfo.BaseName+iStr+"?charset=utf8mb4&parseTime=True&loc=Asia%2FShanghai"),
			&gorm.Config{
				Logger: LoggerConfig,
			})
		if err != nil {
			panic("failed to connect database " + reigsterInfo.BaseName + iStr)
		}
		dbList[i] = db
	}
	// 初始化默认的数据库
	db, err := gorm.Open(mysql.Open(reigsterInfo.User+":"+reigsterInfo.Pwd+"@tcp("+reigsterInfo.Host+")/"+reigsterInfo.BaseName+"?charset=utf8mb4&parseTime=True&loc=Asia%2FShanghai"),
		&gorm.Config{
			Logger: LoggerConfig,
		})
	if err != nil {
		panic("failed to connect database " + reigsterInfo.BaseName)
	}
	dbDefault = db
}

// 获取对应的db对象(default)
func GetDbDefault() *gorm.DB {
	return dbDefault
}
