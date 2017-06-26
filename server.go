package main

import "net/http"
import "io/ioutil"
import "os"
import "fmt"
import "github.com/gin-gonic/gin"
import . "github.com/ahmetb/go-linq"
import "time"
import "regexp"
import "strconv"
import "errors"

// CONFIG

type ServerRoot string
type ServerKey string
type ServerCer string

var serverRoot = ServerRoot("/media/usb/data")
var serverKey = ServerKey("/home/pi/server.key")
var serverCert = ServerCer("/home/pi/server.crt")

// FILE COMMANDS

type fileCommand int

const (
	saveFile fileCommand = iota
	loadFile
)

type HandleResultFn func(result fileCommandResult)

type fileCommandMsg struct {
	Filename string
	Contents string // obviously only makes sense in saving
}

type fileCommandPayload struct {
	fileCommandMsg
	Command      fileCommand
	HandleResult HandleResultFn
}

type fileCommandResult struct {
	fileCommandMsg
	Error error
}

func fileManagerQueue() chan fileCommandPayload {
	queue := make(chan fileCommandPayload)
	go func() {
		for {
			msg := <-queue
			var result fileCommandResult
			switch msg.Command {
			case saveFile:
				result = fileCommandResult{
					Error: saveFileFn(msg),
				}
			case loadFile:
				contents, err := loadFileFn(msg)
				result := fileCommandResult{}
				result.Contents = contents
				result.Error = err
			}
			msg.HandleResult(result)
		}
	}()
	return queue
}

func saveFileFn(msg fileCommandPayload) error {
	now := time.Now().Unix()
	filename := fmt.Sprintf(
		"%s/%s-%d",
		string(serverRoot),
		msg.Filename,
		now)
	return ioutil.WriteFile(filename, []byte(msg.Contents), 0644)
}

// List all files then get the matching one with the latest timestamp
func loadFileFn(msg fileCommandPayload) (string, error) {
	// list all files
	files, err := ioutil.ReadDir(string(serverRoot))
	if err != nil {
		return "", err
	}
	// get only the ones that match the filename and get the latest
	timestampRegexp := regexp.MustCompile("^.+-(.+)$")
	var filenames []string
	From(files).
		SelectT(func(c os.FileInfo) string {
			return c.Name()
		}).
		WhereT(func(c string) bool {
			nameRegexp := regexp.MustCompile("^" + msg.Filename + "-.+$")
			return nameRegexp.MatchString(c)
		}).
		OrderByDescendingT(func(name string) int {
			n, _ := strconv.Atoi(timestampRegexp.FindString(name))
			return n
		}).
		ToSlice(&filenames)
	if len(filenames) == 0 {
		return "", errors.New("file not found")
	}
	latestFilename := filenames[0]
	latestFullFilename := fmt.Sprintf("%s/%s", string(serverRoot), latestFilename)
	result, err := ioutil.ReadFile(latestFullFilename)
	return string(result[:]), err
}

func main() {
	router := gin.Default()
	queue := fileManagerQueue()

	router.GET("/files", func(c *gin.Context) {
		files, err := ioutil.ReadDir(string(serverRoot))
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

	router.GET("/files/:name", func(c *gin.Context) {
		name := c.Param("name")
		filename := fmt.Sprintf("%s/%s", string(serverRoot), name)
		payload := fileCommandPayload{}
		payload.Filename = filename
		payload.HandleResult = func(result fileCommandResult) {
			if result.Error == nil {
				c.String(http.StatusOK, "%s", result.Contents)
			} else {
				c.String(http.StatusNotFound, "")
			}
		}
		queue <- payload
	})

	router.POST("/files/:name", func(c *gin.Context) {
		name := c.Param("name")
		contents, exists := c.GetPostForm("contents")
		if !exists {
			c.String(http.StatusBadRequest, "Need contents")
			return
		}
		payload := fileCommandPayload{}
		payload.Command = saveFile
		payload.Filename = name
		payload.Contents = contents
		payload.HandleResult = func(result fileCommandResult) {
			if result.Error == nil {
				c.String(http.StatusOK, "")
			} else {
				c.String(http.StatusInternalServerError, "%s", result.Error)
			}
		}
		queue <- payload
	})

	router.RunTLS(":8080", string(serverCert), string(serverKey))
}
