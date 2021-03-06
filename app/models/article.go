package models

import (
	"Kronos/library/databases"
	"errors"
)

type Article struct {
	BaseModel

	Title          string           `gorm:"type:varchar(100);"`
	Keyword        string           `gorm:"type:varchar(100);"`
	Description    string           `gorm:"type:varchar(100);"`
	Thumb          string           `gorm:"size:255"` // 设置字段大小为255
	ArticleContent []ArticleContent `gorm:"foreignkey:article_id;association_foreignkey:id;"`
	Category       []Category       `gorm:"many2many:article_cates;"`
	Tags           []Tags           `gorm:"many2many:article_tags;"`
}

func (a Article) Count(where string, vals []interface{}) (int, error) {
	var count = 0
	err := databases.DB.Model(&a).Where(where, vals).Count(&count).Error
	if err != nil {
		return count, errors.New("暂无数据可查")
	}
	return count, nil
}
func (a Article) Lists(where string, vals []interface{}, offset, limit int) ([]Article, error) {
	list := make([]Article, limit)
	// .Preload("ArticleContent")
	err := databases.DB.Model(a).Where(where, vals).Offset(offset).Limit(limit).Find(&list)
	if err.Error != nil {
		return nil, errors.New("暂无数据可查")
	}
	return list, nil
}

func (a Article) Trash(where string, vals []interface{}, offset, limit int) ([]Article, error) {
	list := make([]Article, limit)
	err := databases.DB.Model(a).Unscoped().Where(where, vals).Offset(offset).Limit(limit).Find(&list)
	if err.Error != nil {
		return nil, errors.New("暂无数据可查")
	}
	return list, nil
}

func (a Article) Get(where string, vals []interface{}) (Article, error) {
	first := databases.DB.Model(a).Preload("ArticleContent").Preload("Category").Preload("Tags").Where(where, vals).First(&a)
	if first.Error != nil {
		return a, first.Error
	}
	return a, nil
}

func (a Article) Update(id uint64, data map[string]interface{}) error {
	var find Article
	first := databases.DB.Model(&a).Where("id = ?", id).First(&find)
	if first.Error != nil {
		return first.Error
	}
	update := databases.DB.Model(&find)
	var cate []Category
	databases.DB.Where("id in (?)", data["category_ids"]).Find(&cate)
	update.Association("Category").Replace(cate)
	var tag []Tags
	databases.DB.Where("id in (?)", data["tag_ids"]).Find(&tag)
	update.Association("Tags").Replace(tag)

	if err := databases.DB.Model(&find).Update(a).Error; err != nil {
		return err
	}

	//association := update.Association("ArticleContent").Replace(a.ArticleContent)
	//if association.Error != nil {
	//	return association.Error
	//}
	return nil
}
func (a Article) Create(data map[string]interface{}) error {
	create := databases.DB.Model(&a).Create(&a)
	err := create.Association("ArticleContent").Append(a.ArticleContent).Error
	if err != nil {
		return err
	}
	var cate []Category
	databases.DB.Where("id in (?)", data["category_ids"]).Find(&cate)
	create.Association("Category").Append(cate)
	var tag []Tags
	databases.DB.Where("id in (?)", data["tag_ids"]).Find(&tag)
	create.Association("Tags").Append(tag)
	return nil
}

func (m Article) Delete(id uint64) error {
	first := databases.DB.Model(&m).Where("id = ?", id).First(&m)
	if err := first.Error; err != nil {
		return err
	}
	first.Delete(&m)

	return nil
}

func (m Article) ForceDelete(id uint64) error {
	first := databases.DB.Model(&m).Unscoped().Preload("Category").Preload("ArticleContent").Preload("Tags").Where("id = ?", id).First(&m)
	if err := first.Error; err != nil {
		return err
	}

	databases.DB.Model(&m).Unscoped().Delete(m)
	databases.DB.Model(&m).Unscoped().Association("ArticleContent").Clear()
	databases.DB.Model(&m).Unscoped().Association("Category").Delete(m.Category)
	databases.DB.Model(&m).Unscoped().Association("Tags").Delete(m.Tags)

	return nil
}
