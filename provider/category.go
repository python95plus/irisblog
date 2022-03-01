package provider

import (
	"irisblog/config"
	"irisblog/model"
)

func GetCategories(parentId uint) ([]*model.Category, error) {
	var categories []*model.Category
	db := config.DB
	builder := db.Where("`status` = ?", 1)

	if parentId > 0 {
		builder = builder.Where("`parent_id` = ?", parentId)
	}

	err := builder.Find(&categories).Error
	if err != nil {
		return nil, err
	}

	return categories, nil
}

func GetCategoryByTitle(title string) (*model.Category, error) {
	var category model.Category
	db := config.DB
	err := db.Where("`title`=?", title).First(&category).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}
