package main

import (
	"fmt"
	"io/ioutil"

	"github.com/domdom82/go-responder/http"
	"github.com/domdom82/go-responder/tcp"
	"github.com/domdom82/go-responder/websocket"
	"gopkg.in/yaml.v2"
)

type ServerConfig struct {
	Http      *http.Config      `yaml:"http"`
	Tcp       *tcp.Config       `yaml:"tcp"`
	Websocket *websocket.Config `yaml:"websocket"`
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

func (cfg ServerConfig) String() string {
	s := ""
	if cfg.Http != nil {
		s = fmt.Sprintf("{http: %v}", cfg.Http)
	}
	if cfg.Tcp != nil {
		s = fmt.Sprintf("{tcp: %v}", cfg.Tcp)
	}
	if cfg.Websocket != nil {
		s = fmt.Sprintf("{websocket: %v}", cfg.Websocket)
	}

	return s
}
