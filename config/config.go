package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

var Config DevConfig
type DevConfig struct {
	HTTPS bool `default:"false" env:"HTTPS"`
	Port  uint `default:"7005" env:"PORT"`
	Database    struct {
		Name     string `env:"DBName" default:"qor_example"`
		Adapter  string `env:"DBAdapter" default:"mysql"`
		Host     string `env:"DBHost" default:"localhost"`
		Port     string `env:"DBPort" default:"3306"`
		User     string `env:"DBUser"`
		Password string `env:"DBPassword"`
	}
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
