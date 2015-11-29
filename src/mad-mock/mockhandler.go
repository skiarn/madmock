package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

type Mockhandler struct {
	settings Settings
}

func (h *Mockhandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("Trying to mock request: " + r.URL.String())
	m, err := Load(r, h.settings)
	if err != nil {
		fmt.Printf("%s \n", err)
		http.Error(w, "Resource unavailable: "+r.URL.String()+" Failed with error: "+err.Error(), http.StatusServiceUnavailable)
		return
	}

	d, err := os.Open(h.settings.DataDirPath + "/" + m.GetFileName() + ".data")
	if err != nil {
		log.Printf("%s \n", err)
		http.Error(w, "Resource unavailable: "+r.URL.String()+" Failed with error: "+err.Error(), http.StatusServiceUnavailable)
		return
	}
	defer d.Close()
	dstat, err := d.Stat()
	if err != nil {
		http.Error(w, "Resource unavailable: "+r.URL.String()+" Failed with error: "+err.Error(), http.StatusServiceUnavailable)
	}
	w.Header().Set("Content-Length", strconv.FormatInt(dstat.Size(), 10))
	w.Header().Set("Content-Type", m.ContentType)
	n, err := io.Copy(w, d)
	if err != nil {
		http.Error(w, "Internal error while wringing response: "+r.URL.String()+" Failed with error: "+err.Error(), http.StatusInternalServerError)
	}
	log.Printf("Copied %v bytes\n", n)
}
