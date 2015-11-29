package main

import (
	"flag"
	"os"
	"path/filepath"
)

type Settings struct {
	Port        int
	TargetURL   string
	DataDirPath string
}

func (s *Settings) Init() error {

	var port = flag.Int("p", 8080, "What port the mock should run on.")
	var url = flag.String("u", "", "Base url to system to be mocked (request will be fetched once and stored locally).")
	var dir = flag.String("d", "mad-mock-store", "Directory path to mock data and config files.")
	flag.Parse() // parse the flags

	s.Port = *port
	s.TargetURL = *url

	path, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return err
	}
	dirpath := path + "/" + *dir
	s.DataDirPath = dirpath

	return nil

}

// CreateDir creates directory relativly to execution path.
func (s *Settings) CreateDir() error {
	err := os.MkdirAll(s.DataDirPath, 0777)
	return err
}
