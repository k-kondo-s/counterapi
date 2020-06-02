package modules

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// DummyCounter implementing Counter interface
type DummyCounter struct {
	GenerateCounterFunc  func(to int64) (string, error)
	GetCounterFunc       func(id string) (CounterResult, error)
	ListAllCounterIdFunc func() ([]string, error)
	DeleteCounterFunc    func(id string) error
}

func (d *DummyCounter) GenerateCounter(to int64) (string, error) {
	return d.GenerateCounterFunc(to)
}
func (d *DummyCounter) GetCounter(id string) (CounterResult, error) {
	return d.GetCounterFunc(id)
}
func (d *DummyCounter) ListAllCounterId() ([]string, error) {
	return d.ListAllCounterIdFunc()
}
func (d *DummyCounter) DeleteCounter(id string) error {
	return d.DeleteCounterFunc(id)
}

// return hostname with JSON formatted against the request "/"
func TestRouterGetHostname(t *testing.T) {
	d := &DummyCounter{}
	c := NewController(d, "8080", "test-kenji-kondo.mac.local")
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	c.router.ServeHTTP(w, req)

	assert.Equal(t, "{\"hostname\":\"test-kenji-kondo.mac.local\"}", w.Body.String())

}

// tests of GET /counter
func TestRouterGetAllCounterIDs(t *testing.T) {
	type testCase struct {
		registeredIds      []string
		internalError      error
		expectedBody       string
		expectedHttpStatus int
	}
	var cases = []testCase{
		{
			[]string{"1a0ca312-558f-4a13-987f-ba86930ec9ef", "3f2ead43-5a97-4b14-8bb9-3fbf1dfe1f4e"},
			nil,
			"{\"ids\":[\"1a0ca312-558f-4a13-987f-ba86930ec9ef\",\"3f2ead43-5a97-4b14-8bb9-3fbf1dfe1f4e\"]}",
			200,
		},
		{
			[]string{},
			nil,
			"{\"ids\":[]}",
			200,
		},
		{
			nil,
			errors.New("some error"),
			"{\"error\":\"Internal Server Error\"}",
			500,
		},
	}

	for _, i := range cases {
		d := &DummyCounter{
			ListAllCounterIdFunc: func() (strings []string, err error) {
				strings = i.registeredIds
				err = i.internalError
				return
			},
		}
		c := NewController(d, "", "")
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/counter", nil)
		c.router.ServeHTTP(w, req)
		assert.Equal(t, i.expectedBody, w.Body.String())
		assert.Equal(t, i.expectedHttpStatus, w.Code)
	}
}

// tests of POST /counter?to=[int]
func TestRouterGenerateCounter(t *testing.T) {
	type testCase struct {
		queryString        string
		internalError      error
		generatedId        string
		expectedBody       string
		expectedHttpStatus int
	}
	var cases = []testCase{
		{
			"?to=1000",
			nil,
			"3f2ead43-5a97-4b14-8bb9-3fbf1dfe1f4e",
			"{\"id\":\"3f2ead43-5a97-4b14-8bb9-3fbf1dfe1f4e\"}",
			201,
		},
		{
			"?to=0",
			nil,
			"3f2ead43-5a97-4b14-8bb9-3fbf1dfe1f4e",
			"{\"id\":\"3f2ead43-5a97-4b14-8bb9-3fbf1dfe1f4e\"}",
			201,
		},
		{
			"?to=kondokenji",
			nil,
			"",
			"{\"error\":\"the value kondokenji is invalid\"}",
			400,
		},
		{
			"?to=",
			nil,
			"",
			"{\"error\":\"param to is required\"}",
			400,
		},
		{
			"",
			nil,
			"",
			"{\"error\":\"param to is required\"}",
			400,
		},
		{
			"?to=1000",
			errors.New("some error"),
			"",
			"{\"error\":\"Internal Server Error\"}",
			500,
		},
	}

	for _, i := range cases {
		d := &DummyCounter{GenerateCounterFunc: func(to int64) (s string, err error) {
			s = i.generatedId
			err = i.internalError
			return
		}}
		c := NewController(d, "", "")
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/counter"+i.queryString, nil)
		c.router.ServeHTTP(w, req)
		assert.Equal(t, i.expectedBody, w.Body.String())
		assert.Equal(t, i.expectedHttpStatus, w.Code)
	}
}

// tests of GET /counter/:id
func TestRouterGetCurrentCounter(t *testing.T) {
	type testCase struct {
		inputID        string
		currentCounter CounterResult
		internalError error
		expectedBody   string
		expectedStatus int
	}
	var cases = []testCase{
		{
			"/3f2ead43-5a97-4b14-8bb9-3fbf1dfe1f4e",
			CounterResult{
				Current: 10,
				To:      1000,
				counterExistence: true,
			},
			nil,
			"{\"current\":10,\"to\":1000}",
			200,
		},
		{
			"/3f2ead43-5a97-4b14-8bb9-3fbf1dfe1f4e",
			CounterResult{
				Current: 0,
				To:      0,
				counterExistence: false,
			},
			nil,
			"{\"error\":\"no such counter with 3f2ead43-5a97-4b14-8bb9-3fbf1dfe1f4e\"}",
			404,
		},
		{
			"/3f2ead43-5a97-4b14-8bb9-3fbf1dfe1f4e",
			CounterResult{
				Current: 0,
				To:      0,
				counterExistence: false,
			},
			errors.New("some error"),
			"{\"error\":\"Internal Server Error\"}",
			500,
		},
	}

	for _, i := range cases {
		d := &DummyCounter{GetCounterFunc: func(id string) (result CounterResult, err error) {
			result = i.currentCounter
			err = i.internalError
			return
		}}
		c := NewController(d, "", "")
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/counter"+i.inputID, nil)
		c.router.ServeHTTP(w, req)
		assert.Equal(t, i.expectedBody, w.Body.String())
		assert.Equal(t, i.expectedStatus, w.Code)
	}
}

// tests of POST /counter/:id/stop
func TestRouterDeleteCounter(t *testing.T) {
	type testCase struct {
		inputID        string
		internalError  error
		expectedStatus int
	}
	var cases = []testCase{
		{
			"/3f2ead43-5a97-4b14-8bb9-3fbf1dfe1f4e/stop",
			nil,
			204,
		},
		{
			"/3f2ead43-5a97-4b14-8bb9-3fbf1dfe1f4e/stop",
			errors.New("some error"),
			500,
		},
	}

	for _, i := range cases {
		d := &DummyCounter{DeleteCounterFunc: func(id string) error {
			return i.internalError
		}}
		c := NewController(d, "", "")
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/counter"+i.inputID, nil)
		c.router.ServeHTTP(w, req)
		assert.Equal(t, i.expectedStatus, w.Code)
	}
}

// tests of default routing
func TestRouterNotFound(t *testing.T) {
	type testCase struct {
		path           string
		method         string
		expectedBody   string
		expectedStatus int
	}
	var cases = []testCase{
		{
			"/kenji",
			http.MethodPost,
			"{\"error\":\"Not Found\"}",
			404,
		},
		{
			"/counter/stop",
			http.MethodPost,
			"{\"error\":\"Not Found\"}",
			404,
		},
	}

	for _, i := range cases {
		d := &DummyCounter{}
		c := NewController(d, "", "")
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(i.method, i.path, nil)
		c.router.ServeHTTP(w, req)
		assert.Equal(t, i.expectedBody, w.Body.String())
		assert.Equal(t, i.expectedStatus, w.Code)
	}
}
