// @Author Zihao_Li 2024/3/17 21:46:00
package models

import (
	"encoding/json"
	"food-service/infrastructure/db/mysql"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
)

type MUser struct{}

type User struct {
	F_open_id string `json:"open_id" gorm:"column:F_open_id;primary_key"`
	F_type    int    `json:"type" gorm:"column:F_type"`
}

func (User) TableName() string {
	return "rb_user"
}

const (
	appID     = "wx511e1c56f72e2851"
	appSecret = "e2afbb01100b8d4af81f789e1e679142"
	myOpenID  = "oSLGA6p9a_8EonTJCfYvwz1GatPw"
)

// 登录接口
func (g *MUser) Login(code string) (string, error) {
	var openId string

	resp, err := http.Get("https://api.weixin.qq.com/sns/jscode2session?appid=" + appID + "&secret=" + appSecret + "&js_code=" + code + "&grant_type=authorization_code")
	if err != nil {
		return openId, errors.Wrap(err, "请求微信服务器失败")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return openId, errors.New("解析信息返回值失败")
	}
	var result struct {
		OpenID string `json:"openid"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return openId, errors.New("jason序列化信息失败")
	}
	openId = result.OpenID
	if openId == "" {
		return openId, errors.New("获取用户唯一标识id失败")
	}

	db := mysql.GetDbDefault()

	var num int64
	if err = db.Table(User{}.TableName()).Where("F_open_id = ? ", openId).Count(&num).Error; err != nil {
		return "", errors.Wrap(err, "查询用户信息失败")
	}

	if num == 0 {
		db.Begin()
		if err = db.Create(User{F_open_id: openId, F_type: 2}).Error; err != nil {
			db.Rollback()
			return "", errors.Wrap(err, "创建用户信息失败")
		}
		db.Commit()
	}

	return openId, nil
}

func (g *MUser) VerifyUserID(openId string) (int64, error) {
	db := mysql.GetDbDefault()

	var num int64
	if err := db.Table(User{}.TableName()).Where("F_open_id = ? ", openId).Count(&num).Error; err != nil {
		return 0, errors.Wrap(err, "查询用户信息失败")
	}

	if num > 1 {
		return 0, errors.New("用户id不唯一，请联系后台人员")
	}

	return num, nil
}
