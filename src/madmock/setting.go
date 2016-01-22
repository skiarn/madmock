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
	"flag"
	"os"
	"path/filepath"
	"strings"
)

type Settings struct {
	Port        int
	TargetURL   string
	DataDirPath string
}

func (s *Settings) Init() error {

	var port = flag.Int("p", 9988, "What port the mock should run on.")
	var url = flag.String("u", "", "Base url to system to be mocked (request will be fetched once and stored locally).")
	var dir = flag.String("d", "mad-mock-store", "Directory path to mock data and config files.")
	flag.Parse() // parse the flags

	s.Port = *port
	s.TargetURL = *url

	path, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return err
	}

	// Removes : in "example:8080" since : is not valid directory character.
	dirpath := path + "/" + *dir + "/" + strings.Replace(*url, ":", "_", -1)
	s.DataDirPath = dirpath

	return nil

}

// CreateDir creates directory relativly to execution path.
func (s *Settings) CreateDir() error {
	err := os.MkdirAll(s.DataDirPath, 0777)
	return err
}
