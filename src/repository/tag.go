package repository

import (
	"gorm.io/gorm"
)

type Tag struct {
	gorm.Model
	Name     string     `gorm:"unique" json:"name"`
	ImageURL string     `json:"url"`
	Articles []*Article `gorm:"many2many:article_tags;" json:"articles"`
}

func (tag *Tag) Save() error {
	result := db.Create(tag)
	return result.Error
}

func FindByTagNames(name string) ([]Tag, error) {
	var tags []Tag
	result := db.Where("name = ?", name).Find(&tags)
	return tags, result.Error
}

func FindRelatedByTagNames(name string, pageNation int) (*Tag, error) {
	var tag Tag
	err := db.Model(&Tag{}).Preload("Articles").Preload("Articles.Tags").Offset((pageNation-1)*30).Limit(30).Where("name = ?", name).Find(&tag).Error
	return &tag, err
}
