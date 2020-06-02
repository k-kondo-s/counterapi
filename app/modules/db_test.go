package modules

import (
	"fmt"
	"testing"
)

var redisAddress = "127.0.0.1:6379"
var redisDb = 0

func TestRedisClient_Set(t *testing.T) {
	r, _ := NewRedisClient(redisAddress, redisDb)
	_ = r.Set("keytest", "valuetest", 1000)

}

func TestRedisClient_Get(t *testing.T) {
	r, _ := NewRedisClient(redisAddress, redisDb)
	result, _ := r.Get("keytest")
	fmt.Println(result)
}
