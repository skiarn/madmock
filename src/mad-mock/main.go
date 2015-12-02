package main

import (
	"log"
	"net/http"
	"strconv"
)

func main() {
	settings := Settings{}
	settings.Init()
	err := settings.CreateDir()
	if err != nil {
		log.Fatal(err)
	}
	mux := http.NewServeMux()
	pagehandler := Pagehandler{settings: settings}
	mockhandler := Mockhandler{settings: settings}
	mux.Handle("/mock", &pagehandler)
	mux.HandleFunc("/mock/add", pagehandler.handleMockconfPost)
	mux.Handle("/", &mockhandler)
	log.Println("Server to listen on a port: ", settings.Port)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(settings.Port), mux))
}
