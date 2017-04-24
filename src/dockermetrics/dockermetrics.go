package main

import (
	"config"
	"dao"
	"flag"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"time"
)

func main() {
	app := cli.NewApp()
	app.Name = "dockermetrics"
	app.Usage = "Start the docker metrics components"
	app.Version = "1.0.0"
	app.Flags = []cli.Flag{}
	app.Action = func(c *cli.Context) {
		log.Infof("%s", "starting")
		conf := loadLoggerAndConfig(c, "dockermetrics")
		flag.Parse()
		log.Infof("%s", "running")
		CommandServiceInstance := dao.NewCommandService(conf)
		CommandServiceInstance.Run()
		router := gin.Default()
		router.GET("/ping", func(c *gin.Context) {
			c.String(200, "pong")
		})
		server := &http.Server{
			Addr:           ":8080",
			Handler:        router,
			ReadTimeout:    300 * time.Second,
			WriteTimeout:   300 * time.Second,
			MaxHeaderBytes: 1 << 20,
		}
		server.ListenAndServe()
	}
	app.Run(os.Args)
}

func loadLoggerAndConfig(c *cli.Context, component string) *config.Config {
	var conf *config.Config
	var err error
	conf, err = config.FromEnv()
	if err != nil {
		fmt.Printf("Failed to load config from env : %s", err.Error())
		os.Exit(1)
	}
	return conf
}
