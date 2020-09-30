package main

import (
	"crypto/subtle"
	"fmt"
	"path/filepath"

	"github.com/jessevdk/go-flags"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	log "github.com/sirupsen/logrus"

	"hurracloud.io/souq/internal/controller"
	"hurracloud.io/souq/internal/database"
)

type Options struct {
	Host       string         `short:"h" long:"host" env:"HOST" description:"Host to bind HTTP server to" default:"127.0.0.1"`
	Port       int            `short:"p" long:"port" env:"PORT" description:"Port to listen HTTP server" default:"5060"`
	Database   flags.Filename `short:"d" long:"db" env:"DB" description:"Database filename" default:"souq.db"`
	AppsDir    string         `short:"m" long:"apps_dir" env:"APPS_DIR" description:"Where to store metadta about applications" default:"./apps"`
	EnableAuth bool           `short:"a" long:"enable_auth" env:"ENABLE_AUTH" description:"Sets up Basic Auth in front of API"`
	Verbose    bool           `short:"v" long:"verbose" description:"Enable verbose logging"`
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

	appsDir, err := filepath.Abs(options.AppsDir)
	if err != nil {
		log.Warnf("Could not determine absolute path for metdata directory '%s': %s", options.AppsDir, err)
		appsDir = options.AppsDir
	}

	controller := &controller.Controller{AppsDir: appsDir}
	e := echo.New()
	e.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		// Be careful to use constant time comparison to prevent timing attacks
		if subtle.ConstantTimeCompare([]byte(username), []byte("HURRANET")) == 1 &&
			subtle.ConstantTimeCompare([]byte(password), []byte("bSdh~e9J:FTbLS#w")) == 1 {
			return true, nil
		}
		return false, nil
	}))

	e.GET("/apps", controller.GetApps)
	e.GET("/apps/:id", controller.GetApp)
	e.GET("/apps/:id/image", controller.DownloadApp)
	e.GET("/apps/:id/containers", controller.ListAppContainers)
	e.GET("/containers/:app/:container", controller.DownloadAppContainerImage)
	log.Fatal(e.Start(fmt.Sprintf("%s:%d", options.Host, options.Port)))
}
