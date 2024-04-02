// @Author Zihao_Li 2024/3/15 21:20:00
package models

import (
	"errors"
	"food-service/infrastructure/db/mysql"
)

type MBanner struct{}

type Banner struct {
	FId    int    `json:"id" gorm:"column:F_id;primary_key;"`
	FUrl   string `json:"url" gorm:"column:F_url;"`
	FDescs string `json:"descs" gorm:"column:F_descs"`
}

func (Banner) TableName() string {
	return "rb_banner"
}

func (f *MBanner) GetBanner() ([]Banner, error) {
	db := mysql.GetDbDefault()

	var banners []Banner
	if err := db.Table(Banner{}.TableName()).Find(&banners).Error; err != nil {
		return nil, errors.New("内部查询失败，请联系管理员")
	}

	return banners, nil
}
