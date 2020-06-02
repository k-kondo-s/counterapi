package modules

import (
	"testing"
)

//type Dao interface {
//	Set(key string, value string, expirationSecond int64) error
//	Get(key string) (string, error)
//	GetAllKeys(prefix string) ([]string, error)
//	Del(key string) error
//}


type DummyDao struct {
	SetFunc func(key string, value string, expirationSecond int64) error
	GetFunc func(key string) (string, error)
	GetAllKeysFunc func() ([]string, error)
	DelFunc func(key string) error
}

func (d *DummyDao) Set(key string, value string, expirationSecond int64) error {
	return d.SetFunc(key, value, expirationSecond)
}
func (d *DummyDao) Get(key string) (string, error) {
	return d.GetFunc(key)
}
func (d *DummyDao) GetAllKeys() ([]string, error) {
	return d.GetAllKeysFunc()
}
func (d *DummyDao) Del(key string) error {
	return d.DelFunc(key)
}

func TestCounterCalculator_GenerateCounter(t *testing.T) {
	d := &DummyDao{
		SetFunc: func(key string, value string, expirationSecond int64) error {
			return nil
		},
	}
	c := NewCounterCalculator(d)
	var i int64 = 1000
	_, err := c.GenerateCounter(i)
	if err != nil {
		t.Fatal()
	}
}
