package mysql

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"os"
	. "todo/config"
	"todo/model"
)

// DB Global DB connection
var DB *gorm.DB

func InitMysql() {
	var err error
	DB, err = gorm.Open("mysql",
		fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?parseTime=True&loc=Local",
			Config.Database.User,
			Config.Database.Password,
			Config.Database.Host,
			Config.Database.Port,
			Config.Database.Name))
	if err != nil {
		panic(err)
	}

	if os.Getenv("DEBUG") != "" {
		DB.LogMode(true)
	}
	DB.LogMode(true)

	AutoMigrate(
		&model.Project{},
		&model.Task{},
	)
}

func AutoMigrate(values ...interface{}) {
	for _, one := range values {
		DB.AutoMigrate(one)
	}
}
