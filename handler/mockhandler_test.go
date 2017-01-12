package handler_test

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/skiarn/madmock/handler"
	"github.com/skiarn/madmock/model"
)

const testDataDirPath = "/path/mad-mock-store"

//Fake File system to read dummy mockconf from.
type DummyFileSystem struct{}

func (DummyFileSystem) ReadMockConf(filepath string) (*model.MockConf, error) {
	// RYK2253FOOJI773UXZASWXUM4AA3CTEF = GET /example/query?value=abc&value2=cba
	if filepath == testDataDirPath+"/RYK2253FOOJI773UXZASWXUM4AA3CTEF.mc" {
		dummyMockConf1 := model.MockConf{URI: "/example/query?value=abc&value2=cba", Method: "GET", ContentType: "test-content-v.1", StatusCode: 200}
		return &dummyMockConf1, nil
	}
	return nil, errors.New("Unknown test mockconf.")
}
func (DummyFileSystem) WriteMock(c model.MockConf, content []byte, dirpath string) error {
	log.Printf("Writing mock %v to dummy fs.", c.String())
	return nil
}
func (DummyFileSystem) ReadAllMockConf(dirpath string) (*model.MockConfs, error) {
	var confs model.MockConfs
	dummyMockConf1 := model.MockConf{URI: "/path/query?value=1", Method: "GET", ContentType: "test-content-v.1", StatusCode: 200}
	confs = append(confs, dummyMockConf1)
	return &confs, nil
}
func (DummyFileSystem) ReadResource(filepath string) (io.Reader, error) {
	if filepath == testDataDirPath+"/AZH5Y2SPPMLFNUGNVXRJGON64IUV3FNB.data" {
		mockdata := []byte(`body with some random data`)
		r := bytes.NewReader(mockdata)
		return r, nil
	}
	return nil, errors.New("Unknown test resource.")
}

type DummyClient struct{}

func (DummyClient) RequestTargetInfo(URL string, w http.ResponseWriter, r *http.Request) (*http.Response, error) {
	//Creating faker server.
	server := httptest.NewServer(http.HandlerFunc(func(sw http.ResponseWriter, sr *http.Request) {
		w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
		w.WriteHeader(200)
		fmt.Fprintln(sw, `body with some random data`)
	}))
	defer server.Close()

	// reroutes all traffic to the fake server
	transport := &http.Transport{
		Proxy: func(req *http.Request) (*url.URL, error) {
			return url.Parse(server.URL)
		},
	}
	httpClient := &http.Client{Transport: transport}

	resp, err := httpClient.Get(URL)
	return resp, err
}

func TestServeHTTP_ValidGETRequest(t *testing.T) {
	handler := handler.Mockhandler{TargetURL: "github.com", DataDirPath: testDataDirPath, Fs: DummyFileSystem{}, Client: DummyClient{}}
	test := GenerateHandleTester(t, &handler)
	//w := test("GET", url.Values{})
	w := test("GET", nil)
	if w.Code != http.StatusOK {
		t.Errorf("Mock didn't return %v, response value: %v", http.StatusOK, w)
	}

	expectedBody := `body with some random data`
	gotBody := string(w.Body.Bytes())
	if expectedBody != gotBody {
		t.Errorf("Mocked page didn't return %v, but got: %v", expectedBody, gotBody)
	}
}

func TestBuildTargetRequestURL_NormalURL(t *testing.T) {
	expectedURL := "http://google.se/path"
	mockhandler := handler.NewMockhandler("google.se", "/dir/test")
	req, err := http.NewRequest("GET", "http://myservice.com/path", strings.NewReader(string("body")))
	if err != nil {
		t.Errorf("%v", err)
	}

	url := mockhandler.BuildTargetRequestURL(req)

	if url != expectedURL {
		t.Errorf("Expected: %v but got: %v", expectedURL, url)
	}
}

func TestBuildTargetRequestURL_URLWithQuery(t *testing.T) {
	expectedURL := "http://google.se/example/query?value=abc&value2=cba"
	mockhandler := handler.NewMockhandler("google.se", "/dir/test")
	req, err := http.NewRequest("GET", "http://myservice.com/example/query?value=abc&value2=cba", strings.NewReader(string("body")))
	if err != nil {
		t.Errorf("%v", err)
	}

	url := mockhandler.BuildTargetRequestURL(req)

	if url != expectedURL {
		t.Errorf("Expected: %v but got: %v", expectedURL, url)
	}
}

func TestBuildTargetRequestURL_URLWithPortSpecified(t *testing.T) {
	expectedURL := "http://localhost:8080/path"
	mockhandler := handler.NewMockhandler("localhost:8080", "/dir/test")
	req, err := http.NewRequest("GET", "http://myservice.com/path", strings.NewReader(string("body")))
	if err != nil {
		t.Errorf("%v", err)
	}

	url := mockhandler.BuildTargetRequestURL(req)

	if url != expectedURL {
		t.Errorf("Expected: %v but got: %v", expectedURL, url)
	}
}

type HandleTester func(
	method string,
	//params url.Values,
	body []byte,
) *httptest.ResponseRecorder

// Given the current test runner and an http.Handler, generate a
// HandleTester which will test its given input against the
// handler.

func GenerateHandleTester(
	t *testing.T,
	handleFunc http.Handler,
) HandleTester {

	// Given a method type ("GET", "POST", etc) and
	// parameters, serve the response against the handler and
	// return the ResponseRecorder.

	return func(
		method string,
		//params url.Values,
		body []byte,
	) *httptest.ResponseRecorder {

		req, err := http.NewRequest(
			method,
			"",
			//strings.NewReader(params.Encode()),
			strings.NewReader(string(body)),
		)
		if err != nil {
			t.Errorf("%v", err)
		}
		req.Header.Set(
			"Content-Type",
			"application/x-www-form-urlencoded; param=value",
		)
		w := httptest.NewRecorder()
		handleFunc.ServeHTTP(w, req)
		return w
	}
}
