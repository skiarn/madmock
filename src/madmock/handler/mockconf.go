package handler

import (
	"crypto/sha1"
	"encoding/base32"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

//MockConf represent a http call mock entity.
type MockConf struct {
	URI         string `json:"uri"`
	Method      string `json:"method"`
	ContentType string `json:"contenttype"`

	Errors map[string]string
}

//MockConfs is a list of MockConf.
type MockConfs []MockConf

const ConfEXT = ".mc"
const ContentEXT = ".data"

//WriteToDisk saves a MockConf to disk.
func (c MockConf) WriteToDisk(content []byte, dataDirPath string) error {
	b, err := json.Marshal(c)
	if err != nil {
		return err
	}

	ioutil.WriteFile(dataDirPath+"/"+c.GetFileName()+ContentEXT, content, 0644)
	return ioutil.WriteFile(dataDirPath+"/"+c.GetFileName()+ConfEXT, b, 0644)

}

//GetFileName returns the filename for a MockConf enitiy.
func (c *MockConf) GetFileName() string {

	hasher := sha1.New()
	hasher.Write([]byte(c.Method + "-" + c.URI))
	filename := base32.StdEncoding.EncodeToString(hasher.Sum(nil))
	return filename
}

func GetFileName(r *http.Request) (string, error) {
	hasher := sha1.New()
	hasher.Write([]byte(r.Method + "-" + r.RequestURI))
	filename := base32.StdEncoding.EncodeToString(hasher.Sum(nil))
	return filename, nil
}

//GetRequestURL builds the request target url.
func GetRequestURL(uri string, targetURL string) (string, error) {
	target, err := url.Parse(targetURL)
	if err != nil {
		return "", err
	}
	target.Scheme = "http"
	return target.String() + uri, nil
}

//Load tries to read a MockConf from disk by using the request url to determine the filename.
func Load(r *http.Request, dataDirPath string) (*MockConf, error) {

	filename, err := GetFileName(r)
	if err != nil {
		return nil, err
	}
	dir := dataDirPath
	data, err := ioutil.ReadFile(dir + "/" + filename + ConfEXT)
	if err != nil {
		return nil, err
	}
	var m MockConf
	json.Unmarshal(data, &m)
	return &m, nil

}

//LoadAll loads all MockConf entities available.
func LoadAll(dataDirPath string) (*MockConfs, error) {
	var confs MockConfs

	d, err := os.Open(dataDirPath)
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
				data, err := ioutil.ReadFile(dataDirPath + "/" + file.Name())
				if err != nil {
					return nil, err
				}
				var c MockConf
				json.Unmarshal(data, &c)
				confs = append(confs, c)
			}
		}
	}
	return &confs, nil
}
