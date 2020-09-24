package controller

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

type Controller struct {
	MetadataDir string
}

type App struct {
	ID              string
	Name            string
	Description     string
	LongDescription string `yaml:"long_description"`
	Publisher       string
	Version         string
	Icon            string
}

/* GET /apps */
func (c *Controller) GetApps(ctx echo.Context) error {

	appsMetadata, err := filepath.Glob(path.Join(c.MetadataDir, "**", "metadata.yml"))
	if err != nil {
		log.Errorf("Unexpected scanning metadata directory: %s", err)
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": "unexpected error"})
	}

	var apps []App
	for _, metadataPath := range appsMetadata {
		contents, err := ioutil.ReadFile(metadataPath)
		if err != nil {
			log.Errorf("Error opening :%s: %s", metadataPath, err)
			continue
		}

		app := App{}
		err = yaml.Unmarshal(contents, &app)
		if err != nil {
			log.Errorf("Error parsing metadata file: %s: %s", metadataPath, err)
			continue
		}

		appIcon, err := ioutil.ReadFile(path.Join(path.Dir(metadataPath), "icon.svg"))
		if err != nil {
			log.Errorf("Error reading apps icon file: %s: %s", metadataPath, err)
			continue
		}
		app.ID = path.Base(path.Dir(metadataPath))
		app.ID = strings.TrimSuffix(app.ID, path.Ext(app.ID))
		app.Icon = strings.Replace(string(appIcon), "<svg ", "<svg class=\"appStoreIcon\" ", 1)
		apps = append(apps, app)
	}

	log.Debugf("APPS: %v", apps)
	return ctx.JSON(http.StatusOK, apps)
}

/* GET /apps/:id */
func (c *Controller) DownloadApp(ctx echo.Context) error {
	appID := ctx.Param("id")
	imageFile := path.Join(c.MetadataDir, fmt.Sprintf("%s.tar.gz", appID))

	log.Debugf("Request to download %s", imageFile)
	_, err := os.Stat(imageFile)
	if err != nil {
		if os.IsNotExist(err) {
			return ctx.JSON(http.StatusNotFound, map[string]string{"messsage": "not found"})
		}
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": "unexpected error"})
	}

	return ctx.File(path.Join(c.MetadataDir, fmt.Sprintf("%s.tar.gz", appID)))
}
