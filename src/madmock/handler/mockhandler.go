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
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

//Mockhandler handles request to be mocked.
type Mockhandler struct {
	TargetURL   string
	DataDirPath string
}

func (h *Mockhandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("Request: " + r.URL.String())
	m, err := Load(r, h.DataDirPath)
	if err != nil {
		log.Println(err)
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

func (h *Mockhandler) requestInfo(w http.ResponseWriter, r *http.Request) (*MockConf, error) {
	requestURL, err := GetRequestURL(r.RequestURI, h.TargetURL)
	if err != nil {
		return nil, err
	}

	log.Println("Fetching: " + requestURL)
	client := &http.Client{}

	if r.Method != "GET" {
		errorMsg := "Could not execute request for: " + requestURL + "\n" + r.Method + " should be executed manually using GUI."
		log.Println(errorMsg)
		return nil, errors.New(errorMsg)
	}

	request, err := http.NewRequest(r.Method, requestURL, r.Body)
	if err != nil {
		return nil, err
	}
	copyHeader(r.Header, &request.Header)

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	responseHeaders := make(map[string]string)
	for k, v := range response.Header {
		for _, vv := range v {
			responseHeaders[k] = vv
		}
	}
	c := MockConf{URI: r.RequestURI, Method: r.Method, ContentType: response.Header.Get("Content-Type"), StatusCode: response.StatusCode, Header: responseHeaders}
	err = c.WriteToDisk(contents, h.DataDirPath)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

func copyHeader(source http.Header, dest *http.Header) {
	for k, v := range source {
		for _, vv := range v {
			dest.Add(k, vv)
		}
	}
}

func (h *Mockhandler) sendMockResponse(m *MockConf, w http.ResponseWriter, r *http.Request) {
	d, err := os.Open(h.DataDirPath + "/" + m.GetFileName() + ".data")
	if err != nil {
		log.Printf("%s \n", err)
		http.Error(w, "Resource unavailable: "+r.URL.String()+" Failed with error: "+err.Error(), http.StatusServiceUnavailable)
		return
	}
	defer d.Close()
	dstat, err := d.Stat()
	if err != nil {
		http.Error(w, "Resource unavailable: "+r.URL.String()+" Failed with error: "+err.Error(), http.StatusServiceUnavailable)
		return
	}
	w.Header().Set("Content-Length", strconv.FormatInt(dstat.Size(), 10))
	w.Header().Set("Content-Type", m.ContentType)
	w.WriteHeader(m.StatusCode)
	n, err := io.Copy(w, d)
	if err != nil {
		http.Error(w, "Internal error while wringing response: "+r.URL.String()+" Failed with error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("%v bytes %s\n", n, r.URL.String())
}
