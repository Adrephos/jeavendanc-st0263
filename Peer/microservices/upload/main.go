package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

var port = flag.Int("port", 8082, "Port in which the microservice will run")

func createFile(dir, file string) {
	path := dir + file
	if dir[len(dir)-1] != os.PathSeparator {
		path = dir + string(os.PathSeparator) + file
	}

	fmt.Println(path)

	os.Create(path)
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
	r.POST("/upload", func(c *gin.Context) {
		var b request
		c.Bind(&b)
		createFile(b.Dir, b.File)
		c.JSON(http.StatusOK, gin.H{
			"success": true,
		})
	})

	// Start server
	fmt.Println("Running server on port", *port)
	r.Run(fmt.Sprintf(":%d", *port))
}
