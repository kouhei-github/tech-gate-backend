package repository

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UserName          string `json:"userName"`
	Email             string `gorm:"not null" gorm:"unique" json:"email"`
	Password          string
	IsLogin           bool       `gorm:"default=false"`
	Image             string     `json:"image"`
	SelfIntroduction  string     `json:"selfIntroduction"`
	GithubUser        string     `json:"githubUser"`
	TwitterUser       string     `json:"twitterUser"`
	ArticleLiked      []*Article `gorm:"many2many:like_articles;" json:"article_liked"`
	ArticleBookMarked []*Article `gorm:"many2many:book_mark_articles;" json:"article_book_marked"`
	Comments          []Comment
}

func (receiver *User) Save() error {
	result := db.Create(receiver)
	return result.Error
}

func (receiver *User) FindByEmail(email string) ([]User, error) {
	var users []User
	result := db.Where("email = ?", email).Find(&users)
	return users, result.Error
}

func (receiver *User) FindById(userId uint) error {
	result := db.First(receiver, userId)
	return result.Error
}

func (receiver *User) Update() error {
	result := db.Save(receiver)
	return result.Error
}

func FindUserBookMarkArticle(userId uint64) (*User, error) {
	var user User
	//err := db.Preload("ArticleBookMarked.UserLiked").Preload("ArticleBookMarked.UserBookMarked").First(&user, "id = ?", userId).Error
	err := db.Preload("ArticleBookMarked").First(&user, "id = ?", userId).Error

	return &user, err
}

func FindUserLikeArticle(userId uint64) (*User, error) {
	var user User
	err := db.Preload("ArticleLiked").First(&user, "id = ?", userId).Error
	return &user, err
}
