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

package handler_test

import (
	"errors"
	"io"
	"testing"

	"github.com/skiarn/madmock/handler"
	"github.com/skiarn/madmock/model"
)

// fsMockImpl is implementation of application mocked filesystem.
type fsMockImpl struct{}

func (fsMockImpl) Remove(name string) error {
	return errors.New("File can not be removed because this is a mock test!")
}

func (fsMockImpl) ReadMockConf(filepath string) (*model.MockConf, error) {
	return nil, errors.New("ReadMockConf went horrible testy wrong")
}
func (fsMockImpl) WriteMock(c model.MockConf, content []byte, dirpath string) error {
	return errors.New("WriteMock went horrible testy wrong")
}
func (fsMockImpl) ReadAllMockConf(dataDirPath string) (*model.MockConfs, error) {
	return &model.MockConfs{}, nil
}
func (fsMockImpl) ReadResource(filepath string) (io.Reader, error) {
	return nil, errors.New("ReadResource went horrible testy wrong")
}

func TestPagehandlerHandleGet_Page(t *testing.T) {
	expectedReturnCode := 200
	pagehandler := handler.Pagehandler{DataDirPath: "/test", Fs: fsMockImpl{}}

	test := GenerateHandleTester(t, &pagehandler)

	w := test("GET", nil)
	if w.Code != expectedReturnCode {
		t.Errorf("Expected %v but got %v", expectedReturnCode, w.Code)
	}
	//expect body not empty
	if w.Body.String() == "" {
		t.Errorf("Expected body not empty but got %v", w.Body.String())
	}
}
