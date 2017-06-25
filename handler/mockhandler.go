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
	"compress/gzip"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/skiarn/madmock/filesys"
	"github.com/skiarn/madmock/model"
	"github.com/skiarn/madmock/ws"
)

//HttpClient used for making requests to the target system to be mocked
type HttpClient interface {
	RequestTargetInfo(URL string, w http.ResponseWriter, r *http.Request) (*http.Response, error)
}

type client struct{}

func (client) RequestTargetInfo(targetURL string, w http.ResponseWriter, r *http.Request) (*http.Response, error) {
	client := &http.Client{}
	request, err := http.NewRequest(r.Method, targetURL, r.Body)
	if err != nil {
		log.Println("Failed to create request: ", request)
		return nil, err
	}
	copyHeader(r.Header, &request.Header)

	resp, err := client.Do(request)
	return resp, err
}

//Mockhandler handles request to be mocked.
type Mockhandler struct {
	TargetURL   string
	DataDirPath string
	Fs          filesys.FileSystem
	Client      HttpClient

	WSMockInfoHandler ws.Handler
}

//NewMockhandler handles initzialisation of NewMockhandler.
func NewMockhandler(targetURL string, dirpath string, wsh ws.Handler) Mockhandler {
	return Mockhandler{TargetURL: targetURL, DataDirPath: dirpath, Fs: filesys.LocalFileSystem{}, Client: client{}, WSMockInfoHandler: wsh}
}

func (h *Mockhandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mConfFileName := h.DataDirPath + "/" + model.GetMockFileName(r) + filesys.ConfEXT
	m, err := h.Fs.ReadMockConf(mConfFileName)
	if err != nil {
		if r.Method != "GET" {
			body := "Mock only executes GET requests. " + r.URL.String() + "\n" + r.Method + " is faked and can be edited in mock GUI. see /mock"
			m = &model.MockConf{URI: r.URL.RequestURI(), Method: r.Method, StatusCode: http.StatusOK}
			err := h.Fs.WriteMock(*m, []byte(body), h.DataDirPath)
			if err != nil {
				http.Error(w, "Error occured when requesting: "+r.URL.String()+" Failed with error: "+err.Error(), http.StatusInternalServerError)
			}
		} else {
			//do real request do targetet system, to retrive information from system.
			m, err = h.requestInfo(w, r)
			if err != nil {
				http.Error(w, "Resource unavailable: "+r.URL.String()+" Failed with error: "+err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}

	h.SendMockResponse(m, w, r)
}

//BuildTargetRequestURL replaces this mock service with real target server request url.
func (h *Mockhandler) BuildTargetRequestURL(r *http.Request) string {
	//A browser can issue a relative HTTP request, ex: GET / HTTP/1.1
	//When we on the server access it I like toallways include the url containing real adress ex:
	//GET http://localhost:8080/ HTTP/1.1

	u, err := url.Parse(r.URL.String())
	if err != nil {
		log.Fatal(err)
	}
	u.Host = h.TargetURL
	if u.Scheme == "" {
		u.Scheme = "http"
	}

	return u.String()
	//return strings.Replace(url, r.Host, h.TargetURL, 1)
}

func (h *Mockhandler) requestInfo(w http.ResponseWriter, r *http.Request) (*model.MockConf, error) {
	targetURL := h.BuildTargetRequestURL(r)

	response, err := h.Client.RequestTargetInfo(targetURL, w, r)
	if err != nil {
		return nil, err
	}
	var reader io.ReadCloser
	switch response.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(response.Body)
		if err != nil {
			return nil, err
		}
	default:
		reader = response.Body
	}
	defer reader.Close()

	defer response.Body.Close()
	contents, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	c := model.ReadResponse(response)

	err = h.Fs.WriteMock(*c, contents, h.DataDirPath)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func copyHeader(source http.Header, dest *http.Header) {
	for k, v := range source {
		for _, vv := range v {
			dest.Add(k, vv)
		}
	}
}

//SendMockResponse sending a response based on a mock.
func (h *Mockhandler) SendMockResponse(m *model.MockConf, w http.ResponseWriter, r *http.Request) {
	filename := h.DataDirPath + "/" + model.GetMockFileName(r) + filesys.ContentEXT
	dr, err := h.Fs.ReadResource(filename) //ioutil.ReadFile(filename)
	if err != nil {
		log.Printf("%s \n", err)
		http.Error(w, "Resource unavailable: "+r.URL.String()+" Failed with error: "+err.Error(), http.StatusServiceUnavailable)
		return
	}
	d, err := ioutil.ReadAll(dr)
	if err != nil {
		log.Printf("%s \n", err)
		http.Error(w, "Resource unavailable: "+r.URL.String()+" Failed with error: "+err.Error(), http.StatusServiceUnavailable)
		return
	}

	w.Header().Set("Content-Length", strconv.FormatInt(int64(len(d)), 10)) //strconv.FormatInt(dstat.Size(), 10))
	w.Header().Set("Content-Type", m.ContentType)
	w.WriteHeader(m.StatusCode)
	w.Write(d)
	log.Printf("%v \t %s \t %v bytes\n", m.StatusCode, r.URL.String(), len(d))

	wsmsg, err := json.Marshal(m)
	if err != nil {
		log.Printf("Websocket error: %s \n", err)
		return
	}
	msg := ws.Message{Data: wsmsg}
	h.WSMockInfoHandler.Send(&msg)
}
