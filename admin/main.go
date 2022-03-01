package main

import (
	"github.com/GoAdminGroup/filemanager"
	_ "github.com/GoAdminGroup/go-admin/adapter/gin" // web framework adapter
	"github.com/GoAdminGroup/go-admin/engine"
	_ "github.com/GoAdminGroup/go-admin/modules/db/drivers/mysql" // sql driver
	"github.com/GoAdminGroup/go-admin/template"
	"github.com/GoAdminGroup/go-admin/template/chartjs"
	_ "github.com/GoAdminGroup/themes/sword" // ui theme
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"todo/admin/models"
	"todo/admin/pages"
	"todo/admin/tables"
	"todo/storage/mysql"
)

func main() {
	startServer()
}

func startServer() {
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = ioutil.Discard

	r := gin.Default()

	template.AddComp(chartjs.NewChart())

	eng := engine.Default()

	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	err = os.Mkdir(filepath.Join(dir, "files"), os.FileMode(0755))
	if err != nil {
		panic(err)
	}

	if err := eng.AddConfigFromJSON("./config.json").
		AddGenerators(tables.Generators).
		AddPlugins(filemanager.NewFileManager(filepath.Join(dir, "files"))).
		Use(r); err != nil {
		panic(err)
	}

	r.Static("/uploads", "./uploads")

	eng.HTML("GET", "/admin", pages.GetDashBoard)
	eng.HTMLFile("GET", "/admin/hello", "./html/hello.tmpl", map[string]interface{}{
		"msg": "Hello world",
	})

	models.Init(eng.MysqlConnection())

	mysql.InitMysql()

	err = r.Run(":20080")
	if err != nil {
		panic(err)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Print("closing database connection")
	eng.MysqlConnection().Close()
}
