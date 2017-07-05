package file

import (
	. "github.com/ahmetb/go-linq"
	"io/ioutil"
	"os"

	"fmt"
	"time"
)

type Root string
type BaseFilename string
type Contents string

func LoadDir(dir string) ([]string, error) {
	files, err := ioutil.ReadDir(dir)
	var filenames []string
	From(files).
		SelectT(func(c os.FileInfo) string { return c.Name() }).
		ToSlice(&filenames)
	return filenames, err
}

func SaveFile(root Root, baseFilename BaseFilename, contents Contents) error {
	now := time.Now().Unix()
	filename := fmt.Sprintf(
		"%s/%s-%d",
		root,
		baseFilename,
		now)
	return ioutil.WriteFile(filename, []byte(string(contents)), 0644)
}

// // List all files then get the matching one with the latest timestamp
// func LoadFileFn(msg fileCommandPayload) (string, error) {
// 	// list all files
// 	files, err := ioutil.ReadDir(string(serverRoot))
// 	if err != nil {
// 		return "", err
// 	}
// 	// get only the ones that match the filename and get the latest
// 	timestampRegexp := regexp.MustCompile("^.+-(.+)$")
// 	var filenames []string
// 	From(files).
// 		SelectT(func(c os.FileInfo) string {
// 			return c.Name()
// 		}).
// 		WhereT(func(c string) bool {
// 			nameRegexp := regexp.MustCompile("^" + msg.Filename + "-.+$")
// 			return nameRegexp.MatchString(c)
// 		}).
// 		OrderByDescendingT(func(name string) int {
// 			n, _ := strconv.Atoi(timestampRegexp.FindString(name))
// 			return n
// 		}).
// 		ToSlice(&filenames)
// 	if len(filenames) == 0 {
// 		return "", errors.New("file not found")
// 	}
// 	latestFilename := filenames[0]
// 	latestFullFilename := fmt.Sprintf("%s/%s", string(serverRoot), latestFilename)
// 	result, err := ioutil.ReadFile(latestFullFilename)
// 	return string(result[:]), err
// }
