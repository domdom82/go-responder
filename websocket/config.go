package websocket

import (
	"fmt"
	"strings"
	"time"
)

type WsResponse struct {
	Bufsize string         `yaml:"bufsize"`
	Delay   *time.Duration `yaml:"delay,omitempty"`
	Type    *string        `yaml:"type,omitempty"`
}

type Response struct {
	Read  *WsResponse `yaml:"read,omitempty"`
	Write *WsResponse `yaml:"write,omitempty"`
}

type Config struct {
	Port      int                  `yaml:"port"`
	Responses map[string]*Response `yaml:"responses"`
}

func (cfg *Config) NewServer() *WsServer {

	server := &WsServer{cfg, nil}

	return server
}

func (cfg Config) String() string {
	return fmt.Sprintf("{port: %d, responses: %v}", cfg.Port, cfg.Responses)
}

func (response Response) String() string {
	s := strings.Builder{}

	if response.Read != nil {
		s.WriteString(fmt.Sprintf(" Read: %v", *response.Read))
	}
	if response.Write != nil {
		s.WriteString(fmt.Sprintf(" Write: %v", *response.Write))
	}

	return fmt.Sprintf("{%s}", s.String())
}

func (tcpResponse WsResponse) String() string {
	s := strings.Builder{}

	s.WriteString(fmt.Sprintf("Bufsize: %s", tcpResponse.Bufsize))

	if tcpResponse.Type != nil {
		s.WriteString(fmt.Sprintf(" Type: %s", *tcpResponse.Type))
	}
	if tcpResponse.Delay != nil {
		s.WriteString(fmt.Sprintf(" Delay: %s", *tcpResponse.Delay))
	}

	return fmt.Sprintf("{%s}", s.String())
}
