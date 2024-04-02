// @Author Zihao_Li 2024/3/20 23:47:00
package controllers

import "food-service/models"

type ExpirationController struct {
	BaseController
}

type expirationListParam struct {
	F_tag_id int `form:"tag_id"`
	F_flag   int `form:"flag"` // 0表示临期食品 1表示过期食品
}

// ExpirationList 获取临期或过期食品列表
func (e *ExpirationController) ExpirationList() {
	params := &expirationListParam{}
	if !e.paramsValid(params) {
		e.renderParamsError(nil)
		return
	}

	var expiration *models.MExpiration
	res, err := expiration.ExpirationList(params.F_tag_id, params.F_flag)

	if err != nil {
		e.renderUnknownError(err.Error(), nil)
	}

	dates := make(map[string]interface{})
	dates["result"] = res
	e.renderSuccessJSON(dates)
	return
}

// RemoveExpirationFood 批量下架临期或过期食品
func (e *ExpirationController) RemoveExpirationFood() {
	params := &expirationListParam{}
	if !e.paramsValid(params) {
		e.renderParamsError(nil)
		return
	}

	var expiration *models.MExpiration
	err := expiration.RemoveExpirationFood(params.F_tag_id, params.F_flag)

	if err != nil {
		e.renderUnknownError(err.Error(), nil)
	}

	e.jsonEchoOk()
	return
}

type changeExpirationPriceParam struct {
	F_id     int     `form:"id"`                       // 商品id
	F_tag_id int     `form:"tag_id"`                   // 类型id
	F_number float32 `form:"number" valid:"Required;"` // 数字（可能为价格，可能为打折比例）
	F_flag   int     `form:"flag"`                     // 0:表示直接价格，1:表示打折比例
}

// ChangeExpirationPrice 批量调整临期食品价格
func (e *ExpirationController) ChangeExpirationPrice() {
	params := &changeExpirationPriceParam{}
	if !e.paramsValid(params) {
		e.renderParamsError(nil)
		return
	}

	var expiration *models.MExpiration
	err := expiration.ChangeExpirationPrice(params.F_id, params.F_tag_id, params.F_number, params.F_flag)

	if err != nil {
		e.renderUnknownError(err.Error(), nil)
	}

	e.jsonEchoOk()
	return
}
