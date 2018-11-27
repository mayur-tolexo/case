package score

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

//NewHandler : New object of HTTPHandler
func NewHandler(baseURL string) *HTTPHandler {
	return &HTTPHandler{BaseURL: baseURL, HTTPClient: &http.Client{}}
}

//Run the test case
func (hr *HTTPHandler) runTest(t *testing.T, testCase TestCase, method, path string) {
	var (
		requestBody  []byte
		responseBody []byte
		urlParams    string
		bodyParams   interface{}
		header       map[string]string
	)
	if testCase.Params != nil {
		bodyParams, urlParams, header = testCase.Params()

	}

	url := fmt.Sprintf("%s%s%s", hr.BaseURL, path, urlParams)
	t.Logf("%v %v", method, url)

	if bodyParams != nil {
		jsonValue, err := json.Marshal(bodyParams)
		require.NoError(t, err, "Failed to encode body params")
		requestBody = jsonValue
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(requestBody))
	require.NoError(t, err, "Failed to create HTTP request")

	for k, v := range header {
		req.Header.Set(k, v)
	}

	resp, err := hr.HTTPClient.Do(req)
	require.NoError(t, err, "Failed to send request")
	require.NotNil(t, resp,
		fmt.Sprintf("Request to '%s' returned nil response", url))

	if resp.Body != nil {
		defer resp.Body.Close()
		responseBody, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Error(err.Error())
			return
		}
	}

	if testCase.AssertResponse != nil {
		testCase.AssertResponse(t, testCase.ExpectedData, responseBody, resp.StatusCode)
	}
}

//Get test struct name
func getTestName(test interface{}) string {
	return reflect.TypeOf(test).String()
}

//Run all test cases one by one
func (hr *HTTPHandler) Run(t *testing.T, trial ...Trial) {
	for _, curTrial := range trial {
		for _, curAPI := range curTrial.GetAPI(t) {
			method, path, apiDesc := curAPI.GetDesc()
			for _, testcase := range curAPI.TestCases {
				t.Logf("Case: %s of %s (%s)", testcase.Desc, apiDesc, getTestName(curTrial))
				hr.runTest(t, testcase, method, path)
			}
		}
	}
}

//RunTrial will run specific trial api test cases
func (hr *HTTPHandler) RunTrial(t *testing.T, trial interface{}, apis ...API) {
	for _, api := range apis {
		method, path, apiDesc := api.GetDesc()
		for _, testcase := range api.TestCases {
			t.Logf("Case: %s of %s (%s)", testcase.Desc, apiDesc, getTestName(trial))
			hr.runTest(t, testcase, method, path)
		}
	}
}

//GetAPIDesc : Get API method path and description
func GetAPIDesc(method, path, desc string) func() (string, string, string) {
	return func() (string, string, string) { return method, path, desc }
}
