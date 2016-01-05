package main

import (
	"log"
	"madmock/handler"
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
	pagehandler := handler.Pagehandler{DataDirPath: settings.DataDirPath}
	mockhandler := handler.Mockhandler{TargetURL: settings.TargetURL, DataDirPath: settings.DataDirPath}

	mux.Handle("/mock", &pagehandler)
	mux.HandleFunc("/mock/view/data/", pagehandler.ViewDataHandler)
	mux.HandleFunc("/mock/view/conf/", pagehandler.ViewConfHandler)
	mux.HandleFunc("/mock/edit/", pagehandler.EditHandler)
	mux.HandleFunc("/mock/new/", pagehandler.New)
	mux.HandleFunc("/mock/save/", pagehandler.Save)
	mux.HandleFunc("/mock/delete/", pagehandler.DeleteHandler)
	mux.Handle("/", &mockhandler)
	log.Println("Server to listen on a port: ", settings.Port)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(settings.Port), mux))
}
