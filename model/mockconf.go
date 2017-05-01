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

package model

import (
	"crypto/sha1"
	"encoding/base32"
	"errors"
	"fmt"
	"net/http"
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

//Valid returns true or false. Second argument is map containing errors. Key is field name and value is error message.
func (c MockConf) Valid()(bool, map[string]string){
	paramErrors := make(map[string]string)
	if strings.TrimSpace(c.URI) == "" {
		paramErrors["URI"] = "missing"
	}

	if strings.TrimSpace(c.Method) == "" {
		paramErrors["Method"] = "missing"
	}

	if strings.TrimSpace(c.ContentType) == "" {
		paramErrors["ContentType"] = "missing"
	}

	_, err := ValidateStatusCode(strconv.Itoa(c.StatusCode))
	if err != nil {
		paramErrors["StatusCode"] = err.Error()
	}
	return len(paramErrors) == 0, paramErrors
}

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

//GetFileName returns the filename for a MockConf enitiy.
func (c *MockConf) GetFileName() string {
	hasher := sha1.New()
	hasher.Write([]byte(c.Method + "-" + c.URI))
	filename := base32.StdEncoding.EncodeToString(hasher.Sum(nil))
	return filename
}

//GetMockFileName returns the base name of a mock.
func GetMockFileName(r *http.Request) string {
	hasher := sha1.New()
	hasher.Write([]byte(r.Method + "-" + r.URL.RequestURI()))
	filename := base32.StdEncoding.EncodeToString(hasher.Sum(nil))
	return filename
}

//ReadResponse from reading target response.
func ReadResponse(r *http.Response) *MockConf {
	defer r.Body.Close()

	rHeaders := make(map[string]string)
	for k, v := range r.Header {
		for _, vv := range v {
			rHeaders[k] = vv
		}
	}
	c := MockConf{URI: r.Request.URL.RequestURI(), Method: r.Request.Method, ContentType: r.Header.Get("Content-Type"), StatusCode: r.StatusCode, Header: rHeaders}
	return &c
}

func (c *MockConf) String() string {
	return fmt.Sprintf("URI:%s, Method:%s, Contenttype:%s, Status:%v\n", c.URI, c.Method, c.ContentType, c.StatusCode)
}
