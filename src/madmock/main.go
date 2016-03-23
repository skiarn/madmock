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

package main

import (
	"log"
	"madmock/handler"
	"madmock/setting"
	"net/http"
	"strconv"
)

func main() {
	settings, err := setting.Create()
	if err != nil {
		log.Fatal(err)
	}
	err = settings.CreateDir()
	if err != nil {
		log.Fatal(err)
	}
	mux := http.NewServeMux()
	pagehandler := handler.NewPageHandler(settings.DataDirPath)
	curdhandler := handler.NewMockCURDHandler(settings.DataDirPath)
	viewDataHandler := handler.NewViewDataHandler(settings.DataDirPath)
	mockhandler := handler.NewMockhandler(settings.TargetURL, settings.DataDirPath)
	mux.Handle("/mock", &pagehandler)
	mux.Handle("/mock/api/mock/", &curdhandler)
	mux.Handle("/mock/api/mock/data/", &viewDataHandler)

	mux.HandleFunc("/mock/edit/", pagehandler.EditHandler)
	mux.HandleFunc("/mock/new/", pagehandler.New)
	mux.Handle("/", &mockhandler)
	log.Println("Server to listen on a port: ", settings.Port)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(settings.Port), mux))
}
