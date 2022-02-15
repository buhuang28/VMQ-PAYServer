package db

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var (
	DbLink *gorm.DB
)

func DbInit() {
	var err error
	DbLink, err = gorm.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/oauth?charset=utf8")
	if err != nil {
		fmt.Println(err)
	}
	//DbLink.LogMode(true)
	DbLink.SingularTable(true)
}

func init() {
	DbInit()
}
