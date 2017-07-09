package file

import (
	. "github.com/ahmetb/go-linq"
	"io/ioutil"
	"os"

	"errors"
	"fmt"
	"regexp"
	"strconv"
	"time"
)

type Root string
type BaseFilename string
type Contents string

type nowTime func() time.Time
type writeFile func(filename string, data []byte, perm os.FileMode) error
type readDir func(dirname string) ([]os.FileInfo, error)
type readFile func(filename string) ([]byte, error)

// get only the ones that match the filename and get the latest
var timestampRegexp = regexp.MustCompile("^.+-(.+)$")
var nameRegexp = regexp.MustCompile("^(.+)-.+$")

type FileIO struct {
	nowTime   nowTime
	readDir   readDir
	writeFile writeFile
	readFile  readFile
}

func (f FileIO) LoadDir(dir string) ([]string, error) {
	files, err := f.readDir(dir)
	if err != nil {
		return nil, err
	}
	var filenames []string
	From(files).
		SelectT(func(c os.FileInfo) string { return c.Name() }).
		GroupByT(
			func(name string) string {
				matches := nameRegexp.FindStringSubmatch(name)
				if len(matches) > 0 {
					return matches[len(matches)-1]
				}
				return name
			},
			func(name string) string { return name },
		).
		SelectT(func(g Group) string { return g.Key.(string) }).
		Distinct().
		SortT(func(a string, b string) bool { return a < b }).
		ToSlice(&filenames)
	return filenames, nil
}

func (f FileIO) SaveFile(
	root Root,
	baseFilename BaseFilename,
	contents Contents,
) error {
	now := f.nowTime().Unix()
	filename := fmt.Sprintf(
		"%s/%s-%d",
		root,
		baseFilename,
		now)
	return f.writeFile(filename, []byte(string(contents)), 0644)
}

// List all files then get the matching one with the latest timestamp
func (f FileIO) LoadFile(root Root, baseFilename BaseFilename) (string, error) {
	// list all files
	files, err := f.readDir(string(root))
	if err != nil {
		return "", err
	}
	var filenames []string
	From(files).
		SelectT(func(c os.FileInfo) string {
			return c.Name()
		}).
		WhereT(func(c string) bool {
			nameRegexp := regexp.MustCompile("^" + string(baseFilename) + "-.+$")
			return nameRegexp.MatchString(c)
		}).
		OrderByDescendingT(func(name string) int {
			// Sigh. We have to use this clunky thing
			matches := timestampRegexp.FindStringSubmatch(name)
			timestamp := matches[len(matches)-1]
			n, _ := strconv.Atoi(timestamp)
			return n
		}).
		ToSlice(&filenames)
	if len(filenames) == 0 {
		return "", errors.New("file not found")
	}
	latestFilename := filenames[0]
	latestFullFilename := fmt.Sprintf("%s/%s", root, latestFilename)
	result, err := f.readFile(latestFullFilename)
	return string(result[:]), err
}

func NewFileIO() FileIO {
	return FileIO{
		nowTime:   func() time.Time { return time.Now() },
		readDir:   ioutil.ReadDir,
		writeFile: ioutil.WriteFile,
		readFile:  ioutil.ReadFile,
	}
}
