package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"slices"

	"github.com/gin-gonic/gin"
)

var port = flag.Int("port", 8080, "Port in which the microservice will run")

func dirFiles(dir string) []string {
	entries, err := os.ReadDir(dir)
	var files []string
	if err != nil {
		log.Println(err)
	} else {
		for _, entry := range entries {
			files = append(files, entry.Name())
		}
	}
	return files
}

func download(dir, file string) (string, error) {
	files := dirFiles(dir)
	i := slices.Index(files, file)
	if i == -1 {
		return "", errors.New("file not found")
	}
	path := dir + file
	if dir[len(dir)-1] != os.PathSeparator {
		path = dir + string(os.PathSeparator) + file
	}
	return fmt.Sprint(os.Stat(path)), nil
}

type request struct {
	File string `json:"file"`
	Dir  string `json:"dir"`
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	// Set up router
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Routes
	r.POST("/download", func(c *gin.Context) {
		var b request
		c.Bind(&b)
		string, err := download(b.Dir, b.File)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"error":   err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data":    string,
		})
	})

	// Start server
	fmt.Println("Running server on port", *port)
	r.Run(fmt.Sprintf(":%d", *port))
}
