package main

import (
	"counterapi/modules"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
)

const (
	envPrefix       string = "COUNTERAPI"
	envRedisAddress string = "REDIS_ADDRESS"
	envRedisDB      string = "REDIS_DB"
	envListenPort   string = "PORT"
)

func main() {

	// Get parameters from environment variables
	viper.SetEnvPrefix(envPrefix)
	viper.AutomaticEnv()
	redisAddress := viper.GetString(envRedisAddress)
	redisDB := viper.GetInt(envRedisDB)
	listenPort := viper.GetString(envListenPort)
	hostname, err := os.Hostname()
	if err != nil {
		logrus.Fatal("Can't get hostname. exit")
	}

	// Inject dependencies
	redisClient, err := modules.NewRedisClient(redisAddress, redisDB)
	if err != nil {
		logrus.Fatal(err)
	}
	counter := modules.NewCounterCalculator(redisClient)
	router := modules.NewController(counter, listenPort, hostname)

	// Run
	if err := router.Run(); err != nil {
		logrus.Fatal("Failed to start.")
	}

}
