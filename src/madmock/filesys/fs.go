package filesys

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"madmock/model"
	"os"
	"path/filepath"
)

//ConfEXT is fileextension for config file.
const ConfEXT = ".mc"

//ContentEXT is fileextension for data body file.
const ContentEXT = ".data"

//FileSystem is extracted to make testing easier of the local filesystem.
type FileSystem interface {
	//Open(name string) (os.File, error)
	ReadMockConf(filepath string) (*model.MockConf, error)
	WriteMock(c model.MockConf, content []byte, dirpath string) error
	ReadAllMockConf(dataDirPath string) (*model.MockConfs, error)
	ReadResource(filepath string) (io.Reader, error)
}

//LocalFileSystem is implementation of application filesystem.
type LocalFileSystem struct{}

//func (LocalFileSystem) Open(name string) (*os.File, error) { return os.Open(name) }

//func (LocalFileSystem) Remove(name string) error       { return os.Remove(name) }

//ReadMockConf returns a MockConf model after reading file from disk.
func (LocalFileSystem) ReadMockConf(filepath string) (*model.MockConf, error) {
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	var m model.MockConf
	json.Unmarshal(data, &m)
	return &m, nil
}

//WriteMock writes a MockConf to disk with conf and content in separate files.
func (LocalFileSystem) WriteMock(c model.MockConf, content []byte, dirpath string) error {
	b, err := json.Marshal(c)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(dirpath+"/"+c.GetFileName()+ContentEXT, content, 0644)
	err = ioutil.WriteFile(dirpath+"/"+c.GetFileName()+ConfEXT, b, 0644)
	log.Println("Finished writing mock:", dirpath+"/"+c.GetFileName())
	return err

}

//ReadAllMockConf reads all MockConf entities in dirpath.
func (fs LocalFileSystem) ReadAllMockConf(dirpath string) (*model.MockConfs, error) {
	var confs model.MockConfs

	d, err := os.Open(dirpath)
	if err != nil {
		return nil, err
	}
	defer d.Close()

	files, err := d.Readdir(-1)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if file.Mode().IsRegular() {
			if filepath.Ext(file.Name()) == ConfEXT {
				fmt.Println("Found: " + file.Name())
				data, err := ioutil.ReadFile(dirpath + "/" + file.Name())
				if err != nil {
					return nil, err
				}
				var c model.MockConf
				json.Unmarshal(data, &c)
				confs = append(confs, c)
			}
		}
	}
	return &confs, nil
}

//ReadResource returns a reader to the resource.
func (LocalFileSystem) ReadResource(filepath string) (io.Reader, error) {
	log.Println("Trying to read resource: ", filepath)
	d, err := os.Open(filepath)
	if err != nil {
		log.Printf("Failed to lead resource: %s", filepath)
		return nil, err
	}
	//defer d.Close()
	return d, nil
}
