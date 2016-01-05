package handler_test

import (
	"encoding/json"
	"madmock/handler"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestPagehandlerHandleGet_WhenMissingBodyID(t *testing.T) {
	expectedBody := `{"id":"missing"}`
	expectedReturnCode := 400
	pagehandler := handler.Pagehandler{DataDirPath: "/test"}
	test := GenerateHandleTester(t, &pagehandler)
	//w := test("GET", url.Values{})
	w := test("GET", []byte(`{"id":"123456ID"}`))
	if w.Code != expectedReturnCode {
		t.Errorf("Home page return %v Expected: %v", w.Code, expectedReturnCode)
	}

	response := make(map[string]string)
	//err := json.Unmarshal(w.Body.Bytes(), &response)
	decoder := json.NewDecoder(w.Body)
	err := decoder.Decode(&response)

	if err != nil {
		t.Errorf("Error occured while Unmarshaling response...", err)
	}
	respData, _ := json.Marshal(response)
	got := string(respData)
	if got != expectedBody {
		t.Errorf("Expected: %v but got: %v", expectedBody, got)
	}
}

func xTestPagehandlerHandleGet(t *testing.T) {
	pagehandler := handler.Pagehandler{DataDirPath: "/test"}
	test := GenerateHandleTester(t, &pagehandler)
	//w := test("GET", url.Values{})
	w := test("GET", []byte(`{"id":"123456ID"}`))
	if w.Code != http.StatusOK {
		t.Errorf("Home page didn't return %v, response value: %v", http.StatusOK, w)
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
