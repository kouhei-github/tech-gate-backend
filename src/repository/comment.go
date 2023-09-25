package repository

import "gorm.io/gorm"

type Comment struct {
	gorm.Model
	Message  string     `gorm:"not null" json:"message"`
	Articles []*Article `gorm:"many2many:comment_articles;" json:"articles"`
	UserID   uint       `json:"user_id"`
}
