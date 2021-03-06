package main

import (
	"encoding/json"
	"errors"
	"github.com/frankandrobot/acadia/file"
	"github.com/frankandrobot/acadia/messaging"
	"github.com/frankandrobot/acadia/queue"
	"github.com/labstack/echo"
	"io/ioutil"
	"net/http"
	"os"
)

// CONFIG

type ServerRoot string
type ServerKey string
type ServerCer string

type config struct {
	ServerRoot *ServerRoot `json:"serverRoot"`
	ServerKey  *ServerKey  `json:"serverKey"`
	ServerCert *ServerCer  `json:"serverCert"`
}

type contents struct {
	Contents string `form:"contents" json:"contents" binding:"required"`
}

func main() {
	rawConfig, err := ioutil.ReadFile(os.Getenv("HOME") + "/.acadia.json")
	if err != nil {
		panic(err)
	}
	var config config
	err = json.Unmarshal(rawConfig, &config)
	if err != nil {
		panic(err)
	}
	if config.ServerRoot == nil {
		panic(errors.New("Need a server root"))
	}
	if config.ServerKey == nil {
		panic(errors.New("Need a server key"))
	}
	if config.ServerCert == nil {
		panic(errors.New("Need a server certificate"))
	}
	router := echo.New()
	queue := queue.MakeQueue()
	fileIO := file.NewFileIO()

	router.GET("/files", func(c echo.Context) error {
		action := func() messaging.ChanResult {
			files, err := fileIO.LoadDir(string(*config.ServerRoot))
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
				file.Root(*config.ServerRoot),
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
				file.Root(*config.ServerRoot),
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

	err = router.StartTLS(":8080", string(*config.ServerCert), string(*config.ServerKey))
	if err != nil {
		panic(err)
	}
}
