package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	DB *gorm.DB
)

func OpenDatabase(dbFile string) {
	var err error
	DB, err = gorm.Open(sqlite.Open(dbFile), &gorm.Config{})
	if err != nil {
		panic("failed to open database")
	}
}

func Migrate() {
	// Not using DB for now
}
