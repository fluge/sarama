package sarama

import (
	"sync"
)

type singleton struct {
	Topic_Name string
}
var once sync.Once
var instance *singleton

func GetInstance() *singleton {
	once.Do(func() {
		instance = &singleton{"air_temperature_7"}
	})

	return instance
}

