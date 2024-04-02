package routers

import (
	"food-service/controllers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/plugins/cors"
)

func init() {
	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
		AllowCredentials: true,
	}))

	// ************** 用户端 ****************** //
	// 轮播图
	beego.Router("api/banner", &controllers.BannerController{}, "get:GetBanner")
	// 更新商品
	beego.Router("api/goods/update", &controllers.GoodsController{}, "put:UpdateGood")
	// 删除商品
	beego.Router("api/goods/delete", &controllers.GoodsController{}, "delete:DeleteGood")
	// 添加商品
	beego.Router("api/goods/add", &controllers.GoodsController{}, "post:AddGood")
	// 商品列表
	beego.Router("api/goods", &controllers.GoodsController{}, "get:GetGoodList")
	// 商品详情
	beego.Router("api/goods/details", &controllers.GoodsController{}, "get:GetGoodDetail")
	// 商品模糊搜索列表
	beego.Router("api/goods/search", &controllers.GoodsController{}, "get:GetFuzzSearchGoodList")
	// 搜索历史记录
	beego.Router("api/keywords", &controllers.GoodsController{}, "get:GetKeyWorlds")
	// 获取分类列表
	beego.Router("api/category/list", &controllers.CategoryController{}, "get:GetCategoryList")
	// 获取该类下的所有商品
	beego.Router("api/category/goods", &controllers.CategoryController{}, "get:GetGoodsByCategory")
	// 添加购物车
	beego.Router("api/cart/add", &controllers.CartController{}, "put:AddGoodToCart")
	// 购物车查询
	beego.Router("api/cart/get", &controllers.CartController{}, "get:GetCartList")
	// 删除购物车
	beego.Router("api/cart/del", &controllers.CartController{}, "delete:DelGoodToCart")
	// 购物车结算
	beego.Router("api/cart/buy", &controllers.CartController{}, "post:BuyGoods")
	// 获取销售额情况 (单个)
	beego.Router("api/goods/sale/single", &controllers.CartController{}, "get:SaleSingleTotal")
	// 获取销售额情况 (数组)
	beego.Router("api/goods/sale/list", &controllers.CartController{}, "get:SaleListTotal")
	// 获取商品种类整体销售情况
	beego.Router("api/goods/sale/tag/rank", &controllers.CartController{}, "get:SaleTagRank")
	// 获取具体商品销量情况
	beego.Router("api/goods/sale/good/rank", &controllers.CartController{}, "get:SaleGoodRank")
	// 获取临期或过期食品列表
	beego.Router("api/expiration/list", &controllers.ExpirationController{}, "get:ExpirationList")
	// 批量下架临期或过期食品
	beego.Router("api/expiration/remove", &controllers.ExpirationController{}, "put:RemoveExpirationFood")
	// 批量调整临期食品价格
	beego.Router("api/expiration/change/price", &controllers.ExpirationController{}, "post:ChangeExpirationPrice")
	// 登录
	beego.Router("api/login", &controllers.UserController{}, "post:Login")

	// ************** 骑手端 ****************** //

}
