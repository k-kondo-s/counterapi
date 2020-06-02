package modules

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)


type storedData struct {
	key string
	value string
	expirationSecond int64
}

type DummyDao struct {
	SetFunc func(key string, value string, expirationSecond int64) error
	GetFunc func(key string) (string, error)
	GetAllKeysFunc func() ([]string, error)
	DelFunc func(key string) error
	ExistsFunc func(key string) (int64, error)
	storedData []storedData
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
func (d *DummyDao) Exists(key string) (int64, error) {
	return d.ExistsFunc(key)
}

func TestCountCalculator_GenerateCounter(t *testing.T) {
	type testCase struct {
		id string
		startTime int64
		duration int64
		internalError error
		expectedStoredValue string
		expectedError error
	}
	var cases = []testCase{
		{
			"9dd29757-ed4e-488f-b62c-b8cececbac29",
			int64(1591115560),
			int64(1000),
			nil,
			"{\"start_timestamp\":1591115560,\"end_timestamp\":1591116560}",
			nil,
		},
		{
			"9dd29757-ed4e-488f-b62c-b8cececbac29",
			int64(1591115560),
			int64(1000),
			errors.New("some error"),
			"{\"start_timestamp\":1591115560,\"end_timestamp\":1591116560}", // actually this value should be "", but it's no matter on this test.
			nil,
		},
	}
	d := &DummyDao{}

	for _, i := range cases {
		d.SetFunc = func(key string, value string, expirationSecond int64) error {
			d.storedData = append(d.storedData, storedData{
				key:              key,
				value:            value,
				expirationSecond: expirationSecond,
			})
			return nil
		}
		c := NewCounterCalculator(d)
		c.generateUUID = func() string {return i.id}
		c.generateTimestamp = func() int64 {return i.startTime}
		id, err := c.GenerateCounter(i.duration)

		assert.Equal(t, i.expectedError, err)
		assert.Equal(t, i.id, id)
		assert.Equal(t, i.id, d.storedData[0].key)
		assert.Equal(t, i.expectedStoredValue, d.storedData[0].value)
		assert.Equal(t, i.duration, d.storedData[0].expirationSecond)
	}

}

func TestCountCalculator_GetCounter(t *testing.T) {
	type testCase struct {
		id                         string
		valueInDBCorrespondingToID string
		counterExistenceInDB int64
		daoInternalError           error
		currentTime                int64
		expectedResult             CounterResult
		expectedError              error
	}
	var cases = []testCase{
		{
			"9dd29757-ed4e-488f-b62c-b8cececbac29",
			"{\"start_timestamp\":1591115560,\"end_timestamp\":1591116560}", // end_timestamp = start_timestamp + 1000
			1,
			nil,
			int64(1591115560), // It equals to "start_timestamp",
			CounterResult{
				Current: int64(1), // so its value should be 1 because it's required that a counter has to start from 1.
				To:      1000,
				counterExistence: true,
			},
			nil,
		},
		{
			"9dd29757-ed4e-488f-b62c-b8cececbac29",
			"{\"start_timestamp\":1591115560,\"end_timestamp\":1591115570}", // end_timestamp = start_timestamp + 10
			1,
			nil,
			int64(1591115569), // It equals to start_timestamp + 9
			CounterResult{
				Current: int64(10),
				To:      10,
				counterExistence: true,
			},
			nil,
		},
		{
			"9dd29757-ed4e-488f-b62c-b8cececbac29",
			"{\"start_timestamp\":1591115560,\"end_timestamp\":1591116560}", // end_timestamp = start_timestamp + 1000
			1,
			nil,
			int64(1591116560), // It equals to end_timestamp. This seems strange but can happen if Redis works wrong unexpectedly.
			CounterResult{
				Current: 1001,
				To:      1000,
				counterExistence: false, // This means there is no counter with the given ID.
			},
			nil,
		},
		{
			"9dd29757-ed4e-488f-b62c-b8cececbac29",
			"",
			0, // No counter with the given ID in DB
			nil,
			int64(1591116560),
			CounterResult{
				Current: 0, // This means there is no counter with the given ID.
				To:      0,
				counterExistence: false,
			},
			nil,
		},
		{
			"9dd29757-ed4e-488f-b62c-b8cececbac29",
			"{\"start_timestamp\":1591115560,\"end_timestamp\":1591116560}",
			1,
			errors.New("some errors happen in DB"), // Case when internal error occurred
			int64(1591115560),
			CounterResult{
				Current: 0,
				To:      0,
				counterExistence: false,
			},
			errors.New("some errors happen in DB"),
		},
	}

	for _, i := range cases {
		d := &DummyDao{
			GetFunc: func(key string) (s string, err error) {
				s = i.valueInDBCorrespondingToID
				err = i.daoInternalError
				return
			},
			ExistsFunc: func(key string) (result int64, err error) {
				result = i.counterExistenceInDB
				err = i.daoInternalError
				return
			},
		}
		c := NewCounterCalculator(d)
		c.generateTimestamp = func() int64 {return i.currentTime}
		r, err := c.GetCounter(i.id)

		assert.Equal(t, i.expectedResult, r)
		assert.Equal(t, i.expectedError, err)
	}
}