package file

import (
	. "github.com/ahmetb/go-linq"
	"io/ioutil"
	"os"

	"github.com/frankandrobot/acadia/messaging"
)

func LoadDir(dir string) messaging.Result {
	var result messaging.Result
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		result.Error = err
		return result
	}
	From(files).
		SelectT(func(c os.FileInfo) string { return c.Name() }).
		ToSlice(&result.Filenames)
	return result
}

// func SaveFileFn(msg fileCommandPayload) error {
// 	now := time.Now().Unix()
// 	filename := fmt.Sprintf(
// 		"%s/%s-%d",
// 		string(serverRoot),
// 		msg.Filename,
// 		now)
// 	return ioutil.WriteFile(filename, []byte(msg.Contents), 0644)
// }

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
