package provider

import (
	"irisblog/config"
	"irisblog/model"
)

func GetArticleById(id uint) (*model.Article, error) {
	var article model.Article
	db := config.DB
	err := db.Where("`id`=?", id).First(&article).Error
	if err != nil {
		return nil, err
	}
	article.ArticleData = &model.ArticleData{}
	db.Where("`id`=?", id).First(article.ArticleData)
	article.Category = &model.Category{}
	db.Where("`id`=?", article.CategoryId).First(&article.Category)
	return &article, nil
}

func GetPrevArticleById(categoryId, id uint) (*model.Article, error) {
	var article model.Article
	db := config.DB
	err := db.Where("`category_id`=?", categoryId).Where("id BETWEEN ? AND ?", int(id)-1000, int(id)-1).Where("`status`=1").Last(&article).Error
	if err != nil {
		return nil, err
	}
	return &article, nil
}

func GetNextArticleById(categoryId, id uint) (*model.Article, error) {
	var article model.Article
	db := config.DB
	err := db.Where("`category_id`=?", categoryId).Where("id BETWEEN ? AND ?", id+1, id+1000).Where("`status`=1").First(&article).Error
	if err != nil {
		return nil, err
	}
	return &article, nil
}

func GetArticleList(categoryId uint, order string, currentPage, pageSize int) ([]*model.Article, error) {
	var article []*model.Article
	offset := (currentPage - 1) * pageSize
	builder := config.DB.Debug().Model(&model.Article{})
	if categoryId > 0 {
		builder = builder.Where("`category_id`=?", categoryId)
	}

	if order != "" {
		builder = builder.Order(order)
	}

	err := builder.Limit(pageSize).Offset(offset).Find(&article).Error
	if err != nil {
		return nil, err
	}

	return article, nil
}
