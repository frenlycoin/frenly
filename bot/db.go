package bot

import (
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Initializes DB object
func initDb() *gorm.DB {
	var db *gorm.DB
	var err error

	if conf.SQLite {
		db, err = gorm.Open(sqlite.Open(conf.DbURI), &gorm.Config{})
	} else {
		if conf.Dev {
			db, err = gorm.Open(postgres.Open(conf.DbURI), &gorm.Config{
				Logger: logger.Default.LogMode(logger.Info),
			})
		} else {
			db, err = gorm.Open(postgres.Open(conf.DbURI), &gorm.Config{})
		}
	}

	if err != nil {
		loge(err)
	}

	if err := db.AutoMigrate(&User{}, &Transaction{}, &KeyValue{}, &Channel{}, &Post{}, &Boost{}); err != nil {
		panic(err.Error())
	}

	err = db.SetupJoinTable(&User{}, "Boosts", &Boost{})

	if err != nil {
		loge(err)
	}

	return db
}
