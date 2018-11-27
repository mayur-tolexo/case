package score

import (
	"net/http"
	"testing"
)

//Trial : Interface which contain all test API struct
type Trial interface {
	GetAPI(*testing.T) []API
}

//API : API Model of test cases
type API struct {
	GetDesc   APIDesc
	TestCases []TestCase
}

//TestCase : Deatils of API test case
type TestCase struct {
	Desc           string
	Params         ParamsFunc
	ExpectedData   interface{}
	AssertResponse AssertResponseFunc
}

//APIDesc : Api description function containing method type, path of the api and desctiprion
type APIDesc func() (method, path, desc string)

//ParamsFunc : Parameters of the api
type ParamsFunc func() (body interface{}, url string, header map[string]string)

//AssertResponseFunc : Assert response of the api
type AssertResponseFunc func(t *testing.T, expected interface{},
	responseBody []byte, code int)

//HTTPHandler : Object Containing the base url and http client object
type HTTPHandler struct {
	BaseURL    string
	HTTPClient *http.Client
}
