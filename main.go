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
	"net/http"
	"os"
	"strconv"

	"github.com/skiarn/madmock/handler"
	"github.com/skiarn/madmock/setting"
	"github.com/skiarn/madmock/ws"

	"golang.org/x/net/websocket"
)

var logger *log.Logger
var errorlog *os.File

func main() {
	//errorlog, err := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	//if err != nil {
	//	fmt.Printf("error opening file: %v", err)
	//	os.Exit(1)
	//}
	//defer errorlog.Close()
	logger = log.New(os.Stdout, "applog: ", log.Lshortfile|log.LstdFlags)

	settings, err := setting.Create()
	if err != nil {
		log.Fatal(err)
	}
	err = settings.CreateDir()
	if err != nil {
		log.Fatal(err)
	}
	mux := http.NewServeMux()
	pagehandler := handler.NewPageHandler(settings.DataDirPath, settings.TargetURL)
	curdhandler := handler.NewMockCURDHandler(settings.DataDirPath)
	viewDataHandler := handler.NewViewDataHandler(settings.DataDirPath)

	wsHandler := ws.NewHandler(logger)
	mockhandler := handler.NewMockhandler(settings.TargetURL, settings.DataDirPath, *wsHandler)
	mux.Handle("/mock/wsmockinfo", websocket.Handler(wsHandler.WSMockInfoServer))
	go wsHandler.Run()

	mux.Handle("/mock", &pagehandler)
	mux.Handle("/mock/api/mock/", &curdhandler)
	mux.Handle("/mock/api/mock/data/", &viewDataHandler)

	mux.HandleFunc("/mock/edit/", pagehandler.EditHandler)
	mux.HandleFunc("/mock/new/", pagehandler.New)
	mux.Handle("/", &mockhandler)
	logger.Printf("Server to listen on a port: %v \n", settings.Port)
	logger.Fatal(http.ListenAndServe(":"+strconv.Itoa(settings.Port), mux))
}
