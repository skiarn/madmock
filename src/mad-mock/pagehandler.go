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
	postform := `<h4>New</h4>
					<form method="POST" action="/mock/add">
							<label>URI:
								<input type="text" name="URI" placeholder="/example/123"/>
							</label>
							<select name="Method">
							  <option value="GET">GET</option>
							  <option value="POST">POST</option>
							  <option value="PUT">PUT</option>
							  <option value="DELETE">DELETE</option>
							</select>
							<label>ContentType:<input type="text" name="ContentType" placeholder="application/json; charset=UTF-8"/></label>
							<br>
							<textarea type="text" name="body" placeholder="Response..."></textarea>
    				<input type="submit">
					</form>`

	for _, i := range *responseList {
		resources = resources + fmt.Sprintf(`<div><a href="%s">%s</a> <b>Method:</b> %s <br> <b>ContentType:</b> %s </div> <br>`, i.URI, r.URL.Host+i.URI, i.Method, i.ContentType)
	}
	css := "<style>" + `.container { width: 100%; border: 1px solid #cccccc; }
											textarea {
												width: 100%;
												display: block;
												border: 1px solid #cccccc;
												padding: 1px;
												height: 20vw;
												margin: 0; padding: 0;
												border-width: 0;
											}` + "</style>"
	resources = resources + "</ul>"
	content := fmt.Sprintf(`<h1>%s</h1><div class="container">%s</div><div class="container">%s</div>`, title, resources, postform)
	page := "<html><head>" + css + "</head>" + "<body>" + content + "</body></html>"
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
	http.Redirect(w, r, "/mock", http.StatusFound)
}
