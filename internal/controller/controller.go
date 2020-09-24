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

var (
	AppNotFound = fmt.Errorf("App Not Found")
)

type Controller struct {
	MetadataDir string
}

type App struct {
	UniqueID        string
	Name            string
	Description     string
	LongDescription string `yaml:"long_description"`
	Publisher       string
	Version         string
	Icon            string
	Containers      string
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
		appName := path.Base(path.Dir(metadataPath))

		app, err := c.readAppMetadata(appName)
		if err != nil {
			log.Errorf("Unexpected error while reading app %s metadata: %s", ctx.Param("id"), err)
			continue
		}
		apps = append(apps, *app)
	}

	return ctx.JSON(http.StatusOK, apps)
}

/* GET /apps/:id */
func (c *Controller) GetApp(ctx echo.Context) error {

	app, err := c.readAppMetadata(ctx.Param("id"))
	if err == AppNotFound {
		ctx.JSON(http.StatusNotFound, map[string]string{"message": "not found"})
	} else if err != nil {
		log.Errorf("Unexpected error while reading app %s metadata: %s", ctx.Param("id"), err)
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": "unexpected error"})
	}

	return ctx.JSON(http.StatusOK, app)
}

/* GET /apps/:id/image */
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

/* GET /apps/:id/containers/:container */
func (c *Controller) DownloadAppContainerImage(ctx echo.Context) error {
	appID := ctx.Param("id")
	containerName := ctx.Param("container")
	imageFile := path.Join(c.MetadataDir, appID, "containers", fmt.Sprintf("%s.tar.gz", containerName))

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

func (c *Controller) readAppMetadata(name string) (*App, error) {
	// Read metadata file
	metadataPath := path.Join(c.MetadataDir, name, "metadata.yml")
	if _, err := os.Stat(metadataPath); os.IsNotExist(err) {
		return nil, AppNotFound
	}

	contents, err := ioutil.ReadFile(metadataPath)
	if err != nil {
		return nil, fmt.Errorf("Error opening: %s: %s", metadataPath, err)
	}

	// Parse it
	app := App{}
	err = yaml.Unmarshal(contents, &app)
	if err != nil {
		return nil, fmt.Errorf("Error parsing metadata file: %s: %s", metadataPath, err)
	}

	// Check for containers.yml
	seen := make(map[string]bool)
	var images []string
	containersPath := path.Join(c.MetadataDir, name, "containers", "containers.yml")
	if _, err := os.Stat(containersPath); err == nil {
		contents, err := ioutil.ReadFile(containersPath)
		if err != nil {
			return nil, fmt.Errorf("Error opening: %s: %s", containersPath, err)
		}
		m := make(map[interface{}]interface{})
		err = yaml.Unmarshal(contents, &m)
		if err != nil {
			return nil, fmt.Errorf("Error parsing containers file: %s: %s", containersPath, err)
		}

		for _, service := range m["services"].(map[interface{}]interface{}) {
			image := service.(map[interface{}]interface{})["image"].(string)
			if _, ok := seen[image]; ok {
				continue
			}
			seen[image] = true
			images = append(images, image)
		}

	}

	// AppIcon (icon.svg)
	appIcon, err := ioutil.ReadFile(path.Join(path.Dir(metadataPath), "icon.svg"))
	if err != nil {
		return nil, fmt.Errorf("Error reading apps icon file: %s: %s", metadataPath, err)
	}

	app.UniqueID = path.Base(path.Dir(metadataPath))
	app.UniqueID = strings.TrimSuffix(app.UniqueID, path.Ext(app.UniqueID))
	app.Icon = strings.Replace(string(appIcon), "<svg ", "<svg class=\"appStoreIcon\" ", 1)
	app.Containers = strings.Join(images, ",")
	return &app, nil
}
