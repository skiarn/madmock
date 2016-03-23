package handler_test

import (
	"madmock/handler"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestBuildTargetRequestURL_NormalURL(t *testing.T) {
	expectedURL := "http://google.se/path"
	mockhandler := handler.NewMockhandler("google.se", "/dir/test")
	req, err := http.NewRequest("GET", "http://myservice.com/path", strings.NewReader(string("body")))
	if err != nil {
		t.Errorf("%v", err)
	}

	url, err := mockhandler.BuildTargetRequestURL(req)
	if err != nil {
		t.Errorf("Error occured: %v", err)
	}

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

	url, err := mockhandler.BuildTargetRequestURL(req)
	if err != nil {
		t.Errorf("Error occured: %v", err)
	}

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
