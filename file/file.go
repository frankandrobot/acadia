package file

import (
	. "github.com/ahmetb/go-linq"
	"io/ioutil"
	"os"

	"errors"
	"fmt"
	"regexp"
	"strconv"
)

type Root string
type BaseFilename string
type Contents string

type writeFile func(filename string, data []byte, perm os.FileMode) error
type readDir func(dirname string) ([]os.FileInfo, error)
type readFile func(filename string) ([]byte, error)

var versionRegex = regexp.MustCompile("^.+-(.+)$") // gets the version
var nameRegex = regexp.MustCompile("^(.+)-.+$")    // gets the filename

// baseFilename gets the baseFilename from a versionedFilename.
// Ex: "foobar-0" => "foobar"
func getBaseFilename(versionedFilename string) (string, bool) {
	matches := nameRegex.FindStringSubmatch(versionedFilename)
	if len(matches) == 0 {
		return "", false
	}
	return matches[len(matches)-1], true
}

// getFileVersion gets the file version from a versionedFilename.
// Ex: "foobar-0" => 0
func getFileVersion(versionedFilename string) (int, bool) {
	// Sigh. We have to use this clunky thing
	matches := versionRegex.FindStringSubmatch(versionedFilename)
	if len(matches) == 0 {
		return 0, false
	}
	version := matches[len(matches)-1]
	n, err := strconv.Atoi(version)
	if err != nil {
		return 0, false
	}
	return n, true
}

type FileIO struct {
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
		// group filenames by base filename
		GroupByT(
			func(versionedFilename string) string {
				baseFilename, ok := getBaseFilename(versionedFilename)
				if ok {
					return baseFilename
				}
				// if !ok it's actually not a versioned filename
				return versionedFilename
			},
			func(versionedFilename string) string {
				return versionedFilename
			},
		).
		// then get just the distinct group names
		SelectT(func(g Group) string { return g.Key.(string) }).
		Distinct().
		// then alpha sort
		SortT(func(a string, b string) bool { return a < b }).
		ToSlice(&filenames)
	return filenames, nil
}

func (f FileIO) SaveFile(
	root Root,
	baseFilename BaseFilename,
	contents Contents,
) error {
	// get the latest version of the baseFilename
	files, err := f.readDir(string(root))
	if err != nil {
		return err
	}
	var filenames []string
	From(files).
		SelectT(func(c os.FileInfo) string { return c.Name() }).
		// get all the filenames that match the baseFilename
		WhereT(func(versionedFilename string) bool {
			thisname, ok := getBaseFilename(versionedFilename)
			if ok {
				return thisname == string(baseFilename)
			}
			return false
		}).
		// then order by version
		OrderByDescendingT(func(versionedFilename string) int {
			version, ok := getFileVersion(versionedFilename)
			if ok {
				return version
			}
			return 0
		}).
		ToSlice(&filenames)
	latestVersion := 0
	if len(filenames) > 0 {
		latest := filenames[0]
		var ok bool // look ma, := definition leads to a bug!
		latestVersion, ok = getFileVersion(latest)
		if ok {
			latestVersion = latestVersion + 1
		}
	}
	filename := fmt.Sprintf(
		"%s/%s-%d",
		root,
		baseFilename,
		latestVersion)
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
		SelectT(func(c os.FileInfo) string { return c.Name() }).
		// get all filenames that match the baseFilename
		WhereT(func(versionedFilename string) bool {
			thisname, ok := getBaseFilename(versionedFilename)
			if ok {
				return thisname == string(baseFilename)
			}
			return false
		}).
		// then sort by version
		OrderByDescendingT(func(versionedFilename string) int {
			n, _ := getFileVersion(versionedFilename)
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
		readDir:   ioutil.ReadDir,
		writeFile: ioutil.WriteFile,
		readFile:  ioutil.ReadFile,
	}
}
