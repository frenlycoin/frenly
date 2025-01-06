package bot

import "gorm.io/gorm"

type KeyValue struct {
	gorm.Model
	Key      string `gorm:"size:255;uniqueIndex"`
	ValueInt int64  `gorm:"type:int"`
	ValueStr string `gorm:"type:string"`
}
