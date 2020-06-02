package modules

import (
	"encoding/json"
	"github.com/google/uuid"
	"time"
)

type CounterResult struct {
	Current          int64 `json:"current"`
	To               int64 `json:"to"`
	counterExistence bool
}

type Counter interface {
	GenerateCounter(to int64) (string, error)
	GetCounter(id string) (CounterResult, error)
	ListAllCounterId() ([]string, error)
	DeleteCounter(id string) error
}

type CountCalculator struct {
	dao Dao
	generateUUID func() string
	generateTimestamp func() int64
}

type DaoValueFormat struct {
	StartTimestamp int64 `json:"start_timestamp"`
	EndTimestamp int64 `json:"end_timestamp"`
}

// Initialize CounterCalculator.
func NewCounterCalculator(dao Dao) *CountCalculator {
	c := new(CountCalculator)
	c.dao = dao
	c.generateUUID = func() string { return uuid.New().String() }
	c.generateTimestamp = func() int64 { return time.Now().Unix() }
	return c
}

// Generate a new counter
func (c *CountCalculator) GenerateCounter(to int64) (string, error) {
	id := c.generateUUID()
	//id := uuid.New().String()
	startTimestamp := c.generateTimestamp()
	value, _ := daoValueFormatter(startTimestamp, to)
	err := c.dao.Set(id, value, to)
	if err != nil {
		return "", err
	}
	return id, nil
}

// Get a current counter of a given ID.
// Note: This is the core implementation of "Counter API".
// I got the current counter by calculating the endTimestamp - startTimesamp + 1.
// The architecture can let whole system immutable.
// Note: I didn't implement the behavior when a counter comes to the expire date,
// alternatively, Redis cares about it.
func (c *CountCalculator) GetCounter(id string) (CounterResult, error) {

	counterResult := CounterResult{}

	// Check the counter with the given ID exists in DB
	existence, errExists := c.dao.Exists(id)

	// If internal error occurs in DB, return error
	if errExists != nil {
		return counterResult, errExists
	}

	// Put the existence to the result
	counterResult.counterExistence = convertIntToBool(existence)

	// Get the counter from DB if it exists.
	if counterResult.counterExistence {
		r, errGet := c.dao.Get(id)
		if errGet != nil {
			return counterResult, errGet
		}
		var rFormatted DaoValueFormat
		_ = json.Unmarshal([]byte(r), &rFormatted)

		// Calculate counter
		counterResult.Current = c.generateTimestamp() - rFormatted.StartTimestamp + 1
		counterResult.To = rFormatted.EndTimestamp - rFormatted.StartTimestamp

		// If a case which is something wrong as the following happens, return "the counter doesn't exist".
		if counterResult.Current > counterResult.To {
			counterResult.counterExistence = false
			return counterResult, nil
		}
	}

	return counterResult, nil
}

// List all registered counter IDs
func (c *CountCalculator) ListAllCounterId() ([]string, error) {
	results, err := c.dao.GetAllKeys()
	if err != nil {
		return []string{}, err
	}
	return results, nil
}

// Delete the counter with the given ID
func (c *CountCalculator) DeleteCounter(id string) error {
	err := c.dao.Del(id)
	return err
}

// Formatter for the value in DB
func daoValueFormatter(startTimestamp int64, to int64) (string, error) {
	result := DaoValueFormat{
		StartTimestamp: startTimestamp,
		EndTimestamp:   startTimestamp + to,
	}
	resultJson, err := json.Marshal(result)
	return string(resultJson), err
}

// Convert 0 or else -> false or true
func convertIntToBool(i int64) bool {
	switch i {
	case int64(0):
		return false
	default:
		return true
	}
}