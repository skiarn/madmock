package main

import (
	"fmt"
	"log"
	"net/http"
)

// Pagehandler handles index page.
type Pagehandler struct {
	settings Settings
}

func (h *Pagehandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("Loading all recources..")
	title := "Available rescources:"
	resources := "<ul>"
	responseList, err := LoadAll(h.settings)
	if err != nil {
		for _, r := range *responseList {
			resources = resources + fmt.Sprintf("<div><b>URL:</b> %s <br> <b>Method:</b> %s <br> <b>ContentType:</b> %s </div> <br>", r.URL, r.Method, r.ContentType)
		}
	}
	resources = resources + "</ul>"
	content := fmt.Sprintf("<h1>%s</h1><div>%s</div>", title, resources)
	page := "<html>" + content + "</html>"
	fmt.Fprintf(w, "%s", page)
}
