package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

var port = flag.Int("port", 8081, "Port in which the microservice will run")

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

type request struct {
	Dir  string `json:"dir"`
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	// Set up router
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Routes
	r.POST("/list", func(c *gin.Context) {
		var b request
		c.Bind(&b)
		files := dirFiles(b.Dir)
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data":    files,
		})
	})

	// Start server
	fmt.Println("Running server on port", *port)
	r.Run(fmt.Sprintf(":%d", *port))
}
