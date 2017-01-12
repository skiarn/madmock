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

package setting_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/skiarn/madmock/setting"
)

func TestInitSetting_WhenUrlHasPortNumber(t *testing.T) {

	oldArgs := os.Args

	os.Args = []string{"cmd", "-u=google.com:9090"}
	expectedURL := "google.com:9090"
	expectedParnetDirName := "google.com9090"

	settings, err := setting.Create()
	if err != nil {
		t.Errorf("Error occured while trying to create settings: %v", err)
	}

	if settings.TargetURL != expectedURL {
		t.Errorf("Expected: %v but got: %v", expectedURL, settings.TargetURL)
	}

	segments := strings.Split(settings.DataDirPath, string(filepath.Separator))
	parentDir := segments[len(segments)-1]
	if parentDir != expectedParnetDirName {
		t.Errorf("Expected: %v but got: %v", expectedParnetDirName, parentDir)
	}
	defer func() { os.Args = oldArgs }()
}
