package filesys_test

import (
	"errors"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/skiarn/madmock/filesys"
)

//func (f *File) Read(b []byte) (n int, err error) { return }

//func (fMock) Read(name string) (filesys.File, error) { return bytes.NewReader([]byte("herrow")), nil }

//func (fMock) Close() error                           { return nil }

type fileModeMock struct{}

func (fileModeMock) IsDir() bool     { return true }
func (fileModeMock) IsRegular() bool { return true }

//func (fileModeMock) Perm() os.FileMode { return fileModeMock }
func (fileModeMock) String() string { return "" }

type fileInfo1 struct{}

func (fileInfo1) Name() string { return "filenameone" }
func (fileInfo1) Size() int64  { return 1 }
func (fileInfo1) Mode() os.FileMode {
	var fm os.FileMode = fileModeMock{}
	return fm
}
func (fileInfo1) ModTime() time.Time { return time.Now() }
func (fileInfo1) IsDir() bool {
	return fileInfo1.Mode().IsDir()
}
func (fileInfo1) Sys() interface{} { return "" }

//type fileInfo2 struct{}

//func (fileInfo2) Name() string { return "filenametwo" }

type fileMock1 struct{}

func (fileMock1) Close() (err error) { return nil }
func (fileMock1) Read([]byte) (int, error) {
	return len([]byte("FILE INFORMATION DATA")), nil
}
func (fileMock1) ReadAt([]byte, int64) (int, error) { return len([]byte("FILE INFORMATION DATA")), nil }
func (fileMock1) Readdir(i int) ([]os.FileInfo, error) {
	if i != -1 {
		panic("Expected -1 and not:" + strconv.Itoa(i))
	}
	var fileinfos []os.FileInfo
	fileinfos = append(fileinfos, fileInfo1{})
	//fileinfos = append(fileinfos, fileInfo2{})

	return fileinfos, nil
}

//bytes.NewReader([]byte("FILE INFORMATION DATA")
//func (fMock) Read(b []byte) (n int, err error)
//fsMockImpl is implementation of application mocked filesystem.
type fsMockImpl struct{}

func (fsMockImpl) Open(name string) (filesys.File, error) {
	//return nil, errors.New("File not mocked!")

	var f filesys.File = fileMock1{}

	return f, nil
}
func (fsMockImpl) Remove(name string) error { return errors.New("Could not remove mocked file.") }

func TestGetMockConfFilepaths(t *testing.T) {
	expectedFilePaths := []string{"fileone.mc", "filetwo.mc"}

	fileinfos, err := filesys.GetFileInfoWithExtension(fsMockImpl{}, "path/to/files", ".mc")
	if err != nil {
		t.Errorf("Error occured:%v", err)
	}

	for _, f := range expectedFilePaths {
		if !contains(fileinfos, f) {
			t.Errorf("Expected to find: %v in: %v", f, fileinfos)
		}
	}
}

func contains(s []os.FileInfo, expectedfilename string) bool {
	for _, a := range s {
		if a.Name() == expectedfilename {
			return true
		}
	}
	return false
}
