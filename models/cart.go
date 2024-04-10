// @Author Zihao_Li 2024/3/19 13:41:00
package models

import (
	"fmt"
	"food-service/infrastructure/db/mysql"
	"food-service/infrastructure/helper"
	"github.com/pkg/errors"
	"strings"
	"time"
)

type MCart struct{}

type Cart struct {
	F_id          int       `json:"id" gorm:"column:F_id;primary_key;"`
	F_open_id     string    `json:"open_id" gorm:"column:F_open_id"`
	F_good_id     int       `json:"good_id" gorm:"column:F_good_id"`
	F_price       float64   `json:"price" gorm:"column:F_price"`
	F_tag_id      int       `json:"tag_id" gorm:"column:F_tag_id"`
	F_tag_name    string    `json:"tag_name" gorm:"column:F_tag_name"`
	F_good_num    int       `json:"good_num" gorm:"column:F_good_num"`
	F_add_time    time.Time `json:"add_time" gorm:"column:F_add_time"`
	F_is_checkout int       `json:"is_checkout" gorm:"column:F_is_checkout"`
	F_isdel       int       `json:"isdel" gorm:"column:F_isdel"`
}

func (Cart) TableName() string {
	return "rb_cart"
}

type GetCartListResp struct {
	F_title        string  `json:"title" gorm:"column:F_title"`
	F_price        float32 `json:"price" gorm:"column:F_price"`
	F_price_module string  `json:"price_module" gorm:"column:F_price_module"`
	F_url          string  `json:"url" gorm:"column:F_url"`
	F_open_id      string  `json:"open_id" gorm:"column:F_open_id"`
	F_good_id      int     `json:"id" gorm:"column:F_good_id"`
	F_good_num     int     `json:"good_num" gorm:"column:F_good_num"`
}

func (c *MCart) GetCartList(openId string) ([]GetCartListResp, error) {
	db := mysql.GetDbDefault()

	// 校验用户id
	var user MUser
	num, err := user.VerifyUserID(openId)
	if err != nil {
		return nil, err
	}

	if num == 0 {
		return nil, errors.New("该用户不存在")
	}

	// 获取结果
	var result []GetCartListResp
	db.Raw(`SELECT *
FROM rb_cart
JOIN rb_goods ON rb_cart.F_good_id = rb_goods.F_id
WHERE rb_cart.F_open_id = ?
AND rb_cart.F_is_checkout = 0
AND rb_cart.F_isdel = 0
ORDER BY rb_cart.F_add_time DESC;`, openId).Find(&result)

	return result, nil
}

func (c *MCart) AddGoodToCart(openId string, goodId, goodNum int, goodPrice float64) error {
	db := mysql.GetDbDefault()

	// 校验id
	var user MUser
	num, err := user.VerifyUserID(openId)
	if err != nil {
		return err
	}

	if num == 0 {
		return errors.New("该用户不存在")
	}

	var mGood MGood
	num, err = mGood.VerifyGoodID(goodId)
	if err != nil {
		return err
	}
	if num == 0 {
		return errors.New("该商品不存在")
	}

	// 增加至购物车内
	tx := db.Begin()

	var midCart []Cart
	cart := Cart{
		F_open_id:     openId,
		F_good_id:     goodId,
		F_is_checkout: 0,
		F_isdel:       0,
	}
	if err := tx.Where(cart).Find(&midCart).Error; err != nil {
		tx.Rollback()
		return errors.Wrap(err, "查询购物车记录失败")

	}
	switch len(midCart) {
	case 0:
		//	没有存储记录
		cart.F_price = float64(goodNum) * goodPrice
		cart.F_good_num = goodNum
		cart.F_add_time = time.Now()

		// 补全种类列
		var good Good
		if err := tx.Table(Good{}.TableName()).Where("F_id = ?", goodId).First(&good).Error; err != nil {
			tx.Rollback()
			return errors.Wrap(err, "查询商品种类失败")
		}
		cart.F_tag_id = good.F_tag_id
		cart.F_tag_name = good.F_tag_name
		if err := tx.Create(&cart).Error; err != nil {
			tx.Rollback()
			return errors.Wrap(err, "加入购物车失败")
		}

	case 1:
		//	已填入至购物车
		midCart[0].F_good_num = goodNum + midCart[0].F_good_num
		midCart[0].F_price = float64(midCart[0].F_good_num) * goodPrice
		midCart[0].F_add_time = time.Now()
		if err := tx.Save(midCart[0]).Error; err != nil {
			tx.Rollback()
			return errors.Wrap(err, "加入购物车失败")
		}

	default:
		tx.Rollback()
		return errors.Wrap(err, "查询该购物车数据同一件商品不止一条数据存储")

	}

	tx.Commit()

	return nil
}

func (c *MCart) DelGoodToCart(id int) error {
	db := mysql.GetDbDefault()

	tx := db.Begin()
	where := Cart{
		F_id:          id,
		F_is_checkout: 0,
		F_isdel:       0,
	}
	if err := tx.Table(Cart{}.TableName()).Where(where).Update("F_isdel", 1).Error; err != nil {
		tx.Rollback()
		return errors.Wrap(err, "更新购物车状态失败")
	}

	tx.Commit()

	return nil
}

func (c *MCart) BuyGoods(id, openId string) error {
	db := mysql.GetDbDefault()

	goodStrList := strings.Split(id, ",")

	if len(goodStrList) != 0 {
		goodIntList, err := helper.ConvertStringSliceToIntSlice(goodStrList)
		if err != nil {
			return err
		}
		db = db.Where("F_good_id in (?)", goodIntList)

	}
	if err := db.Table(Cart{}.TableName()).Where("F_open_id = ? and F_isdel = 0 ", openId).Update("F_is_checkout", 1).Error; err != nil {
		return errors.Wrap(err, "更新购物车状态失败")
	}

	return nil
}

// 直接购买
func (c *MCart) DirectBuyGood(openId string, goodID int, price float64, tagID int, tagName string, goodNum int) error {
	db := mysql.GetDbDefault()
	var muser *MUser
	num, err := muser.VerifyUserID(openId)
	if err != nil {
		return err
	}
	if num != 1 {
		return errors.New("用户信息有误，请联系管理员")
	}

	tx := db.Begin()
	var tableGoodNum int
	if err := tx.Table(Good{}.TableName()).Where("F_id = ?", goodID).Pluck("F_num", &tableGoodNum).Error; err != nil {
		tx.Rollback()
		return errors.Wrap(err, "查询商品库存数量失败")
	}
	if goodNum > tableGoodNum {
		tx.Rollback()
		return errors.New("商品库存数量不足")
	}

	cart := Cart{
		F_open_id:     openId,
		F_good_id:     goodID,
		F_price:       price,
		F_tag_id:      tagID,
		F_tag_name:    tagName,
		F_good_num:    goodNum,
		F_add_time:    time.Now(),
		F_is_checkout: 1,
	}

	if err := tx.Table(Cart{}.TableName()).Create(&cart).Error; err != nil {
		tx.Rollback()
		return errors.Wrap(err, "创建购物记录失败")
	}

	if err := tx.Table(Good{}.TableName()).Where("F_id = ?", goodID).Update("F_num", tableGoodNum-goodNum).Error; err != nil {
		tx.Rollback()
		return errors.Wrap(err, "更新商品数量失败")
	}

	tx.Commit()
	return nil
}

type saleTagRankResp struct {
	F_tag_id   int64  `json:"tag_id" gorm:"column:F_id"`
	F_tag_name string `json:"tag_name" gorm:"column:F_tag_name"`
	F_num      int    `json:"num" gorm:"column:F_num"`
}

// SaleTagRank 获取商品类型销量排行
func (c *MCart) SaleTagRank(timeRange, flag string) ([]saleTagRankResp, error) {
	db := mysql.GetDbDefault()

	timeSlice := strings.Split(timeRange, ",")
	if len(timeSlice) != 2 {
		return nil, errors.New("时间参数格式错误")
	}
	afterDay, err := helper.AddOneDay(timeSlice[1])
	if err != nil {
		return nil, errors.Wrap(err, "加1天时间失败")
	}

	if flag == "" {
		// 默认销量从高到低
		flag = "DESC"
	}

	if flag == "DESC" {
		db = db.Order("F_num DESC")
	} else {
		db = db.Order("F_num")
	}

	var result []saleTagRankResp
	if err := db.Table(Cart{}.TableName()).
		Select("F_tag_id, F_tag_name, SUM(F_good_num) AS F_num").
		Where("F_add_time between ? and ? AND F_is_checkout = 1 AND F_isdel = 0", timeSlice[0], afterDay).
		Group("F_tag_id, F_tag_name").
		Find(&result).Error; err != nil {
		return nil, errors.Wrap(err, "获取销量类型失败")
	}
	return result, nil
}

type saleGoodRankResp struct {
	F_good_id int64  `json:"good_id" gorm:"column:F_good_id"`
	F_tag_id  int64  `json:"tag_id" gorm:"column:F_tag_id"`
	F_title   string `json:"title" gorm:"column:F_title"`
	F_url     string `json:"url" gorm:"column:F_url"`
	F_num     int    `json:"num" gorm:"column:F_num"`
}

// SaleGoodRank 获取商品销量排行
func (c *MCart) SaleGoodRank(timeRange, flag string, tagID int) ([]saleGoodRankResp, error) {
	db := mysql.GetDbDefault()

	timeSlice := strings.Split(timeRange, ",")
	if len(timeSlice) != 2 {
		return nil, errors.New("时间参数格式错误")
	}
	afterDay, err := helper.AddOneDay(timeSlice[1])
	if err != nil {
		return nil, errors.Wrap(err, "加1天时间失败")
	}

	var tagSql string
	if tagID != 0 {
		tagSql = fmt.Sprintf("AND rb_cart.F_tag_id = %d", tagID)
	}

	if flag == "" {
		flag = "DESC"
	}
	sql := fmt.Sprintf(`SELECT F_good_id, F_title, rb_cart.F_tag_id, F_url, SUM(F_good_num) AS F_num
FROM rb_cart
JOIN rb_goods ON rb_cart.F_good_id = rb_goods.F_id
WHERE F_add_time BETWEEN ? AND ?
  AND F_is_checkout = 1
  AND rb_cart.F_isdel = 0
  %s
GROUP BY F_good_id, F_title, rb_cart.F_tag_id, F_url
ORDER BY F_num %s;`, tagSql, flag)

	var result []saleGoodRankResp
	if err := db.Raw(sql, timeSlice[0], afterDay).Find(&result).Error; err != nil {
		return nil, errors.New("查询result失败")
	}

	return result, nil
}

// SaleSingleTotal 获取商品销量排行
func (c *MCart) SaleSingleTotal(timeRange string, tagID int) (int64, string, error) {
	db := mysql.GetDbDefault()

	timeSlice := strings.Split(timeRange, ",")
	if len(timeSlice) != 2 {
		return 0, "", errors.New("时间参数格式错误")
	}
	afterDay, err := helper.AddOneDay(timeSlice[1])
	if err != nil {
		return 0, "", errors.Wrap(err, "加1天时间失败")
	}
	tagName := "总营业额"
	if tagID != 0 {
		db.Table(Category{}.TableName()).Select("F_name").Where("F_id = ?", tagID).Scan(&tagName)
		db = db.Where("F_tag_id = ?", tagID)
		tagName += "类总营业额"
	}

	var total int64
	if err := db.Table(Cart{}.TableName()).
		Select("SUM(F_price) as F_price").
		Where("F_add_time BETWEEN ? and ? AND F_is_checkout = 1 AND  F_isdel = 0", timeSlice[0], afterDay).
		Scan(&total).Error; err != nil {
		return 0, "", errors.Wrap(err, "获取营业额失败")
	}

	return total, tagName, nil
}

type saleListTotalResp struct {
	F_tag_id int64  `json:"tag_id" gorm:"column:F_tag_id"`
	F_name   string `json:"name" gorm:"column:F_tag_name"`
	F_total  int64  `json:"total" gorm:"column:F_total"`
}

// SaleListTotal 获取销售额情况 (数组)
func (c *MCart) SaleListTotal(timeRange string) ([]saleListTotalResp, error) {
	db := mysql.GetDbDefault()

	timeSlice := strings.Split(timeRange, ",")
	if len(timeSlice) != 2 {
		return nil, errors.New("时间参数格式错误")
	}
	afterDay, err := helper.AddOneDay(timeSlice[1])
	if err != nil {
		return nil, errors.Wrap(err, "加1天时间失败")
	}

	// 获取总营业额
	total, totalName, err := c.SaleSingleTotal(timeRange, 0)
	if err != nil {
		return nil, err
	}
	var result []saleListTotalResp

	if err := db.Table(Cart{}.TableName()).Select("F_tag_id, F_tag_name, SUM(F_price) as F_total").Where("F_add_time BETWEEN ? and ?  AND F_is_checkout = 1 AND F_isdel = 0", timeSlice[0], afterDay).Group("F_tag_id, F_tag_name").Find(&result).Error; err != nil {
		return nil, err
	}

	result = append([]saleListTotalResp{{
		F_tag_id: 0,
		F_name:   totalName,
		F_total:  total,
	}}, result...)

	return result, nil
}
