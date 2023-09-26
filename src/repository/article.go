package repository

import (
	"gorm.io/gorm"
	"time"
)

type Article struct {
	gorm.Model
	ZennSlug       string     `gorm:"not null" gorm:"uniqueIndex" json:"zenn_slug"`
	Title          string     `gorm:"unique" json:"title"`
	ImageUrl       string     `gorm:"not null" json:"image"`
	Url            string     `gorm:"not null" json:"url"`
	PublishedAt    time.Time  `gorm:"not null" json:"date"`
	Tags           []*Tag     `gorm:"many2many:article_tags;" json:"tags"`
	UserLiked      []*User    `gorm:"many2many:like_articles;" json:"good"`
	UserBookMarked []*User    `gorm:"many2many:book_mark_articles;" json:"book_marked"`
	Comments       []*Comment `gorm:"many2many:comment_articles;" json:"comment"`
}

func (article *Article) Save() error {
	result := db.Create(article)
	return result.Error
}

func (article *Article) Update() error {
	result := db.Save(article)
	return result.Error
}

func FindByTitles(title string) ([]Article, error) {
	var articles []Article
	result := db.Where("title = ?", title).Find(&articles)
	return articles, result.Error
}

func FindByArticles(pageNation int) (*[]Article, error) {
	var articles []Article
	var result *gorm.DB
	if pageNation < 1 {
		result = db.Model(&Article{}).Preload("Tags").Preload("UserLiked").Preload("UserBookMarked").Order("created_at desc").Offset(1).Limit(30).Find(&articles)
	} else {

		result = db.Model(&Article{}).Preload("Tags").Preload("UserLiked").Preload("UserBookMarked").Order("created_at desc").Offset((pageNation - 1) * 30).Limit(30).Find(&articles)
	}
	return &articles, result.Error
}

func FindLikeArticleById(articleId uint) (*Article, error) {
	var article Article
	err := db.Preload("UserLiked").Preload("UserBookMarked").First(&article, "id = ?", articleId).Error
	return &article, err
}

func FindLikeArticleByIds(articleIds []uint) (*[]Article, error) {
	var articles []Article
	err := db.Preload("UserLiked").Preload("UserBookMarked").Find(&articles, articleIds).Error
	return &articles, err
}

func FindBookMarkArticleById(articleId uint) (*Article, error) {
	var article Article
	err := db.Preload("UserLiked").Preload("UserBookMarked").First(&article, "id = ?", articleId).Error
	return &article, err
}
