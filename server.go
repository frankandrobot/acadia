package main

import (
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

type contents struct {
	Contents string `form:"contents" json:"contents" binding:"required"`
}

func main() {
	router := echo.New()
	queue := queue.MakeQueue()

	router.GET("/files", func(c echo.Context) error {
		action := func() messaging.ChanResult {
			files, err := file.LoadDir(string(serverRoot))
			result := messaging.ChanResult{}
			result.Filenames = files
			result.Error = err
			return result
		}
		result := queue.Add(action)
		if result.Error != nil {
			return result.Error
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

	router.POST("/files/:name", func(c echo.Context) error {
		name := c.Param("name")
		doc := new(contents)
		if err := c.Bind(doc); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}
		action := func() messaging.ChanResult {
			err := file.SaveFile(
				file.Root(serverRoot),
				file.BaseFilename(name),
				file.Contents(doc.Contents),
			)
			result := messaging.ChanResult{}
			result.Error = err
			return result
		}
		result := queue.Add(action)
		if result.Error != nil {
			return result.Error
		}
		return c.JSON(http.StatusOK, doc)
	})

	router.StartTLS(":8080", string(serverCert), string(serverKey))
}
