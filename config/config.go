package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

var Config DevConfig

type DevConfig struct {
	HTTPS bool `default:"false" env:"HTTPS" json:"https"`
	Port  uint `default:"7005" env:"PORT" json:"port"`
	Database    struct {
		Default   struct {
			Host  	string `json:"host"`
			Port     string `json:"port"`
			User     string `json:"user"`
			Pwd 	 string `json:"pwd"`
			Name 	 string `json:"name"`
		} `json:"default"`
	} `json:"database"`
}

func init() {
	content, err := ioutil.ReadFile("config.json")
	println(string(content))
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(content, &Config)
	if err != nil {
		panic(err)
	}
	fmt.Println(Config)
}
