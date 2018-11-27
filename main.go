package main

//go:generate go-bindata -o static.go static/... views/...

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

const (
	envBind      = "BYT_BIND"
	envUploadDir = "BYT_UPLOAD_DIR"
	envHost      = "BYT_HOST"

	defaultBind      = ":8080"
	defaultUploadDir = "./upload"
	defaultHost      = "localhost:8080"

	cacheHeader = "public, max-age=86400"
)

var (
	bind      = defaultBind
	uploadDir = defaultUploadDir
	host      = defaultHost
)

func usage() {
	fmt.Printf("Usage: %s\n", os.Args[0])
	fmt.Printf("\tenv variables:\n")
	fmt.Printf("\t%s: Server bind address, default: %s\n", envBind, defaultBind)
	fmt.Printf("\t%s: Uploaded file save directory, default: %s\n", envUploadDir, defaultUploadDir)
	fmt.Printf("\t%s: Host url for uploaded file link, default: %s\n", envHost, defaultHost)
}

// envDefault returns an environment variable with `name`.
// returns default value if missing.
func envDefault(name, def string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}

	return def
}

func setupServer() *gin.Engine {

	s := gin.Default()
	s.GET("/static/*path", handleStatic)
	s.GET("/favicon.ico", handleAsset("static/favicon.ico"))

	s.GET("/", handleAsset("views/index.html"))

	s.POST("/upload", handleUpload)

	s.GET("/f/:id/*filename", handleFile)
	s.GET("/f/:id", handleFile)

	return s
}

func handleStatic(c *gin.Context) {

	path := filepath.Join("static", c.Param("path"))

	handleAsset(path)(c)
}

func handleAsset(assetPath string) gin.HandlerFunc {
	return func(c *gin.Context) {

		data, err := Asset(assetPath)
		if err != nil {
			if os.IsNotExist(err) {
				c.AbortWithStatus(http.StatusNotFound)
				return
			}
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.Header("Cache-Control", cacheHeader)
		contentType := DetectContentType(filepath.Base(assetPath), data)
		c.Data(http.StatusOK, contentType, data)
	}
}

func handleUpload(c *gin.Context) {
	header, err := c.FormFile("file")
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	f, err := header.Open()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer f.Close()

	id := uuid.NewV4().String()

	hdPath := filepath.Join(envDefault(envUploadDir, defaultUploadDir), id)
	fd, err := os.Create(hdPath)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer fd.Close()

	if _, err := io.Copy(fd, f); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	u := &url.URL{
		Scheme: RequestScheme(c.Request),
		Host:   envDefault(envHost, defaultHost),
		Path:   fmt.Sprintf("/f/%s/%s", id, url.QueryEscape(header.Filename)),
	}

	c.JSON(http.StatusCreated, gin.H{
		"file": u.String(),
	})
}

func handleFile(c *gin.Context) {

	id := c.Param("id")

	path := filepath.Join(envDefault(envUploadDir, defaultUploadDir), id)

	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Header("Cache-Control", cacheHeader)
	c.File(path)
}

func main() {

	flag.Usage = usage
	flag.Parse()

	s := setupServer()
	s.Run(envDefault(envBind, defaultBind))
}
