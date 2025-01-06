package db

import (
	"fmt"

	"gorm.io/gorm"
)

var db *gorm.DB

func GetDb() *gorm.DB {
	return db
}


func InitDb(){
	
}
