//Copyright (C) 2016  Andreas Westberg

//This program is free software: you can redistribute it and/or modify
//it under the terms of the GNU Lesser General Public License as published by
//the Free Software Foundation, version 3 of the License.

//This program is distributed in the hope that it will be useful,
//but WITHOUT ANY WARRANTY; without even the implied warranty of
//MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//GNU Lesser General Public License for more details.

//You should have received a copy of the GNU Lesser General Public License
//along with this program.  If not, see <http://www.gnu.org/licenses/lgpl-3.0.txt>.

package handler

import (
	"crypto/sha1"
	"encoding/base32"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

//MockConf represent a http call mock entity.
type MockConf struct {
	URI         string            `json:"uri"`
	Method      string            `json:"method"`
	ContentType string            `json:"contenttype"`
	StatusCode  int               `json:"status"`
	Header      map[string]string `json:"header"`

	Errors map[string]string
}

//MockConfs is a list of MockConf.
type MockConfs []MockConf

//ConfEXT is fileextension for config file.
const ConfEXT = ".mc"

//ContentEXT is fileextension for data body file.
const ContentEXT = ".data"

//ValidStatusCodes is valid http status codes.
var ValidStatusCodes = [...]int{100, 101, 200, 201, 202, 203, 204, 205, 206, 300, 301, 302, 303, 304, 305, 307, 400, 401, 402, 403, 404, 405, 406, 407, 408, 409, 410, 411, 412, 413, 414, 415, 416, 417, 418, 500, 501, 502, 503, 504, 505}

//ValidateStatusCode checks if string is valid http code.
func ValidateStatusCode(code string) (int, error) {
	if strings.TrimSpace(code) == "" {
		return 0, errors.New("missing")
	}
	statuscode, err := strconv.Atoi(code)
	if err != nil {
		return 0, err
	}

	isvalid := false
	for _, c := range ValidStatusCodes {
		if c == statuscode {
			isvalid = true
		}
	}

	if !isvalid {
		return 0, fmt.Errorf("%v is unknown", statuscode)
	}

	return statuscode, nil
}

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
