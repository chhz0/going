package main

import (
	"fmt"

	"github.com/chhz0/going"
)

type Config struct {
	Env   string `yaml:"env"`
	App   string `yaml:"app"`
	Http  *Http  `yaml:"http"`
	Mysql *Mysql `yaml:"mysql"`
}

type Http struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}
type Mysql struct {
	Url      string `yaml:"url"`
	DB       string `yaml:"db"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

func main() {
	conf := Config{}
	cv, err := going.NewConfigV(conf)
	if err != nil {
		panic(err)
	}

	err = cv.Load(".", "config", "yaml")
	if err != nil {
		panic(err)
	}

	cv.Watch()

	for {
		fmt.Println(conf.Mysql.Url)
	}
}
