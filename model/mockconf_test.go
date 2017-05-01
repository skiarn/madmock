package model_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/skiarn/madmock/model"
)

func TestValid(t *testing.T) {
        m := model.MockConf {URI: "/", Method: "POST", ContentType: "application/json", StatusCode: 200}
        got, err := m.Valid()
        if got != true {
		t.Errorf("Expected valid but got errors: %v", err)
        } 
}

func TestValid_StatusCode(t *testing.T) {
        m := model.MockConf {StatusCode: 200}
        got, err := m.Valid()
        if got == false {
                if err["StatusCode"] != "" {
                        t.Errorf("Expected not missing but got: %v", err["StatusCode"])
                }
        } else {
                t.Errorf("Expected validation errors. But no errors detected.", m)
        }
}

func TestValid_StatusCode_Invalid(t *testing.T) {
        m := model.MockConf {StatusCode: 1}
        got, err := m.Valid()
        if got == false {
                if err["StatusCode"] != "1 is unknown" {
                        t.Errorf("Expected 1 is unknown but got: %v", err["StatusCode"])
                }
        } else {
                t.Errorf("Expected validation errors. But no errors detected.", m)
        }
}

//TestValid_ContentType
func TestValid_ContentType(t *testing.T) {
        m := model.MockConf {ContentType: "application/json"}
        got, err := m.Valid()
        if got == false {
                if err["ContentType"] != "" {
                        t.Errorf("Expected not missing but got: %v", err["ContentType"])
                }
        } else {
                t.Errorf("Expected validation errors. But no errors detected.", m)
        }
}

//TestValid_ContentType_Missing
func TestValid_ContentType_Missing(t *testing.T) {
        m := model.MockConf {ContentType: ""}
        got, err := m.Valid()
        if got == false {
                if err["ContentType"] != "missing" {
                        t.Errorf("Expected missing but got: %v", err["ContentType"])
                }
        } else {
                t.Errorf("Expected validation errors. But no errors detected.", m)
        }
}

//TestValid_Method
func TestValid_Method(t *testing.T) {
        m := model.MockConf {Method: "POST"}
        got, err := m.Valid()
        if got == false {
                if err["Method"] != "" {
                        t.Errorf("Expected not missing but got: %v", err["Method"])
                }
        } else {
                t.Errorf("Expected validation errors. But no errors detected.", m)
        }
}

//TestValid_Method_Missing
func TestValid_Method_Missing(t *testing.T) {
        m := model.MockConf {URI: "", Method: "", ContentType: "", StatusCode: 200}
        got, err := m.Valid()
        if got == false {
                if err["Method"] != "missing" {
                        t.Errorf("Expected missing but got: %v", err["Method"])
                }
        } else {
                t.Errorf("Expected validation errors. But no errors detected.", m)
        }
}

//TestValid_URI
func TestValid_URI(t *testing.T) {
	m := model.MockConf {URI: "/"}
	got, err := m.Valid()
	if got == false {
		if err["URI"] != "" {
			t.Errorf("Expected not missing but got: %v", err["URI"])
		}
	} else {
                t.Errorf("Expected validation errors. But no errors detected.", m)
	}
}

//TestValid_URI_Missing
func TestValid_URI_Missing(t *testing.T) {
	m := model.MockConf {URI: ""}
	got, err := m.Valid()
	if got == false {
		if err["URI"] != "missing" {
			t.Errorf("Expected missing but got: %v", err["URI"])
		}
	} else {
                t.Errorf("Expected validation errors. But no errors detected.", m)
	}
}

func TestGetMockFileName_ShouldReturnFilename(t *testing.T) {
	//Equals to base 32 uri = / method = GET
	expectedFilename := "AZH5Y2SPPMLFNUGNVXRJGON64IUV3FNB"
	req, err := http.NewRequest("GET", "http://myservice.com/", strings.NewReader(string("body")))
	if err != nil {
		t.Errorf("%v", err)
	}
	got := model.GetMockFileName(req)

	if got != expectedFilename {
		t.Errorf("Expected: %v but got: %v", expectedFilename, got)
	}
}

func TestGetFileName_ShouldReturnEncodedFilename(t *testing.T) {
	//Equals to base 32 uri = / method = GET
	expectedFilename := "AZH5Y2SPPMLFNUGNVXRJGON64IUV3FNB"
	m := model.MockConf{URI: "/", Method: "GET", ContentType: "contentType", StatusCode: 418}
	got := m.GetFileName()
	if got != expectedFilename {
		t.Errorf("Expected: %v but got: %v", expectedFilename, got)
	}
}

func TestGetFileNameWithQueries_ShouldReturnEncodedFilename(t *testing.T) {
	//GET-/example/query?value=abc&value2=cba = RYK2253FOOJI773UXZASWXUM4AA3CTEF
	expectedFilename := "RYK2253FOOJI773UXZASWXUM4AA3CTEF"
	m := model.MockConf{URI: "/example/query?value=abc&value2=cba", Method: "GET", ContentType: "contentType", StatusCode: 418}
	got := m.GetFileName()
	if got != expectedFilename {
		t.Errorf("Expected: %v but got: %v", expectedFilename, got)
	}
}

func TestReadResponse_ShouldBuildMockFromResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "RandomContentType")
		w.WriteHeader(418)
		fmt.Fprintln(w, `body with some random data`)
	}))
	defer server.Close()

	// Make a transport that reroutes all traffic to the example server
	transport := &http.Transport{
		Proxy: func(req *http.Request) (*url.URL, error) {
			return url.Parse(server.URL)
		},
	}

	// Make a http.Client with the transport
	httpClient := &http.Client{Transport: transport}

	resp, err := httpClient.Get("http://google.se/test/get")
	if err != nil {
		t.Errorf("Error occured unexpectedly with error: %s", err)
	}

	//Test
	m := model.ReadResponse(resp)
	expectedURI := "/test/get"
	expectedMethod := "GET"
	expectedStatus := 418
	expectedContentType := "RandomContentType"
	if m.URI != expectedURI {
		t.Errorf("Expected: %v but got: %v", expectedURI, m.URI)
	}
	if m.Method != expectedMethod {
		t.Errorf("Expected: %v but got: %v", expectedMethod, m.Method)
	}
	if m.StatusCode != expectedStatus {
		t.Errorf("Expected: %v but got: %v", expectedStatus, m.StatusCode)
	}
	if m.ContentType != expectedContentType {
		t.Errorf("Expected: %v but got: %v", expectedContentType, m.ContentType)
	}
	//verify headers.
	for k, v := range resp.Header {
		for _, vv := range v {
			if m.Header[k] == "" {
				t.Errorf("Expected header: %v with value %v but got %v.", k, vv, m.Header[k])
			}
		}
	}
}

func Test_String(t *testing.T) {
	expectedString := "URI:/test, Method:GET, Contenttype:contentType, Status:418\n"
	m := model.MockConf{URI: "/test", Method: "GET", ContentType: "contentType", StatusCode: 418}
	got := m.String()
	if got != expectedString {
		t.Errorf("Expected: %v but got: %v", expectedString, got)
	}
}
