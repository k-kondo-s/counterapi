package main

import (
	"counterapi/modules"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
)

const (
	envPrefix string = "COUNTERAPI"
	envRedisAddress string = "REDIS_ADDRESS"
	envRedisDB string = "REDIS_DB"
	envListenPort string = "PORT"
)

func main() {

	viper.SetEnvPrefix(envPrefix)
	viper.AutomaticEnv()
	redisAddress := viper.GetString(envRedisAddress)
	redisDB :=viper.GetInt(envRedisDB)
	listenPort := viper.GetString(envListenPort)
	hostname, err  := os.Hostname()
	if err != nil {
		logrus.Fatal("can't get hostname. exit")
	}

	redisClient, _ := modules.NewRedisClient(redisAddress, redisDB)
	counter := modules.NewCounterCalculator(redisClient)
	router := modules.NewController(counter, listenPort, hostname)

	router.Run()

}
