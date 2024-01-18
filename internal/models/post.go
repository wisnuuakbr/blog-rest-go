package models

type Post struct {
	ID     uint   `gorm:"primarykey"`
	Title  string `gorm:"not null" json:"title"`
	Body   string `gorm:"type:text" json:"body"`
	UserID uint   `gorm:"foreignkey:UserID" json:"userID"`
	User   User   `gorm:"foreignkey:UserID"`
}