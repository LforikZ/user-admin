// @Author Zihao_Li 2024/3/18 21:59:00
package controllers

import "food-service/models"

type UserController struct {
	BaseController
}

type loginParam struct {
	F_code string `form:"code" valid:"Required;"`
}

// GetBanner 获取轮播图数据(临期食品 或者 热销产品)
func (u *UserController) Login() {
	params := &loginParam{}
	if !u.paramsValid(params) {
		u.renderParamsError(nil)
		return
	}

	var user *models.MUser
	openID, err := user.Login(params.F_code)

	if err != nil {
		u.renderUnknownError(err.Error(), nil)
	}

	dates := make(map[string]interface{})
	dates["result"] = openID
	u.renderSuccessJSON(dates)
	return
}
