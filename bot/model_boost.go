package bot

import "gorm.io/gorm"

type Boost struct {
	gorm.Model
	PostID int `gorm:"primaryKey"`
	UserID int `gorm:"primaryKey"`
}
