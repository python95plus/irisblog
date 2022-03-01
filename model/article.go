package model

import (
	"time"

	"gorm.io/gorm"
)

type Article struct {
	Model
	Title       string       `json:"title" gorm:"column:title;type:varchar(255) not null;default:''"`
	Keywords    string       `json:"keywords" gorm:"column:keywords;type:varchar(255) not null;default:''"`
	Description string       `json:"description" gorm:"column:description;type:varchar(255) not null;deafult:''"`
	CategoryId  uint         `json:"category_id" gorm:"column:category_id;type:int(10) unsigned not null;default:0;index:idx_category_id"`
	Views       uint         `json:"views" gorm:"column:views;type:int(10) unsigned not null;default:0;index:idx_views"`
	Status      uint         `json:"status" gorm:"column:status;type:tinyint(1) unsigned not null;default:0;index:idx_status"`
	Category    *Category    `json:"category" gorm:"-"`
	ArticleData *ArticleData `json:"data" gorm:"-"`
}

type ArticleData struct {
	Id      uint   `json:"id" gorm:"column:id;type:int(10) unsigned not null;primary_key"`
	Content string `json:"content" gorm:"column:content;type:longtext default null"`
}

func (article *Article) Save(db *gorm.DB) error {
	if article.Id == 0 {
		article.CreatedTime = time.Now().Unix()
	}
	article.UpdatedTime = time.Now().Unix()
	if err := db.Debug().Save(article).Error; err != nil {
		return err
	}
	if article.ArticleData != nil {
		article.ArticleData.Id = article.Id
		if err := db.Debug().Save(article.ArticleData).Error; err != nil {
			return err
		}
	}
	return nil
}

func (article *Article) AddViews(db *gorm.DB) error {
	article.Views++
	db.Model(Article{}).Where("`id`=?", article.Id).Update("views", article.Views)
	return nil
}
