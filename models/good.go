// @Author Zihao_Li 2024/3/15 21:11:00
package models

import (
	"food-service/infrastructure/db/mysql"
	"food-service/infrastructure/helper"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"time"
)

type MGood struct{}

type Good struct {
	F_id           int       `json:"id" gorm:"column:F_id;primary_key;"`
	F_title        string    `json:"title" gorm:"column:F_title;"`
	F_num          int       `json:"num" gorm:"column:F_num"`
	F_price        float32   `json:"price" gorm:"column:F_price"`
	F_price_module string    `json:"price_module" gorm:"column:F_price_module"`
	F_url          string    `json:"url" gorm:"column:F_url"`
	F_url_detail   string    `json:"url_detail" gorm:"column:F_url_detail"`
	F_tag_id       int       `json:"tag_id" gorm:"column:F_tag_id"`         // 类别id
	F_tag_name     string    `json:"tag_name" gorm:"column:F_tag_name"`     // 类别名
	F_store        string    `json:"store" gorm:"column:F_store"`           // 存储方式
	F_area         string    `json:"area" gorm:"column:F_area"`             // 商品产地
	F_shelf_life   int       `json:"shelf_life" gorm:"column:F_shelf_life"` // 保质期，单位为（天），默认 1 天 到期
	F_sale_time    time.Time `json:"sale_time" gorm:"column:F_sale_time"`   // 上架时间
	F_isdel        int       `json:"isdel" gorm:"column:F_isdel"`
}

func (Good) TableName() string {
	return "rb_goods"
}

type KeyWords struct {
	F_id      uint      `json:"id" gorm:"primaryKey;column:F_id"`
	F_content string    `json:"content" gorm:"column:F_content"`
	F_user_id string    `json:"user_id" gorm:"column:F_user_id"`
	F_time    time.Time `json:"time" gorm:"column:F_time"`
}

func (KeyWords) TableName() string {
	return "rb_keywords"
}

// AddGood 添加商品
func (g *MGood) AddGood(F_title string, F_num int, F_price float32, F_price_module, F_url, F_url_detail string, F_tag_id int, F_tag_name, F_store, F_area string, F_shelf_life int) error {
	db := mysql.GetDbDefault()

	var temp []Good
	if err := db.Table(Good{}.TableName()).Where("F_title = ? and F_isdel = 0", F_title).Find(&temp).Error; err != nil {
		return errors.Wrap(err, "内部查询失败，请联系管理员")
	}
	if F_shelf_life == 0 {
		F_shelf_life = 1
	}
	good := Good{
		F_title:        F_title,
		F_num:          F_num,
		F_price:        F_price,
		F_price_module: F_price_module,
		F_url:          F_url,
		F_url_detail:   F_url_detail,
		F_tag_id:       F_tag_id,
		F_tag_name:     F_tag_name,
		F_store:        F_store,
		F_area:         F_area,
		F_shelf_life:   F_shelf_life,
		F_sale_time:    time.Now(),
		F_isdel:        0,
	}

	tx := db.Begin()
	if len(temp) == 0 {
		if err := db.Table(Good{}.TableName()).Create(&good).Error; err != nil {
			tx.Rollback()
			return err
		}
	} else if len(temp) == 1 {
		tx.Rollback()
		return errors.New("该商品已存在")

	} else {
		tx.Rollback()
		return errors.New("该商品不唯一,请联系管理员")
	}

	tx.Commit()

	return nil

}

// UpdateGood 更新商品
func (g *MGood) UpdateGood(F_id int, F_title string, F_num int, F_price float32, F_price_module, F_url, F_url_detail string, F_tag_id int, F_tag_name, F_store, F_area string, F_shelf_life int) error {
	db := mysql.GetDbDefault()

	var num int64
	if err := db.Table(Good{}.TableName()).Where("F_id != ? and F_title = ?  and F_isdel = 0", F_id, F_title).Count(&num).Error; err != nil {
		return errors.Wrap(err, "内部查询失败，请联系管理员")
	}
	if num != 0 {
		return errors.New("更新的名字与其他商品重名")
	}

	var temp Good
	if err := db.Table(Good{}.TableName()).Where("F_id = ? and F_isdel = 0", F_id).First(&temp).Error; err != nil {
		return errors.Wrap(err, "内部查询失败，请联系管理员")
	}
	if F_shelf_life == 0 {
		F_shelf_life = 1
	}

	tx := db.Begin()
	if temp.F_id != 0 {
		temp.F_title = F_title
		temp.F_num = F_num
		temp.F_price = F_price
		temp.F_price_module = F_price_module
		temp.F_url = F_url
		temp.F_url_detail = F_url_detail
		temp.F_tag_id = F_tag_id
		temp.F_tag_name = F_tag_name
		temp.F_store = F_store
		temp.F_area = F_area
		temp.F_shelf_life = F_shelf_life
		temp.F_sale_time = time.Now()

		if err := db.Table(Good{}.TableName()).Updates(&temp).Error; err != nil {
			tx.Rollback()
			return err
		}

	} else {
		tx.Rollback()
		return errors.New("商品不存在")
	}

	tx.Commit()

	return nil

}

// DeleteGood 删除商品
func (g *MGood) DeleteGood(F_id, F_flag int) error {
	db := mysql.GetDbDefault()

	tx := db.Begin()

	if F_flag == 1 {
		// 表示下架商品
		if err := tx.Table(Good{}.TableName()).Where("F_id = ?", F_id).Update("F_isdel", 1).Error; err != nil {
			tx.Rollback()
			return errors.Wrap(err, "更新商品信息失败")
		}

	} else {
		// 从数据库中删除商品
		if err := tx.Delete(&Good{}, F_id).Error; err != nil {
			tx.Rollback()
			return errors.New("删除商品失败")
		}

	}

	tx.Commit()

	return nil

}

// GetGoodList 获取商品列表
func (g *MGood) GetGoodList(page, limit int) ([]Good, error) {
	db := mysql.GetDbDefault()

	var result []Good

	//处理分页参数
	limit, page = helper.InitPage(limit, page)
	offset := helper.PageOffset(limit, page)

	if err := db.Table(Good{}.TableName()).Where("DATE_ADD(F_sale_time, INTERVAL F_shelf_life  DAY) >= DATE_ADD(CURDATE(), INTERVAL 1  DAY)").Limit(limit).Offset(offset).Find(&result).Error; err != nil {
		return nil, errors.New("内部查询失败，请联系管理员")

	}

	return result, nil

}

// GetFuzzSearchGoodList 获取模糊搜索商品列表
func (g *MGood) GetFuzzSearchGoodList(F_search string) ([]Good, error) {
	db := mysql.GetDbDefault()

	var result []Good

	if err := db.Table(Good{}.TableName()).Where("DATE_ADD(F_sale_time, INTERVAL F_shelf_life  DAY) >= DATE_ADD(CURDATE(), INTERVAL 1  DAY)").Where("F_title like ?", "%"+F_search+"%").Find(&result).Error; err != nil {
		return nil, errors.New("内部查询失败，请联系管理员")

	}

	return result, nil
}

// GetKeyWorlds 获取搜索关键字
func (g *MGood) GetKeyWorlds(userID string) ([]string, error) {
	db := mysql.GetDbDefault()

	// 校验用户是否存在
	var num int64
	if err := db.Table(User{}.TableName()).Where("F_open_id = ?", userID).Count(&num).Error; err != nil {
		return nil, errors.Wrap(err, "查询用户信息失败")
	}
	if num == 0 {
		return nil, errors.New("改用户不存在，请登录正确账号")
	}

	var result []string
	if err := db.Table(KeyWords{}.TableName()).Where("F_user_id = ?", userID).Order("F_time DESC").Limit(10).Pluck("F_content", &result).Error; err != nil {
		return nil, errors.Wrap(err, "查询关键词失败")
	}

	return result, nil
}

// AddKeyWorld 存储搜索关键字
func (g *MGood) AddKeyWorld(userID, content string) error {
	db := mysql.GetDbDefault()

	// 校验用户是否存在
	var num int64
	if err := db.Table(User{}.TableName()).Where("F_open_id = ?", userID).Count(&num).Error; err != nil {
		return errors.Wrap(err, "查询用户信息失败")
	}
	if num == 0 {
		return errors.New("改用户不存在，请登录正确账号")
	}

	var result KeyWords
	if err := db.Table(KeyWords{}.TableName()).Where("F_user_id = ? and F_content = ?", userID, content).First(&result).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.Wrap(err, "查询搜索记录失败")
	}
	if result.F_id == 0 {
		if err := db.Table(KeyWords{}.TableName()).Create(&KeyWords{
			F_content: content,
			F_user_id: userID,
			F_time:    time.Now(),
		}).Error; err != nil {
			return errors.Wrap(err, "存储搜索记录失败")
		}
	} else {
		result.F_time = time.Now()
		if err := db.Save(&result).Error; err != nil {
			return errors.Wrap(err, "更新搜索记录失败")
		}
	}

	return nil
}

// GetGoodDetail 获取商品详情
func (g *MGood) GetGoodDetail(id int) (Good, error) {
	db := mysql.GetDbDefault()

	var result Good

	if err := db.Table(Good{}.TableName()).Where("F_id = ? and F_isdel = 0 ", id).First(&result).Error; err != nil {
		return Good{}, errors.New("内部查询失败，请联系管理员")

	}

	return result, nil
}

// VerifyGoodID 校验商品id是否存在
func (g *MGood) VerifyGoodID(goodId int) (int64, error) {
	db := mysql.GetDbDefault()

	var num int64
	if err := db.Table(Good{}.TableName()).Where("F_id = ? ", goodId).Count(&num).Error; err != nil {
		return 0, errors.Wrap(err, "查询商品信息失败")
	}

	if num != 1 {
		return 0, errors.New("未查询到该商品")
	}

	return num, nil
}
