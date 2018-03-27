package config

import (
	"fmt"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Parser struct {
		TimeSleep    int `yaml:"timeSleep"`
		Gol          int
		RootDir      string `yaml:"rootDir"`
		HrefAllLinks string `yaml:"hrefAllLinks"`
		MobAgent     string `yaml:"mobAgent"`
		DescAgent    string `yaml:"descAgent"`
		Links        map[string]string
		CountLink    int `yaml:"countLink"`
		Index        string
		Script       string
		NotIframe    bool `yaml:"notIframe"`
		Ajax         bool
		Dir          string
	}
}

func GetConfig() *Config {
	var conf Config
	source, err := ioutil.ReadFile("./config.yml")
	if err != nil {
		panic(err)
	}
	errY := yaml.Unmarshal(source, &conf)
	if errY != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Printf("--- config:\n%v\n\n", conf)
	return &conf
}
