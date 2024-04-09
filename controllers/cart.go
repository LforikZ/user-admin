// @Author Zihao_Li 2024/3/19 13:41:00
package controllers

import (
	"food-service/models"
)

type CartController struct {
	BaseController
}

type getCartListParam struct {
	F_open_id string `form:"open_id" valid:"Required;"`
}

// GetCartList 获取用户购物车数据
func (u *CartController) GetCartList() {
	params := &getCartListParam{}
	if !u.paramsValid(params) {
		u.renderParamsError(nil)
		return
	}

	var cart *models.MCart
	res, err := cart.GetCartList(params.F_open_id)

	if err != nil {
		u.renderUnknownError(err.Error(), nil)
	}

	dates := make(map[string]interface{})
	dates["F_result"] = res
	u.renderSuccessJSON(dates)
	return
}

type addGoodToCartParam struct {
	F_open_id    string  `form:"open_id" valid:"Required;"`
	F_good_id    int     `form:"good_id" valid:"Required"`
	F_good_num   int     `form:"good_num" valid:"Required;"`
	F_good_price float64 `form:"good_price" valid:"Required;"`
}

// AddGoodToCart 添加购物车
func (u *CartController) AddGoodToCart() {
	params := &addGoodToCartParam{}
	if !u.paramsValid(params) {
		u.renderParamsError(nil)
		return
	}

	var cart *models.MCart
	err := cart.AddGoodToCart(params.F_open_id, params.F_good_id, params.F_good_num, params.F_good_price)

	if err != nil {
		u.renderUnknownError(err.Error(), nil)
	}

	u.jsonEchoOk()
	return
}

type delGoodToCartParam struct {
	F_id int `form:"id" valid:"Required;"`
}

// DelGoodToCart 删除购物车
func (u *CartController) DelGoodToCart() {
	params := &delGoodToCartParam{}
	if !u.paramsValid(params) {
		u.renderParamsError(nil)
		return
	}

	var cart *models.MCart
	err := cart.DelGoodToCart(params.F_id)

	if err != nil {
		u.renderUnknownError(err.Error(), nil)
	}

	u.jsonEchoOk()
	return
}

type buyGoodsParam struct {
	F_id      string `form:"id" valid:"Required;"`
	F_open_id string `form:"open_id" valid:"Required;"`
}

// BuyGoods 购物车结算
func (u *CartController) BuyGoods() {
	params := &buyGoodsParam{}
	if !u.paramsValid(params) {
		u.renderParamsError(nil)
		return
	}

	var cart *models.MCart
	err := cart.BuyGoods(params.F_id, params.F_open_id)

	if err != nil {
		u.renderUnknownError(err.Error(), nil)
	}

	u.jsonEchoOk()
	return
}

type directBuyGoodParam struct {
	F_open_id  string  `form:"open_id" valid:"Required;"`
	F_good_id  int     `form:"good_id" valid:"Required;"`
	F_price    float64 `form:"price" valid:"Required;"`
	F_tag_id   int     `form:"tag_id" valid:"Required;"`
	F_tag_name string  `form:"tag_name" valid:"Required;"`
	F_good_num int     `form:"good_num" valid:"Required;"`
}

// DirectBuyGood 直接购买结算
func (u *CartController) DirectBuyGood() {
	params := &directBuyGoodParam{}
	if !u.paramsValid(params) {
		u.renderParamsError(nil)
		return
	}

	var cart *models.MCart
	err := cart.DirectBuyGood(params.F_open_id, params.F_good_id, params.F_price, params.F_tag_id, params.F_tag_name, params.F_good_num)

	if err != nil {
		u.renderUnknownError(err.Error(), nil)
	}

	u.jsonEchoOk()
	return
}

type saleTagRankParam struct {
	F_time string `form:"time"` // 时间范围（不传表示所有时间，格式：2023-1-1,2023-1-31）
	F_flag string `form:"flag"` // 销量正，反统计标志（ASC，DESC）
}

// SaleTagRank 获取商品类型销量排行
func (u *CartController) SaleTagRank() {
	params := &saleTagRankParam{}
	if !u.paramsValid(params) {
		u.renderParamsError(nil)
		return
	}

	var cart *models.MCart
	res, err := cart.SaleTagRank(params.F_time, params.F_flag)

	if err != nil {
		u.renderUnknownError(err.Error(), nil)
	}

	dates := make(map[string]interface{})
	dates["result"] = res
	u.renderSuccessJSON(dates)
	return
}

type saleGoodRankParam struct {
	F_time   string `form:"time" valid:"Required;"` // 时间范围（不传表示所有时间，格式：2023-1-1,2023-1-31）
	F_flag   string `form:"flag"`                   // 销量正，反统计标志（ASC，DESC）
	F_tag_id int    `form:"tag_id"`                 // 种类id （不传表示全部）
}

// SaleGoodRank 获取商品销量排行
func (u *CartController) SaleGoodRank() {
	params := &saleGoodRankParam{}
	if !u.paramsValid(params) {
		u.renderParamsError(nil)
		return
	}

	var cart *models.MCart
	res, err := cart.SaleGoodRank(params.F_time, params.F_flag, params.F_tag_id)

	if err != nil {
		u.renderUnknownError(err.Error(), nil)
	}

	dates := make(map[string]interface{})
	dates["result"] = res
	u.renderSuccessJSON(dates)
	return
}

type saleSingleTotalParam struct {
	F_time   string `form:"time" valid:"Required;"` // 时间范围（不传表示所有时间，格式：2023-1-1,2023-1-31）
	F_tag_id int    `form:"tag_id"`                 // 种类id （不传表示总额）
}

// SaleSingleTotal 获取销售额情况 (单个)
func (u *CartController) SaleSingleTotal() {
	params := &saleSingleTotalParam{}
	if !u.paramsValid(params) {
		u.renderParamsError(nil)
		return
	}

	var cart *models.MCart
	total, name, err := cart.SaleSingleTotal(params.F_time, params.F_tag_id)

	if err != nil {
		u.renderUnknownError(err.Error(), nil)
	}

	dates := make(map[string]interface{})
	dates["total"] = total
	dates["name"] = name
	u.renderSuccessJSON(dates)
	return
}

type saleListTotalParam struct {
	F_time string `form:"time" valid:"Required;"` // 时间范围（不传表示所有时间，格式：2023-1-1,2023-1-31）
}

// SaleListTotal 获取销售额情况 (数组)
func (g *CartController) SaleListTotal() {
	params := &saleListTotalParam{}
	if !g.paramsValid(params) {
		g.renderParamsError(nil)
		return
	}

	var cart *models.MCart
	res, err := cart.SaleListTotal(params.F_time)

	if err != nil {
		g.renderUnknownError(err.Error(), nil)
	}

	dates := make(map[string]interface{})
	dates["result"] = res
	g.renderSuccessJSON(dates)
	return
}
