package main

import (
	"crypto/sha1"
	"encoding/base32"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

type MockConf struct {
	URL         string `json:"url"`
	Method      string `json:"method"`
	ContentType string `json:"contenttype"`

	Errors map[string]string
}

type MockConfs []MockConf

func (c *MockConf) GetFileName() string {
	hasher := sha1.New()
	hasher.Write([]byte(c.URL))
	filename := base32.StdEncoding.EncodeToString(hasher.Sum(nil))
	return filename
}
func Load(r *http.Request, settings Settings) (*MockConf, error) {
	url := r.URL.String()
	hasher := sha1.New()
	hasher.Write([]byte(url))
	filename := base32.StdEncoding.EncodeToString(hasher.Sum(nil))
	dir := settings.DataDirPath
	data, err := ioutil.ReadFile(dir + "/" + filename + ".mc")
	if err != nil {
		return nil, err
	}
	var m MockConf
	json.Unmarshal(data, &m)
	return &m, nil

}

func LoadAll(settings Settings) (*MockConfs, error) {
	var confs MockConfs

	d, err := os.Open(settings.DataDirPath)
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
			if filepath.Ext(file.Name()) == ".mc" {
				fmt.Println("Found: " + file.Name())
				data, err := ioutil.ReadFile(settings.DataDirPath + "/" + file.Name())
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
