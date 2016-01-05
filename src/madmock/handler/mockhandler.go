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
	//TODO: copy all request headers?.
	request.Header.Set("Content-Type", r.Header.Get("Content-Type"))
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	c := MockConf{URI: r.RequestURI, Method: r.Method, ContentType: response.Header.Get("Content-Type")}
	err = c.WriteToDisk(contents, h.DataDirPath)
	if err != nil {
		return nil, err
	}

	return &c, nil
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
	n, err := io.Copy(w, d)
	if err != nil {
		http.Error(w, "Internal error while wringing response: "+r.URL.String()+" Failed with error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("%v bytes %s\n", n, r.URL.String())
}
