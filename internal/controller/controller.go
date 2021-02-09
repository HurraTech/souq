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
	AppsDir string
	OSDir   string
}

type App struct {
	UniqueID        string
	Name            string `yaml:"Name"`
	Description     string `yaml:"Description"`
	LongDescription string `yaml:"LongDescription"`
	Publisher       string `yaml:"Publisher"`
	Version         string `yaml:"Version"`
	WebApp          WebApp `yaml:"WebApp"`
	Icon            string
	Containers      string
}

type WebApp struct {
	Type            string `yaml:"Type"`
	TargetPort      int    `yaml:"TargetPort"`
	TargetContainer string `yaml:"TargetContainer"`
}

/* GET /apps */
func (c *Controller) GetApps(ctx echo.Context) error {

	appsMetadata, err := filepath.Glob(path.Join(c.AppsDir, "**", "metadata.yml"))
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
	arch := ctx.QueryParam("arch")
	imageFile := path.Join(c.AppsDir, fmt.Sprintf("%s-%s.tar.gz", appID, arch))

	log.Debugf("Request to download %s", imageFile)
	_, err := os.Stat(imageFile)
	if err != nil {
		if os.IsNotExist(err) {
			return ctx.JSON(http.StatusNotFound, map[string]string{"messsage": "not found"})
		}
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": "unexpected error"})
	}

	return ctx.File(imageFile)
}

/* GET /containers/:app/:container */
func (c *Controller) DownloadAppContainerImage(ctx echo.Context) error {
	appID := ctx.Param("app")
	containerName := ctx.Param("container")
	arch := ctx.QueryParam("arch")
	imageFile := path.Join(c.AppsDir, appID, "containers", fmt.Sprintf("%s-%s.tar.gz", containerName, arch))

	log.Debugf("Request to download %s", imageFile)
	_, err := os.Stat(imageFile)
	if err != nil {
		if os.IsNotExist(err) {
			return ctx.JSON(http.StatusNotFound, map[string]string{"messsage": "not found"})
		}
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": "unexpected error"})
	}

	return ctx.File(imageFile)
}

/* GET /apps/:id/containers */
func (c *Controller) ListAppContainers(ctx echo.Context) error {
	appID := ctx.Param("id")
	containersFile := path.Join(c.AppsDir, appID, "containers", "containers.yml")

	_, err := os.Stat(containersFile)
	if err != nil {
		if os.IsNotExist(err) {
			return ctx.JSON(http.StatusNotFound, map[string]string{"message": "no containers found"})
		}
		log.Errorf("Error while checking for containers.yml: %s", err)
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": "unexpected error"})
	}

	return ctx.File(containersFile)
}

func (c *Controller) readAppMetadata(name string) (*App, error) {
	// Read metadata file
	metadataPath := path.Join(c.AppsDir, name, "metadata.yml")
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
	containersPath := path.Join(c.AppsDir, name, "containers", "containers.yml")
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

/* GET /hurraos */
func (c *Controller) ListHurraOSVersions(ctx echo.Context) error {
	log.Debugf("Listing directory %s", c.OSDir)
	files, err := ioutil.ReadDir(c.OSDir)
	if err != nil {
		log.Error("Error while listing OS update files: ", err)
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": "unexpected error"})
	}

	res := make(map[string]interface{})
	for _, f := range files {
		match, _ := filepath.Match("*.img", f.Name())
		if !match {
			continue
		}

		stat, err := os.Stat(path.Join(c.OSDir, f.Name()))
		if err != nil {
			log.Error("Error while stating file: ", f.Name())
			return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": "unexpected error"})
		}

		// Read full image's sha file
		file, err := os.Open(path.Join(c.OSDir, fmt.Sprintf("%s.sha", f.Name())))
		if err != nil {
			log.Errorf("Error while opening SHA file %s.sha: %s: %s", f.Name(), err)
			continue
		}
		defer func() {
			if err = file.Close(); err != nil {
				log.Errorf("Error while closing SHA file %s.sha: %s", f.Name(), err)
			}
		}()

		imageSha, err := ioutil.ReadAll(file)
		if err != nil {
			log.Errorf("Error while reading SHA file %s.sha: %s: %s", f.Name(), err)
			continue
		}

		// Read update image's sha file
		file2, err := os.Open(path.Join(c.OSDir, fmt.Sprintf("%s.mender.sha", f.Name())))
		if err != nil {
			log.Errorf("Error while opening SHA file %s.mender.sha: %s: %s", f.Name(), err)
			continue
		}
		defer func() {
			if err = file2.Close(); err != nil {
				log.Errorf("Error while closing SHA file %s.mender.sha: %s", f.Name(), err)
			}
		}()
		updateImageSha, err := ioutil.ReadAll(file2)
		if err != nil {
			log.Errorf("Error while reading SHA file %s.mender.sha: %s: %s", f.Name(), err)
			continue
		}

		ver := strings.TrimSuffix(f.Name(), filepath.Ext(f.Name()))
		res[ver] = map[string]string{
			"version":          ver,
			"release_date":     stat.ModTime().Format("2006-01-02"),
			"full_image":       fmt.Sprintf("%s://%s/%s/files/%s", ctx.Scheme(), ctx.Request().Host, filepath.Base(c.OSDir), f.Name()),
			"full_image_sha":   strings.TrimSpace(string(imageSha[:])),
			"bmap":             fmt.Sprintf("%s://%s/%s/files/%s", ctx.Scheme(), ctx.Request().Host, filepath.Base(c.OSDir), fmt.Sprintf("%s.bmap", f.Name())),
			"update_image":     fmt.Sprintf("%s://%s/%s/files/%s", ctx.Scheme(), ctx.Request().Host, filepath.Base(c.OSDir), fmt.Sprintf("%s.mender", f.Name())),
			"update_image_sha": strings.TrimSpace(string(updateImageSha[:])),
		}
	}

	return ctx.JSON(http.StatusOK, res)
}

/* GET /hurraos/:version */
func (c *Controller) GetHurraOSVersionInfo(ctx echo.Context) error {
	ver := ctx.Param("version")
	img := fmt.Sprintf("%s.img", ver)
	stat, err := os.Stat(path.Join(c.OSDir, img))
	if err != nil {
		log.Error("Error while stating file: ", img)
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": "unexpected error"})
	}

	// Read full image's sha file
	file, err := os.Open(path.Join(c.OSDir, fmt.Sprintf("%s.sha", img)))
	if err != nil {
		log.Errorf("Error while opening SHA file %s.sha: %s: %s", img, err)
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": "unexpected error"})
	}
	defer func() {
		if err = file.Close(); err != nil {
			log.Errorf("Error while closing SHA file %s.sha: %s", img, err)
		}
	}()

	imageSha, err := ioutil.ReadAll(file)
	if err != nil {
		log.Errorf("Error while reading SHA file %s.sha: %s: %s", img, err)
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": "unexpected error"})
	}

	// Read update image's sha file
	file2, err := os.Open(path.Join(c.OSDir, fmt.Sprintf("%s.mender.sha", img)))
	if err != nil {
		log.Errorf("Error while opening SHA file %s.mender.sha: %s: %s", img, err)
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": "unexpected error"})
	}
	defer func() {
		if err = file2.Close(); err != nil {
			log.Errorf("Error while closing SHA file %s.mender.sha: %s", img, err)
		}
	}()
	updateImageSha, err := ioutil.ReadAll(file2)
	if err != nil {
		log.Errorf("Error while reading SHA file %s.mender.sha: %s: %s", img, err)
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": "unexpected error"})
	}

	res := map[string]string{
		"version":          ver,
		"release_date":     stat.ModTime().Format("2006-01-02"),
		"full_image":       fmt.Sprintf("%s://%s/%s/files/%s", ctx.Scheme(), ctx.Request().Host, filepath.Base(c.OSDir), img),
		"full_image_sha":   strings.TrimSpace(string(imageSha[:])),
		"bmap":             fmt.Sprintf("%s://%s/%s/files/%s", ctx.Scheme(), ctx.Request().Host, filepath.Base(c.OSDir), fmt.Sprintf("%s.bmap", img)),
		"update_image":     fmt.Sprintf("%s://%s/%s/files/%s", ctx.Scheme(), ctx.Request().Host, filepath.Base(c.OSDir), fmt.Sprintf("%s.mender", img)),
		"update_image_sha": strings.TrimSpace(string(updateImageSha[:])),
	}

	return ctx.JSON(http.StatusOK, res)
}

/* GET /hurrsos/files/:image */
func (c *Controller) DownloadHurraOS(ctx echo.Context) error {
	image := ctx.Param("image")
	imageFile := path.Join(c.OSDir, image)

	log.Debugf("Request to download %s", imageFile)
	_, err := os.Stat(imageFile)
	if err != nil {
		if os.IsNotExist(err) {
			return ctx.JSON(http.StatusNotFound, map[string]string{"messsage": "not found"})
		}
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": "unexpected error"})
	}

	return ctx.File(imageFile)
}
