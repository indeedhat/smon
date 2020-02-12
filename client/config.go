package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type Config struct {
	Clients []ClientConfig `yaml:",flow"`
}

type ClientConfig struct {
	Name     string
	Interval int64
	Server   struct {
		Host   string
		Port   uint16
		Socket string
		Bin    string
		User   string
		SSH    struct {
			Keyfile  string
			Password string
		}
	}
	Modules struct {
		Cpu     bool
		Memory  bool
		Date    bool
		Uptime  bool
		Disk    []string `yaml:",flow"`
		Network struct {
			Enable    bool
			Interface []string `yaml:",flow"`
		}
	}
}

var config Config

func init() {
	config = Config{}

	data, err := ioutil.ReadFile("config.yml")
	if nil != err {
		log.Fatalf("Failed to load config file: %s\n", err)
	}

	err = yaml.Unmarshal(data, &config)
	if nil != err {
		log.Fatalf("Failed to parse config file: %s\n", err)
	}
}
