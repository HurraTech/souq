package main

import (
	"fmt"

	"github.com/jessevdk/go-flags"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"

	"hurracloud.io/souq/internal/controller"
	"hurracloud.io/souq/internal/database"
)

type Options struct {
	Host        string         `short:"h" long:"host" env:"HOST" description:"Host to bind HTTP server to" default:"127.0.0.1"`
	Port        int            `short:"p" long:"port" env:"PORT" description:"Port to listen HTTP server" default:"5060"`
	Database    flags.Filename `short:"d" long:"db" env:"DB" description:"Database filename" default:"souq.db"`
	MetadataDir string         `short:"m" long:"metadta_dir" env:"METADATA_DIR" description:"Where to store metadta about applications" default:"./apps"`
	Verbose     bool           `short:"v" long:"verbose" description:"Enable verbose logging"`
}

var options Options

func main() {
	_, err := flags.Parse(&options)

	if err != nil {
		panic(err)
	}

	if options.Verbose {
		log.SetLevel(log.DebugLevel)
	}

	database.OpenDatabase(string(options.Database))
	database.Migrate()

	controller := &controller.Controller{MetadataDir: options.MetadataDir}
	e := echo.New()
	e.GET("/apps", controller.GetApps)
	e.GET("/apps/:id", controller.DownloadApp)
	log.Fatal(e.Start(fmt.Sprintf("%s:%d", options.Host, options.Port)))
}
