// @Author Zihao_Li 2024/3/18 14:58:00
package models

import (
	"food-service/infrastructure/db/mysql"
	"github.com/pkg/errors"
)

type MCategory struct{}

type Category struct {
	F_id   int    `json:"id" gorm:"column:F_id;primary_key;"`
	F_name string `json:"name" gorm:"column:F_name"`
}

func (Category) TableName() string {
	return "rb_category"
}

// GetCategoryList 获取所有分类列表
func (g *MCategory) GetCategoryList() ([]Category, error) {
	db := mysql.GetDbDefault()

	res := []Category{{
		F_id:   0,
		F_name: "全部商品",
	}}

	var result []Category
	if err := db.Table(Category{}.TableName()).Find(&result).Error; err != nil {
		return nil, errors.New("内部查询失败，请联系管理员")
	}
	result = append(res, result...)
	return result, nil
}

// GetGoodsByCategory 获取所有分类列表
func (g *MCategory) GetGoodsByCategory(id int) ([]Good, error) {
	db := mysql.GetDbDefault()

	db1 := db.Table(Good{}.TableName())
	if id != 0 {
		// 校验id 是否正常
		var num int64
		if err := db.Table(Category{}.TableName()).Where("F_id = ?", id).Count(&num).Error; err != nil {
			return nil, errors.Wrap(err, "校验类id失败")
		}
		if num == 0 {
			return nil, errors.New("改类id不存在")
		}
		db1 = db1.Where("F_tag_id = ?", id)
	}

	var result []Good
	if err := db1.Where("F_isdel = 0").Find(&result).Error; err != nil {
		return nil, errors.New("内部查询失败，请联系管理员")

	}

	return result, nil
}
