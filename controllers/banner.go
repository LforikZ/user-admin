// @Author Zihao_Li 2024/3/15 21:38:00
package controllers

import (
	"food-service/models"
)

type BannerController struct {
	BaseController
}

// GetBanner 获取轮播图数据(临期食品 或者 热销产品)
func (u *BannerController) GetBanner() {
	var banner *models.MBanner
	res, err := banner.GetBanner()

	if err != nil {
		u.renderUnknownError(err.Error(), nil)
	}

	dates := make(map[string]interface{})
	dates["F_list"] = res
	u.renderSuccessJSON(dates)
	return
}
