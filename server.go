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
	fileIO := file.NewFileIO()

	router.GET("/files", func(c echo.Context) error {
		action := func() messaging.ChanResult {
			files, err := fileIO.LoadDir(string(serverRoot))
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

	router.GET("/files/:name", func(c echo.Context) error {
		name := c.Param("name")
		action := func() messaging.ChanResult {
			contents, err := fileIO.LoadFile(
				file.Root(serverRoot),
				file.BaseFilename(name),
			)
			result := messaging.ChanResult{}
			result.Contents = contents
			result.Error = err
			return result
		}
		result := queue.Add(action)
		if result.Error != nil {
			return result.Error
		}
		return c.JSON(http.StatusOK, contents{Contents: result.Contents})
	})

	router.POST("/files/:name", func(c echo.Context) error {
		name := c.Param("name")
		doc := new(contents)
		if err := c.Bind(doc); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}
		action := func() messaging.ChanResult {
			err := fileIO.SaveFile(
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
