package modules

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

//type Counter interface {
//	GenerateCounter(to int64) (string, error)
//	GetCounter(id string) (CounterResult, error)
//	ListAllCounterId() ([]string, error)
//	DeleteCounter(id string) error
//}

type DummyCounter struct {
	GenerateCounterFunc func(to int64) (string, error)
	GetCounterFunc func(id string) (CounterResult, error)
	ListAllCounterIdFunc func() ([]string, error)
	DeleteCounterFunc func(id string) error
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


func TestHostname(t *testing.T) {
	d := &DummyCounter{}
	c := NewController(d, "8080", "test-kenji-kondo.mac.local")
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	c.router.ServeHTTP(w, req)

	assert.Equal(t, "{\"hostname\":\"test-kenji-kondo.mac.local\"}", w.Body.String())

}