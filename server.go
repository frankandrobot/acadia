package main

import "net/http"
import "io/ioutil"
import "os"
import "github.com/gin-gonic/gin"
import . "github.com/ahmetb/go-linq"

type ServerRoot string

var serverRoot = ServerRoot{"/home/pi/data"}

func main() {
	router := gin.Default()

	router.GET("/files", func(c *gin.Context) {
		files, err := ioutil.ReadDir(serverRoot)
		if err != nil {
			c.String(http.StatusInternalServerError, "%s", err)
		} else {
			var filenames []string
			From(files).
				SelectT(func(c os.FileInfo) string { return c.Name() }).
				ToSlice(&filenames)
			c.JSON(200, filenames)
		}
	})

	router.GET("/pass", func(c *gin.Context) {
		//name := c.Param("name")
		file, err := ioutil.ReadFile("./test.txt")
		if err == nil {
			c.String(http.StatusOK, "Hello %s", file)
		} else {
			c.String(http.StatusOK, "Bye")
		}
	})

	// However, this one will match /user/john/ and also /user/john/send
	// If no other routers match /user/john, it will redirect to /user/john/
	router.GET("/user/:name/*action", func(c *gin.Context) {
		name := c.Param("name")
		action := c.Param("action")
		message := name + " is " + action
		c.String(http.StatusOK, message)
	})

	router.Run(":8080")
}
