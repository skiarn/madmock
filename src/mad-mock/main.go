package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	var port = flag.String("p", "8080", "What port the mock should run on.")
	var url = flag.String("u", "", "Base url to system to be mocked (request will be fetched once and stored locally).")
	var dir = flag.String("d", "mad-mock-store", "Directory path to mock data and config files.")
	flag.Parse() // parse the flags

	err := CreateDir(*dir)
	if err != nil {
		log.Fatal(err)
	}
	mux := http.NewServeMux()
	mockPagehandler := MockPagehandler{DataDir: *dir, TargetURL: *url}
	mux.Handle("/mock", &mockPagehandler)
	mux.HandleFunc("/", mockhandler)
	log.Println("Server to listen on a port: ", *port)
	log.Fatal(http.ListenAndServe(":"+*port, mux))
}

// CreateDir creates directory relativly to execution path.
func CreateDir(dir string) error {
	path, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return err
	}
	dirpath := path + "/" + dir
	fmt.Println(dirpath)
	err = os.MkdirAll(dirpath, 0777)
	return err
}

// MockPagehandler handles index page.
type MockPagehandler struct {
	DataDir   string
	TargetURL string
}

func (h *MockPagehandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
