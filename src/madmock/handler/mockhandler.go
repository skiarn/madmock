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
	"io/ioutil"
	"log"
	"madmock/filesys"
	"madmock/model"
	"net/http"
	"strconv"
)

//Mockhandler handles request to be mocked.
type Mockhandler struct {
	TargetURL   string
	DataDirPath string
	Fs          filesys.FileSystem
}

//NewMockhandler handles initzialisation of NewMockhandler.
func NewMockhandler(targetURL string, dirpath string) Mockhandler {
	return Mockhandler{TargetURL: targetURL, DataDirPath: dirpath, Fs: filesys.LocalFileSystem{}}
}

func (h *Mockhandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mConfFileName := h.DataDirPath + "/" + model.GetMockFileName(r) + filesys.ConfEXT
	log.Println("Trying to read conf file:", mConfFileName)
	m, err := h.Fs.ReadMockConf(mConfFileName)
	if err != nil {
		log.Println("Request kunde ej hittas försöker slå upp mot target", err)
		if r.Method != "GET" {
			errorMsg := "Could not execute request for: " + r.URL.String() + "\n" + r.Method + " should be executed manually using GUI."
			log.Println(errorMsg)
			w.Write([]byte(errorMsg))
			return
		}
		//do real request do targetet system, to retrive information from system.
		m, err = h.requestInfo(w, r)
		if err != nil {
			log.Printf("%s \n", err)
			http.Error(w, "Resource unavailable: "+r.URL.String()+" Failed with error: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	h.sendMockResponse(m, w, r)
}

func (h *Mockhandler) BuildTargetRequestURL(r *http.Request) (string, error) {
	//fmt.Println(r.URL.IsAbs())
	//if r.URL.IsAbs() {
	return "http://" + h.TargetURL + r.URL.Path, nil
	//}
	//return r.URL.Scheme + "://" + h.TargetURL + r.URL.Path, nil
}

func (h *Mockhandler) requestInfo(w http.ResponseWriter, r *http.Request) (*model.MockConf, error) {
	targetURL, err := h.BuildTargetRequestURL(r)
	if err != nil {
		return nil, err
	}

	log.Println("Fetching: " + targetURL)
	client := &http.Client{}

	request, err := http.NewRequest(r.Method, targetURL, r.Body)
	if err != nil {
		log.Println("Failed to create request: ", request)
		return nil, err
	}
	copyHeader(r.Header, &request.Header)

	response, err := client.Do(request)
	if err != nil {
		log.Println("Failed to make request: ", request)
		return nil, err
	}
	log.Println("Got response: ", response.Request.URL)
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	c := model.ReadResponse(response)

	log.Println("Created:", c)
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

func (h *Mockhandler) sendMockResponse(m *model.MockConf, w http.ResponseWriter, r *http.Request) {
	filename := h.DataDirPath + "/" + model.GetMockFileName(r) + filesys.ContentEXT
	log.Println("Trying to open:", filename)
	d, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Printf("%s \n", err)
		http.Error(w, "Resource unavailable: "+r.URL.String()+" Failed with error: "+err.Error(), http.StatusServiceUnavailable)
		return
	}

	w.Header().Set("Content-Length", strconv.FormatInt(int64(len(d)), 10)) //strconv.FormatInt(dstat.Size(), 10))
	w.Header().Set("Content-Type", m.ContentType)
	w.WriteHeader(m.StatusCode)
	w.Write(d)
	log.Printf("Writing status code:%v\n", m.StatusCode)
	log.Printf("%v bytes %s\n", len(d), r.URL.String())
}
