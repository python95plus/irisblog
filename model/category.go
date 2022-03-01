package model

import (
	"time"

	"gorm.io/gorm"
)

type Category struct {
	Model
	Title       string `json:"title" gorm:"column:title;type:varchar(255) not null; default:''"`
	Description string `json:"description" gorm:"column:description;varchar(255) not null; default:''"`
	Content     string `json:"content" gorm:"column:content;teyp:longtext default null"`
	ParentId    uint   `json:"parent_id" gorm:"column:parent_id;type:int(10) unsigned not null;default:0;index:idx_parent_id"`
	Status      uint   `json:"status" gorm:"column:status;type:tinyint(1) unsigned not null;default:0;index:idx_status"`
}

func (category *Category) Save(db *gorm.DB) error {
	if category.Id == 0 {
		category.CreatedTime = time.Now().Unix()
	}
	category.UpdatedTime = time.Now().Unix()
	if err := db.Save(category).Error; err != nil {
		return err
	}
	return nil
}
