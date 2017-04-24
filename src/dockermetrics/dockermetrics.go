package main

import (
	"common"
	"config"
	"dao"
	"flag"
	"fmt"
	"github.com/cloudfoundry/gosteno"
	"github.com/codegangsta/cli"
	"github.com/gin-gonic/gin"
	"logger"
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
		logger, _, conf := loadLoggerAndConfig(c, "dockermetrics")
		common.Logger = logger
		flag.Parse()
		CommandServiceInstance := dao.NewCommandService(conf, logger)
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

func loadLoggerAndConfig(c *cli.Context, component string) (logger.Logger, *gosteno.Logger, *config.Config) {
	var conf *config.Config
	var err error
	conf, err = config.FromEnv()
	if err != nil {
		fmt.Printf("Failed to load config from env : %s", err.Error())
		os.Exit(1)
	}
	stenoConf := &gosteno.Config{Sinks: []gosteno.Sink{gosteno.NewIOSink(os.Stdout)}, Level: conf.LogLevel(), Codec: gosteno.NewJsonCodec()}
	gosteno.Init(stenoConf)
	steno := gosteno.NewLogger(component)
	sysdigLogger := logger.NewRealLogger(steno)
	return sysdigLogger, steno, conf
}
