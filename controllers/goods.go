// @Author Zihao_Li 2024/3/15 21:07:00
package controllers

import (
	"food-service/models"
)

type GoodsController struct {
	BaseController
}
type addGoodParam struct {
	F_title        string  `form:"title" valid:"Required;"`
	F_num          int     `form:"num" valid:"Required;"`
	F_price        float32 `form:"price" valid:"Required;"`
	F_price_module string  `form:"price_module" valid:"Required;"`
	F_url          string  `form:"url" `
	F_url_detail   string  `form:"url_detail" `
	F_tag_id       float32 `form:"tag_id" valid:"Required;"`   // 类别id
	F_tag_name     string  `form:"tag_name" valid:"Required;"` // 类别名
	F_store        string  `form:"store"`                      // 存储方式
	F_area         string  `form:"area"`                       // 商品产地
	F_shelf_life   int     `form:"shelf_life"`                 // 保质期，单位为（天），默认 1 天 到期
}

// AddGood 增加商品
func (u *GoodsController) AddGood() {
	params := &addGoodParam{}
	if !u.paramsValid(params) {
		u.renderParamsError(nil)
		return
	}

	var good *models.MGood
	err := good.AddGood(params.F_title, params.F_num, params.F_price, params.F_price_module, params.F_url, params.F_url_detail, params.F_tag_id, params.F_tag_name, params.F_store, params.F_area, params.F_shelf_life)

	if err != nil {
		u.renderUnknownError(err.Error(), nil)
	}

	u.jsonEchoOk()
	return
}

type updateGoodParam struct {
	F_id           int     `form:"id" valid:"Required;"`
	F_num          int     `form:"num" valid:"Required;"`
	F_title        string  `form:"title" valid:"Required;"`
	F_price        float32 `form:"price" valid:"Required;"`
	F_price_module string  `form:"price_module" valid:"Required;"`
	F_url          string  `form:"url" `
	F_url_detail   string  `form:"url_detail" `
	F_tag_id       float32 `form:"tag_id" valid:"Required;"`   // 类别id
	F_tag_name     string  `form:"tag_name" valid:"Required;"` // 类别名
	F_store        string  `form:"store"`                      // 存储方式
	F_area         string  `form:"area"`                       // 商品产地
	F_shelf_life   int     `form:"shelf_life"`                 // 保质期，单位为（天），默认 1 天 到期
}

// UpdateGood 增加商品
func (u *GoodsController) UpdateGood() {
	params := &updateGoodParam{}
	if !u.paramsValid(params) {
		u.renderParamsError(nil)
		return
	}

	var good *models.MGood
	err := good.UpdateGood(params.F_id, params.F_title, params.F_num, params.F_price, params.F_price_module, params.F_url, params.F_url_detail, params.F_tag_id, params.F_tag_name, params.F_store, params.F_area, params.F_shelf_life)

	if err != nil {
		u.renderUnknownError(err.Error(), nil)
	}

	u.jsonEchoOk()
	return
}

type deleteGoodParam struct {
	F_id   int `form:"id" valid:"Required;"`
	F_flag int `form:"flag" valid:"Required"` // 1表示下架商品，2表示真删除商品
}

// DeleteGood 删除或者下架商品
func (u *GoodsController) DeleteGood() {
	params := &deleteGoodParam{}
	if !u.paramsValid(params) {
		u.renderParamsError(nil)
		return
	}

	var good *models.MGood
	err := good.DeleteGood(params.F_id, params.F_flag)
	if err != nil {
		u.renderUnknownError(err.Error(), nil)
	}

	u.jsonEchoOk()
	return
}

type getGoodListParam struct {
	F_Page  int `form:"page"`
	F_Limit int `form:"limit"`
}

// GetGoodList 获取商品列表
func (u *GoodsController) GetGoodList() {
	params := &getGoodListParam{}
	if !u.paramsValid(params) {
		u.renderParamsError(nil)
		return
	}

	var banner *models.MGood
	res, err := banner.GetGoodList(params.F_Page, params.F_Limit)

	if err != nil {
		u.renderUnknownError(err.Error(), nil)
	}

	dates := make(map[string]interface{})
	dates["F_list"] = res
	u.renderSuccessJSON(dates)
	return
}

type getFuzzSearchGoodParam struct {
	F_Search string `form:"search"`
}

// GetFuzzSearchGoodList 获取模糊搜索商品列表
func (u *GoodsController) GetFuzzSearchGoodList() {
	params := &getFuzzSearchGoodParam{}
	if !u.paramsValid(params) {
		u.renderParamsError(nil)
		return
	}

	var banner *models.MGood
	res, err := banner.GetFuzzSearchGoodList(params.F_Search)

	if err != nil {
		u.renderUnknownError(err.Error(), nil)
	}

	dates := make(map[string]interface{})
	dates["F_list"] = res
	u.renderSuccessJSON(dates)
	return
}

type getKeyWorldsParam struct {
	F_user_id string `form:"user_id" valid:"Required;"`
}

// GetKeyWorlds 获取搜索历史记录
func (u *GoodsController) GetKeyWorlds() {
	params := &getKeyWorldsParam{}
	if !u.paramsValid(params) {
		u.renderParamsError(nil)
		return
	}

	var banner *models.MGood
	res, err := banner.GetKeyWorlds(params.F_user_id)

	if err != nil {
		u.renderUnknownError(err.Error(), nil)
	}

	dates := make(map[string]interface{})
	dates["result"] = res
	u.renderSuccessJSON(dates)
	return
}

type getGoodDetailParam struct {
	F_ID int `form:"id" valid:"Required;"`
}

// GetGoodDetail 获取商品详情
func (u *GoodsController) GetGoodDetail() {
	params := &getGoodDetailParam{}
	if !u.paramsValid(params) {
		u.renderParamsError(nil)
		return
	}

	var banner *models.MGood
	res, err := banner.GetGoodDetail(params.F_ID)

	if err != nil {
		u.renderUnknownError(err.Error(), nil)
	}

	dates := make(map[string]interface{})
	dates["result"] = res
	u.renderSuccessJSON(dates)
	return
}
