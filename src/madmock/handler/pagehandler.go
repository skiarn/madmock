package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
)

// Pagehandler handles index page.
type Pagehandler struct {
	DataDirPath string
}

const title = "Mad Mock"
const css = "<style>" + `
									html,
									body {
										height: 100%;
										background-color: #3F51B5;
									}
									body {
										color: #fff;
										text-align: center;
										text-shadow: 0 1px 3px rgba(0,0,0,.5);
									}
									.container .content {
											display: none;
											padding : 5px;
									}
									//margin-top: 3px; margin-bottom: 3px; margin-right: 3px; margin-left: 3px;
									form {text-align: left; text-shadow: none; color: #000000; }
									.container {text-align: left; text-shadow: none; color: #000000; width: 100%; }
										textarea {
											width: 100%;
											display: block;
											border: 1px solid #cccccc;
											padding: 1px;
											height: 20vw;
											margin: 0; padding: 0;
											border-width: 0;
										}
										.delete-btn {
											border:1px solid transparent;
											background-color: transparent;
											display: inline-block;
											vertical-align: middle;
											outline: 0;
											cursor: pointer;
										}
										.delete-btn:after {
											content: "X";
											display: block;
											width: 20px;
											height: 20px;
											background-color: #C62828;
											right: 35px;
											top: 0;
											bottom: 0;
											margin: auto;
											padding: 2px;
											border-radius: 50%;
											text-align: center;
											color: white;
											font-weight: normal;
											font-size: 15px;
											box-shadow: 0 0 2px #E57373;
											cursor: pointer;
										}` + "</style>"

//statusCodeHTMLSelect returns html selector for http status code.
func statusCodeHTMLSelect() string {
	var statuscodehtml = `<select id="statuscodeS" name="StatusCode">`
	for _, c := range ValidStatusCodes {
		statuscodehtml += fmt.Sprintf("<option value=\"%v\">%v, %s</option>", c, c, http.StatusText(c))
	}
	return statuscodehtml + `</select>`
}
func saveformHTML() string {
	saveform := `<form method="POST" action="/mock/save/">`
	statuscodeSelect := statusCodeHTMLSelect()
	r := `<label>URI:
			<input id="uriI" style="width:80%;" type="text" name="URI" placeholder="/example/123"/>
		</label>
		<select id="methodS" name="Method">
			<option value="GET">GET</option>
			<option value="POST">POST</option>
			<option value="PUT">PUT</option>
			<option value="DELETE">DELETE</option>
		</select>
		<br>
		<label>ContentType:
			<input id="contenttypeI" type="text" name="ContentType" placeholder="application/json; charset=UTF-8"/ style="width:40%;"></label>
		<br>
		<textarea id="contentTA" type="text" name="body" placeholder="Response..."></textarea>
	<input type="submit" value="Submit">
</form>`
	return saveform + statuscodeSelect + r
}
func (h *Pagehandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.servePage(w, r)
}

//ViewDataHandler load resource data.
func (h *Pagehandler) ViewDataHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Path[len("mock/view/data/"):]

	d, err := os.Open(h.DataDirPath + "/" + name + ContentEXT)
	if err != nil {
		http.Error(w, "Was not able to find resource: "+r.URL.String()+" Failed with error: "+err.Error(), http.StatusBadRequest)
		return
	}
	_, err = io.Copy(w, d)
	if err != nil {
		http.Error(w, "Internal error while wringing response: "+r.URL.String()+" Failed with error: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

//ViewConfHandler loads resource configuration.
func (h *Pagehandler) ViewConfHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Path[len("mock/view/data/"):]

	d, err := os.Open(h.DataDirPath + "/" + name + ConfEXT)
	if err != nil {
		http.Error(w, "Was not able to find resource: "+r.URL.String()+" Failed with error: "+err.Error(), http.StatusBadRequest)
		return
	}
	_, err = io.Copy(w, d)
	if err != nil {
		http.Error(w, "Internal error while wringing response: "+r.URL.String()+" Failed with error: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

//EditHandler handles resource page to edit item.
func (h *Pagehandler) EditHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Path[len("/mock/edit/"):]
	editform := "<div class=\"container\"> <h4>Editing</h4>" + "<h5>Headers</h5><p id=\"headersP\"></p> <br>" + saveformHTML() + "</div>"
	content := fmt.Sprintf(`<h1>%s</h1><div>%s</div>`, title, editform)
	script := fmt.Sprintf(`<script>
		function getResource() {
				//ajaxcall get resource data.
				var xmlhttp;
		    // compatible with IE7+, Firefox, Chrome, Opera, Safari
		    datahttp = new XMLHttpRequest();
		    datahttp.onreadystatechange = function(){
		        if (datahttp.readyState == 4 && datahttp.status == 200){
		            console.log(datahttp)
								document.getElementById('contentTA').value = datahttp.response;
		        }
		    }
		    datahttp.open("GET", '/mock/view/data/%s', true);
		    datahttp.send();

				//ajaxcall get resource configuration.
				var confhttp;
		    // compatible with IE7+, Firefox, Chrome, Opera, Safari
		    confhttp = new XMLHttpRequest();
		    confhttp.onreadystatechange = function(){
		        if (confhttp.readyState == 4 && confhttp.status == 200){
								var json = JSON.parse(confhttp.response);
		            console.log(json)
								console.log("uri:" + json.uri)
								console.log("contenttype:" + json.contenttype)
								console.log("method:" + json.method)
								document.getElementById('uriI').value = json.uri;
								document.getElementById('contenttypeI').value = json.contenttype;
								document.getElementById('methodS').value = json.method;
								document.getElementById('statuscodeS').value = json.status;

								document.getElementById('headersP').innerHTML = JSON.stringify(json.header);
		        }
		    }
		    confhttp.open("GET", '/mock/view/conf/%s', true);
		    confhttp.send();


		}
	</script>`, name, name)

	page := "<html><head>" + script + css + "</head>" + `<body onload="getResource()">` + content + "</body></html>"
	fmt.Fprintf(w, "%s", page)

}

//New serves a page for creating a new mock entity.
func (h *Pagehandler) New(w http.ResponseWriter, r *http.Request) {
	postform := "<div class=\"container\"> <h4>New</h4>" + saveformHTML() + "</div>"
	content := fmt.Sprintf(`<h1>%s</h1><div>%s</div>`, title, postform)

	page := "<html><head>" + css + "</head>" + "<body>" + content + "</body></html>"
	fmt.Fprintf(w, "%s", page)
}

//Save a new mock entity on serverside. Will create new resource if not existing and if it already exist it will update the resource.
func (h *Pagehandler) Save(w http.ResponseWriter, r *http.Request) {
	paramErrors := make(map[string]string)
	uri := r.FormValue("URI")
	method := r.FormValue("Method")
	contentType := r.FormValue("ContentType")
	body := r.FormValue("body")
	statuscodeFV := r.FormValue("StatusCode")
	if strings.TrimSpace(uri) == "" {
		paramErrors["URI"] = "missing"
	}

	if strings.TrimSpace(contentType) == "" {
		paramErrors["ContentType"] = "missing"
	}

	statuscode, err := ValidateStatusCode(statuscodeFV)
	if err != nil {
		paramErrors["StatusCode"] = err.Error()
	}

	if len(paramErrors) != 0 {
		validationErrors, err := json.Marshal(paramErrors)
		if err != nil {
			log.Println(err)
		}
		http.Error(w, string(validationErrors), 400)
		return
	}

	c := MockConf{URI: uri, Method: method, ContentType: contentType, StatusCode: statuscode}
	err = c.WriteToDisk([]byte(body), h.DataDirPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	http.Redirect(w, r, "/mock", http.StatusFound)
}

//DeleteHandler saves new mock entity.
func (h *Pagehandler) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Path[len("/mock/delete/"):]

	//validate input!
	reg := regexp.MustCompile("[0-9A-Za-z_]+")
	match := reg.FindAllStringSubmatch(name, -1)
	if len(match) != 1 {
		http.Error(w, "Invalid request, name may only be [0-9A-Za-z_]: "+r.URL.String(), http.StatusBadRequest)
		return
	}

	dir := h.DataDirPath
	err := os.Remove(dir + "/" + name + ConfEXT)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
	err = os.Remove(dir + "/" + name + ContentEXT)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
	return
}

func (h *Pagehandler) servePage(w http.ResponseWriter, r *http.Request) {
	log.Println("Loading all resources..")
	resources := "<ul>"
	responseList, err := LoadAll(h.DataDirPath)
	if err != nil {
		log.Println("Error while loading resources: ", err)
	}
	//newform := "<div class=\"container\"> <h4>New</h4>" + saveform + "</div>"
	newform := `<a href="mock/new/"><button type="button">New</button></a>`
	for _, i := range *responseList {
		resources = resources + fmt.Sprintf(`<div class="container" style="border: 1px solid black">
																						<div style="float:right"> <button class="delete-btn" onclick="deleteResource(this.value)" value="%s"/></div>
																						<div>%v <b>%s</b> <a href="mock/view/data/%s">View %s</a> </div>
																						<div>%s</div>
																						<a href="mock/edit/%s"><button type="button">Edit</button></a>
																					</div>`, i.GetFileName(), i.StatusCode, i.Method, i.GetFileName(), i.URI, i.ContentType, i.GetFileName())
	}

	script := `<script>
		function deleteResource(filename) {
		    if(confirm("Confirm deletion?")){
					//ajaxcall delete resource.
					console.log("we removed " + filename);

					var xmlhttp;
			    // compatible with IE7+, Firefox, Chrome, Opera, Safari
			    xmlhttp = new XMLHttpRequest();
			    xmlhttp.onreadystatechange = function(){
			        if (xmlhttp.readyState == 4 && xmlhttp.status == 200){
			            location.reload();
			        }
			    }
			    xmlhttp.open("DELETE", '/mock/delete/'+filename, true);
			    xmlhttp.send();
				}
		}
	</script>`
	resources = resources + "</ul>"

	content := fmt.Sprintf(`<h1>%s</h1><div>%s</div><div style="float:center">%s</div>`, title, resources, newform)
	page := "<html><head>" + script + css + "</head>" + "<body>" + content + "</body></html>"
	fmt.Fprintf(w, "%s", page)
}
