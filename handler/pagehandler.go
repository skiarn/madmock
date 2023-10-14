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
	"fmt"
	"log"
	"net/http"

	"github.com/skiarn/madmock/filesys"
	"github.com/skiarn/madmock/model"
)

// Pagehandler handles index page.
type Pagehandler struct {
	TargetURL   string
	DataDirPath string
	Fs          filesys.FileSystem
}

// NewPageHandler handles initzialisation of PageHandler.
func NewPageHandler(path string, targeturl string) Pagehandler {
	return Pagehandler{DataDirPath: path, TargetURL: targeturl, Fs: filesys.LocalFileSystem{}}
}

const title = "Mad Mock"
const css = "<style>" + `
									html,
									body {
										height: 100%;
										background-color: white;
									}
									body {
										color: black;
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

// statusCodeHTMLSelect returns html selector for http status code.
func statusCodeHTMLSelect() string {
	var statuscodehtml = `<select id="statuscodeS" name="StatusCode">`
	for _, c := range model.ValidStatusCodes {
		statuscodehtml += fmt.Sprintf("<option value=\"%v\">%v, %s</option>", c, c, http.StatusText(c))
	}
	return statuscodehtml + `</select>`
}
func saveformHTML() string {
	saveform := `<form method="POST" action="/mock/api/mock/">`
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

// EditHandler handles resource page to edit item.
func (h *Pagehandler) EditHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Path[len("/mock/edit/"):]
	editform := "<div class=\"container\"> <h4>Editing</h4>" + "<h5>Headers</h5><p id=\"headersP\"></p> <br>" + saveformHTML() + "</div>"
	content := fmt.Sprintf(`<h1>%s</h1><h4>%s</h4><div>%s</div>`, title, h.TargetURL, editform)
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
		    datahttp.open("GET", '/mock/api/mock/data/%s', true);
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
		    confhttp.open("GET", '/mock/api/mock/%s', true);
		    confhttp.send();


		}
	</script>`, name, name)

	page := "<html><head>" + script + css + "</head>" + `<body onload="getResource()">` + content + "</body></html>"
	fmt.Fprintf(w, "%s", page)

}

// New serves a page for creating a new mock entity.
func (h *Pagehandler) New(w http.ResponseWriter, r *http.Request) {
	postform := "<div class=\"container\"> <h4>New</h4>" + saveformHTML() + "</div>"
	content := fmt.Sprintf(`<h1>%s</h1><h4>%s</h4><div>%s</div>`, title, h.TargetURL, postform)

	page := "<html><head>" + css + "</head>" + "<body>" + content + "</body></html>"
	fmt.Fprintf(w, "%s", page)
}

func (h *Pagehandler) servePage(w http.ResponseWriter, r *http.Request) {
	resources := "<ul>"
	responseList, err := h.Fs.ReadAllMockConf(h.DataDirPath)
	if err != nil {
		log.Println("Error while loading resources: ", err)
	}
	//newform := "<div class=\"container\"> <h4>New</h4>" + saveform + "</div>"
	newform := `<a href="mock/new/"><button type="button">New</button></a>`
	for _, i := range *responseList {
		resources = resources + fmt.Sprintf(`<div class="container" style="border: 1px solid black">
																						<div style="float:right"> <button class="delete-btn" onclick="deleteResource(this.value)" value="%s"/></div>
																						<div>%v <b>%s</b> <a href="mock/api/mock/data/%s">View %s</a> </div>
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
			    xmlhttp.open("DELETE", '/mock/api/mock/'+filename, true);
			    xmlhttp.send();
				}
		}
	</script>`

	ongoingForm := fmt.Sprintf(`<div id="ongoing"> </div>`)
	wsScript := `<script>
		document.addEventListener("DOMContentLoaded", function() {
  		initWS();
		});
	function initWS(){
		var ws = null;
		if (!window["WebSocket"]) {
			alert("Error: Your browser does not support web sockets.")
		} else {
			console.log("creating ws -> ws://` + r.Host + `/mock/wsmockinfo")
			ws = new WebSocket("ws://` + r.Host + `/mock/wsmockinfo");
			ws.onmessage = function(e) {

				var localRaw = localStorage.getItem('mockresponses');
				if (typeof(localRow) !== 'undefined') {
					var list = Object.keys(localRaw).map(function(k) { return localRaw[k] });
					if (! (Array.isArray(list))) {
						list = [];
					}
				} else {
					list = [];
				}

				var mockItem = JSON.parse(e.data)
				console.log("event" + mockItem)
				list.push(mockItem)

				var mstr = JSON.stringify(list[0])
				console.log("request:" + mstr);
				addToOngoing(list[0])

				//localStorage.setItem('mockresponses', list);

			};
			ws.onerror = function(e) {
					console.log("got error:", e);
			};
			ws.onclose = function(e) {
					console.log("got close:", e);
			};
		}
	}
	function addToOngoing(item) {
	var a = document.createElement('a');
	var linkText = document.createTextNode(item.uri);
	a.appendChild(linkText);
	a.title = item.uri;
	a.href = item.uri;

	var totalSeconds = new Date().getTime() / 1000;
	var hours = parseInt( totalSeconds / 3600 ) % 24;
	var minutes = parseInt( totalSeconds / 60 ) % 60;
	var seconds = parseInt ( totalSeconds % 60 );
	var timestamp = (hours < 10 ? "0" + hours : hours) + ":" + (minutes < 10 ? "0" + minutes : minutes) + ":" + (seconds  < 10 ? "0" + seconds : seconds);

	var timestampNode = document.createTextNode(timestamp);
	var timestampSpan = document.createElement("span");
	timestampSpan.style.cssFloat = "left";
	timestampSpan.appendChild(timestampNode);

	var statusNode = document.createElement("span");
	statusNode.appendChild(document.createTextNode(item.status));
	statusNode.style.paddingRight = "5em";
	statusNode.style.paddingLeft = "5em";
	statusNode.style.cssFloat = "left";
	var boldcontent = document.createElement("b");
	boldcontent.style.cssFloat = "left";
	boldcontent.appendChild(document.createTextNode(item.method));

  var newDiv = document.createElement("div");
	newDiv.appendChild(timestampSpan);
	newDiv.appendChild(statusNode);
  newDiv.appendChild(boldcontent); //add the text node to the newly created div.
	newDiv.appendChild(a);
  // add the newly created element and its content into the DOM
  var currentDiv = document.getElementById("ongoing");
	//currentDiv.appendChild(newDiv);
	currentDiv.insertBefore(newDiv, currentDiv.firstChild);

  //document.body.insertBefore(newDiv, currentDiv);
}</script>`

	resources = resources + "</ul>"

	content := fmt.Sprintf(`<h1>%s</h1><h4>%s</h4><div>%s</div><div>%s</div><div style="float:center">%s</div>`, title, h.TargetURL, ongoingForm, resources, newform)
	page := "<html><head>" + script + wsScript + css + "</head>" + "<body>" + content + "</body></html>"
	fmt.Fprintf(w, "%s", page)
}
