package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

func main() {
	var port = flag.String("p", "8080", "What port the mock should run on.")
	var url = flag.String("u", "", "Base url to system to be mocked (request will be fetched once and stored locally).")
	var dir = flag.String("d", "./mad-mock", "Directory path to mock data and config files.")
	flag.Parse() // parse the flags

	fmt.Println("Port:" + *port + " Url:" + *url + " Dir:" + *dir)

	mux := http.NewServeMux()
	mux.HandleFunc("/mock", mockPagehandler)
	mux.HandleFunc("/", mockhandler)
	log.Println("Server to listen on a port: ", *port)
	log.Fatal(http.ListenAndServe(":"+*port, mux))
}

func mockPagehandler(w http.ResponseWriter, r *http.Request) {
	title := "Available rescources:"
	resources := "<ul>"
	resources = resources + fmt.Sprintf("<div> <p> MOCK <p></div> <br>")
	resources = resources + "</ul>"
	content := fmt.Sprintf("<h1>%s</h1><div>%s</div>", title, resources)
	page := "<html>" + content + "</html>"
	fmt.Fprintf(w, "%s", page)
}

func mockhandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "TODO")
}
