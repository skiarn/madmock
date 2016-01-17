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
