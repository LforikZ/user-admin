// @Author Zihao_Li 2024/3/18 14:57:00
package controllers

import "food-service/models"

type CategoryController struct {
	BaseController
}

// GetCategoryList 获取所有分类列表
func (g *CategoryController) GetCategoryList() {

	var category *models.MCategory
	res, err := category.GetCategoryList()

	if err != nil {
		g.renderUnknownError(err.Error(), nil)
	}

	dates := make(map[string]interface{})
	dates["result"] = res
	g.renderSuccessJSON(dates)
	return
}

type getGoodsByCategoryParam struct {
	F_id int `form:"id"`
}

// GetGoodsByCategory 获取分类列表下的所有商品
func (g *CategoryController) GetGoodsByCategory() {
	params := &getGoodsByCategoryParam{}
	if !g.paramsValid(params) {
		g.renderParamsError(nil)
		return
	}

	var category *models.MCategory
	res, err := category.GetGoodsByCategory(params.F_id)

	if err != nil {
		g.renderUnknownError(err.Error(), nil)
	}

	dates := make(map[string]interface{})
	dates["result"] = res
	g.renderSuccessJSON(dates)
	return
}
