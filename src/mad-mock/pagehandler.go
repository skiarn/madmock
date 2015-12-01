package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

// Pagehandler handles index page.
type Pagehandler struct {
	settings Settings
}

func (h *Pagehandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("Loading all resources..")
	title := "Available resources:"
	resources := "<ul>"
	responseList, err := LoadAll(h.settings)
	if err != nil {
		log.Println("Error while loading resources: ", err)
	}
	for _, i := range *responseList {
		resources = resources + fmt.Sprintf(`<div><a href="%s">%s</a> <b>Method:</b> %s <br> <b>ContentType:</b> %s </div> <br>`, i.URI, r.URL.Host+i.URI, i.Method, i.ContentType)
	}
	resources = resources + "</ul>"
	content := fmt.Sprintf("<h1>%s</h1><div>%s</div>", title, resources)
	page := "<html>" + content + "</html>"
	fmt.Fprintf(w, "%s", page)
}

func (h *Pagehandler) handleMockconfPost(w http.ResponseWriter, r *http.Request) {
	paramErrors := make(map[string]string)
	uri := r.FormValue("URI")
	method := r.FormValue("Method")
	contentType := r.FormValue("ContentType")
	body := r.FormValue("body")
	if strings.TrimSpace(uri) == "" {
		paramErrors["URI"] = "missing"
	}

	//valid methods: GET,POST,PUT,DELETE
	validmethods := map[string]string{"GET": "GET", "POST": "POST", "PUT": "PUT", "DELETE": "DELETE"}
	if _, ok := validmethods[method]; !ok {
		paramErrors["Method"] = method + "is invalid."
	}

	if strings.TrimSpace(contentType) == "" {
		paramErrors["ContentType"] = "missing"
	}

	if len(paramErrors) != 0 {
		validationErrors, err := json.Marshal(paramErrors)
		if err != nil {
			log.Println(err)
		}
		http.Error(w, string(validationErrors), 400)
		return
	}

	c := MockConf{URI: uri, Method: method, ContentType: contentType}
	err := c.WriteToDisk([]byte(body), h.settings)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
