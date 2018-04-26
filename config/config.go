package config

import (
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
		CountLink    int    `yaml:"countLink"`
		IndexDesc    string `yaml:"indexDesc"`
		IndexMob     string `yaml:"indexMob"`
		Script       string
		NotIframe    bool `yaml:"notIframe"`
		Ajax         bool
		Dir          string
		Mob          bool
		AddHeader    []string `yaml:"addHeader"`
		RemoveAllTag  []string `yaml:"removeAllTag"`
	}
}

func GetConfig() *Config {
	var conf Config
	source, err := ioutil.ReadFile("./config.yml")
	if err != nil {
		log.Fatalf("File config not found!\n")
	}
	errY := yaml.Unmarshal(source, &conf)
	if errY != nil {
		log.Fatalf("File config no valid - error: %v", err)
	}
	return &conf
}
