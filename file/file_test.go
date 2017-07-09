package file

import "testing"
import "os"
import "time"
import "reflect"

type fileInfo struct {
	name string
}

func (f fileInfo) Name() string       { return f.name }
func (f fileInfo) Size() int64        { return int64(0) }
func (f fileInfo) Mode() os.FileMode  { return os.FileMode(uint32(0)) }
func (f fileInfo) ModTime() time.Time { return time.Now() }
func (f fileInfo) IsDir() bool        { return true }
func (f fileInfo) Sys() interface{}   { return nil }

func TestLoadDir(t *testing.T) {
	cases := []struct {
		message       string
		fileIO        FileIO
		expected      []string
		expectedError error
	}{
		{
			message: "basic case",
			fileIO: FileIO{
				readDir: func(dirname string) ([]os.FileInfo, error) {
					if dirname != "root" {
						t.Errorf("LoadFile wanted root got %s", dirname)
					}
					return []os.FileInfo{
							fileInfo{name: "file-1"},
							fileInfo{name: "file-3"},
							fileInfo{name: "file-2"},
							fileInfo{name: "anotherfile-1"},
						},
						nil
				},
			},
			expectedError: nil,
			expected:      []string{"anotherfile", "file"},
		},
		{
			message: "edge case",
			fileIO: FileIO{
				readDir: func(dirname string) ([]os.FileInfo, error) {
					if dirname != "root" {
						t.Errorf("LoadDir wanted root got %s", dirname)
					}
					return []os.FileInfo{
							fileInfo{name: "edgecase"},
						},
						nil
				},
			},
			expectedError: nil,
			expected:      []string{"edgecase"},
		},
	}
	for _, c := range cases {
		results, err := c.fileIO.LoadDir("root")
		if err != c.expectedError {
			t.Errorf("LoadDir %s: got %s instead of %s",
				c.message, c.expectedError, err)
		}
		if !reflect.DeepEqual(results, c.expected) {
			t.Errorf("LoadDir %s: got %s instead of %s",
				c.message, c.expected, results)
		}
	}
}

func TestLoadFile(t *testing.T) {
	cases := []struct {
		message       string
		input         string
		expectedError error
		fileIO        FileIO
	}{
		{
			message:       "Basic case",
			input:         "file",
			expectedError: nil,
			fileIO: FileIO{
				readDir: func(dirname string) ([]os.FileInfo, error) {
					if dirname != "root" {
						t.Errorf("LoadFile wanted root got %s", dirname)
					}
					return []os.FileInfo{
							fileInfo{name: "file-1"},
							fileInfo{name: "file-3"},
							fileInfo{name: "file-2"},
							fileInfo{name: "anotherfile-1"},
						},
						nil
				},
				readFile: func(filename string) ([]byte, error) {
					if filename != "root/file-3" {
						t.Errorf("LoadFile wanted file got %s", filename)
					}
					return []byte("Hello world!"), nil
				},
			},
		},
		{
			message:       "works with wierd names",
			input:         "file-1",
			expectedError: nil,
			fileIO: FileIO{
				readDir: func(dirname string) ([]os.FileInfo, error) {
					if dirname != "root" {
						t.Errorf("LoadFile wanted root got %s", dirname)
					}
					return []os.FileInfo{
							fileInfo{name: "file-1-5"},
							fileInfo{name: "file-1-2"},
							fileInfo{name: "file-1-1"},
							fileInfo{name: "anotherfile-1"},
						},
						nil
				},
				readFile: func(filename string) ([]byte, error) {
					if filename != "root/file-1-5" {
						t.Errorf("LoadFile got %s", filename)
					}
					return []byte("Hello world!"), nil
				},
			},
		},
	}
	for _, c := range cases {
		contents, err := c.fileIO.LoadFile(Root("root"), BaseFilename(c.input))
		if err != c.expectedError {
			t.Errorf("LoadFile %s: got %s instead of %s",
				c.message, err, c.expectedError)
		}
		if string(contents) != "Hello world!" {
			t.Errorf("LoadFile %s: expected Hello world!, got %s",
				c.message, contents)
		}
	}
}
