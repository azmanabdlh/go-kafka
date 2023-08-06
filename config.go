package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type KafkaConfig struct {
	Hostname string `yaml:"hostname"`
	Port     string `yaml:"port"`
}

type Config struct {
	Kafka KafkaConfig `yaml:"kafka"`
}

func GetConfigFromFile() (cfg Config, err error) {
	// get config file
	fileName := fmt.Sprintf("%s.yaml", "config")
	f, err := os.Open(fileName)
	if err != nil {
		return
	}

	defer f.Close()

	err = yaml.NewDecoder(f).Decode(&cfg)
	if err != nil {
		return
	}

	return
}
