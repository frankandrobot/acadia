package main

import (
	"fmt"
	"github.com/frankandrobot/acadia/file"
	"github.com/frankandrobot/acadia/messaging"
	"github.com/frankandrobot/acadia/queue"
	"github.com/labstack/echo"
	"net/http"
)

// CONFIG

type ServerRoot string
type ServerKey string
type ServerCer string

var serverRoot = ServerRoot("/home/pi/data")
var serverKey = ServerKey("/home/pi/server.key")
var serverCert = ServerCer("/home/pi/server.crt")

func main() {
	router := echo.New()
	queue := queue.Queue()

	router.GET("/files", func(c echo.Context) error {
		payload := messaging.Payload{
			Action: func() messaging.Result {
				return file.LoadDir(string(serverRoot))
			},
			Done: make(chan messaging.Result),
		}
		go func() { queue <- payload }()
		result := <-payload.Done
		if result.Error != nil {
			return c.String(http.StatusInternalServerError, fmt.Sprintf("%s", result.Error))
		}
		return c.JSON(http.StatusOK, result.Filenames)
	})

	// router.GET("/files/:name", func(c *gin.Context) {
	// 	name := c.Param("name")
	// 	filename := fmt.Sprintf("%s/%s", string(serverRoot), name)
	// 	payload := fileCommandPayload{}
	// 	payload.Context = c.Copy()
	// 	payload.Filename = filename
	// 	payload.HandleResult = func(result fileCommandResult) {
	// 		if result.Error == nil {
	// 			result.Context.String(http.StatusOK, "%s", result.Contents)
	// 		} else {
	// 			result.Context.String(http.StatusNotFound, "")
	// 		}
	// 	}
	// 	channels.Load <- payload
	// })

	// router.POST("/files/:name", func(c *gin.Context) {
	// 	context := c.Copy()
	// 	name := c.Param("name")
	// 	var json contents
	// 	err := c.BindJSON(&json)
	// 	if err != nil {
	// 		c.String(http.StatusInternalServerError, "%s", err)
	// 		return
	// 	}
	// 	payload := fileCommandPayload{}
	// 	payload.Filename = name
	// 	payload.Contents = json.Contents
	// 	payload.HandleResult = func(result fileCommandResult) {
	// 		if result.Error == nil {
	// 			context.String(http.StatusOK, "")
	// 		} else {
	// 			context.String(http.StatusInternalServerError, "%s", result.Error)
	// 		}
	// 	}
	// 	channels.Save <- payload
	// })

	router.StartTLS(":8080", string(serverCert), string(serverKey))
}
