package main

import (
	"io/ioutil"

	"github.com/domdom82/go-responder/http"
	"github.com/domdom82/go-responder/tcp"
	"github.com/domdom82/go-responder/websocket"
	"gopkg.in/yaml.v2"
)

type ServerConfig struct {
	Http      *http.Config
	Tcp       *tcp.Config
	Websocket *websocket.Config
}

type Config struct {
	ServerConfigs []*ServerConfig `yaml:"servers"`
}

func NewServerConfigFromFile(filename string) (*Config, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var config *Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
