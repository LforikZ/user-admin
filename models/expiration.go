// @Author Zihao_Li 2024/3/20 23:48:00
package models

import (
	"fmt"
	"food-service/infrastructure/db/mysql"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type MExpiration struct{}

// ExpirationList 临期食品获取
func (g *MExpiration) ExpirationList(F_tag_id, F_flag int) ([]Good, error) {
	db := mysql.GetDbDefault()

	var tagSql string
	if F_tag_id != 0 {
		tagSql = fmt.Sprintf("AND F_tag_id = %d", F_tag_id)
	}

	flag := "="
	if F_flag != 0 {
		// 挑选临期食品
		flag = "<"
	}
	sql := fmt.Sprintf(`SELECT *
FROM rb_goods
WHERE  DATE_ADD(F_sale_time, INTERVAL F_shelf_life  DAY) %s DATE_ADD(CURDATE(), INTERVAL 1  DAY)  AND F_isdel = 0  %s ;`, flag, tagSql)

	var result []Good
	db.Raw(sql).Find(&result)

	return result, nil

}

// RemoveExpirationFood 批量下架临期或过期食品
func (g *MExpiration) RemoveExpirationFood(F_tag_id, F_flag int) error {
	db := mysql.GetDbDefault()

	tx := db.Begin()

	if F_tag_id != 0 {
		tx = tx.Where("F_tag_id = ?", F_tag_id)
	}

	if F_flag == 0 {
		// 挑选临期食品
		tx = tx.Where("DATE_ADD(F_sale_time, INTERVAL F_shelf_life  DAY) = DATE_ADD(CURDATE(), INTERVAL 1  DAY)")
	} else {
		// 过期产品
		tx = tx.Where("DATE_ADD(F_sale_time, INTERVAL F_shelf_life  DAY) < DATE_ADD(CURDATE(), INTERVAL 1  DAY)")
	}
	if err := tx.Table(Good{}.TableName()).Where("F_isdel = 0").Update("F_isdel", 1).Error; err != nil {
		tx.Rollback()
		return errors.Wrap(err, "下架商品失败")
	}
	tx.Commit()
	return nil
}

// ChangeExpirationPrice 批量调整临期食品价格
func (g *MExpiration) ChangeExpirationPrice(F_id, F_tag_id int, F_number float32, F_flag int) error {
	db := mysql.GetDbDefault()

	tx := db.Begin()
	if F_id != 0 {
		tx = tx.Where("F_id = ?", F_id)
	}
	if F_tag_id != 0 {
		tx = tx.Where("F_tag_id = ?", F_tag_id)
	}
	tx = tx.Where("DATE_ADD(F_sale_time, INTERVAL F_shelf_life  DAY) = DATE_ADD(CURDATE(), INTERVAL 1  DAY) AND F_isdel = 0")

	if F_flag == 0 {
		// 直接修改价格
		if err := tx.Table(Good{}.TableName()).Update("F_price", F_number).Error; err != nil {
			tx.Rollback()
			return errors.Wrap(err, "直接修改价格失败")
		}

	} else {
		// 进行打折
		if err := tx.Table(Good{}.TableName()).Update("F_price", gorm.Expr("F_price * ?", F_number/10)).Error; err != nil {
			tx.Rollback()
			return errors.Wrap(err, "打折失败")
		}
	}

	tx.Commit()
	return nil
}
