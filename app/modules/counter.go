package modules

import (
	"encoding/json"
	"github.com/google/uuid"
	"time"
)

type CounterResult struct {
	Current int64 `json:"current"`
	To int64 `json:"to"`
}

type Counter interface {
	GenerateCounter(to int64) (string, error)
	GetCounter(id string) (CounterResult, error)
	ListAllCounterId() ([]string, error)
	DeleteCounter(id string) error
}

type CountCalculator struct {
	dao Dao
}

type DaoValueFormat struct {
	StartTimestamp int64 `json:"start_timestamp"`
	EndTimestamp int64 `json:"end_timestamp"`
}


func NewCounterCalculator(dao Dao) *CountCalculator {
	c := new(CountCalculator)
	c.dao = dao
	return c
}

func (c *CountCalculator) GenerateCounter(to int64) (string, error) {
	id := uuid.New().String()
	startTimestamp := time.Now().Unix()
	value, errDaoValueFormatter := daoValueFormatter(startTimestamp, to)
	if errDaoValueFormatter != nil {
		return "", errDaoValueFormatter
	}
	errSet := c.dao.Set(id, value, to)
	if errSet != nil {
		return "", errSet
	}
	return id, nil
}

func (c *CountCalculator) GetCounter(id string) (CounterResult, error) {
	r, errGet := c.dao.Get(id)
	if errGet != nil {
		return CounterResult{}, errGet
	}
	var rFormatted DaoValueFormat
	errUnmarshal := json.Unmarshal([]byte(r), &rFormatted)
	if errUnmarshal != nil {
		return CounterResult{}, errUnmarshal
	}
	result := CounterResult{
		Current: time.Now().Unix() - rFormatted.StartTimestamp + 1,
		To:      rFormatted.EndTimestamp - rFormatted.StartTimestamp,
	}
	return result, nil
}

func (c *CountCalculator) ListAllCounterId() ([]string, error) {
	results, err := c.dao.GetAllKeys()
	if err != nil {
		return []string{}, err
	}
	return results, nil
}

func (c *CountCalculator) DeleteCounter(id string) error {
	err := c.dao.Del(id)
	return err
}

func daoValueFormatter(startTimestamp int64, to int64) (string, error) {
	result := DaoValueFormat{
		StartTimestamp: startTimestamp,
		EndTimestamp:   startTimestamp + to,
	}
	resultJson, err := json.Marshal(result)
	return string(resultJson), err
}