package server

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/mitchellh/cli"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"syscall"
	"todo/storage/mysql"
)

type Command struct{}

func CommandFactory() (cli.Command, error) {
	return new(Command), nil
}

type Service struct {
	mode       string
	httpServer *http.Server
	storage    *gorm.DB
}

func NewService(args ...string) *Service {
	serviceMode := "debug"
	if len(args) > 0 && args[0] == "release" {
		serviceMode = "release"
	}

	gin.SetMode(serviceMode)

	handler := gin.Default()
	handler.Use(CorsMiddleware())

	addr := fmt.Sprintf(":%d", 8000)
	server := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	mysql.InitMysql()
	service := &Service{
		mode:       serviceMode,
		httpServer: server,
		storage:    mysql.DB,
	}

	service.initRouter(handler)

	return service
}

func (s *Service) Run() int {

	startRestService := func() {
		fmt.Println("Start rest api server", s.httpServer.Addr)

		if err := s.httpServer.ListenAndServe(); err != nil {
			fmt.Errorf("RestServer.Run %s", err)
		}
		fmt.Println("rest server shutdown.")
	}

	waitToStop := func() {
		signalChan := make(chan os.Signal, 1)
		signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
		<-signalChan
	}

	stopAllService := func() {
		s.httpServer.Shutdown(context.Background())
		fmt.Println("all service has stopped, exit.")
	}

	go startRestService()

	waitToStop()
	stopAllService()

	return 1
}

func (c *Command) Run(args []string) int {
	return NewService(args...).Run()
}

func (c *Command) Help() string {
	return help
}

func (c *Command) Synopsis() string {
	return synopsis
}

const synopsis = "front restful api service tips"
const help = `
Usage: front xxxx
`

func CorsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")
		filterHost := []string{
			"http://localhost",
		}
		var isAccess = false
		for _, v := range filterHost {
			match, _ := regexp.MatchString(v, origin)
			if match {
				isAccess = true
			}
		}
		if isAccess {
			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
			c.Header("Access-Control-Allow-Methods", "GET, OPTIONS, POST, PUT, DELETE")
			c.Set("content-type", "application/json")
		}
		if method == "OPTIONS" {
			c.JSON(http.StatusOK, "Options Request!")
		}

		c.Next()
	}
}
