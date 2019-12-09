package tcp

import (
	"fmt"
	"strings"
	"time"
)

//Response represents a tcp response
type Response struct {
	Bufsize string         `yaml:"bufsize"`
	Delay   *time.Duration `yaml:"delay,omitempty"`
	Type    *string        `yaml:"type,omitempty"`
}

//Stream represents a tcp read/write stream
type Stream struct {
	Read  *Response `yaml:"read,omitempty"`
	Write *Response `yaml:"write,omitempty"`
}

//Config represents a tcp server configuration
type Config struct {
	Port      int     `yaml:"port"`
	Responses *Stream `yaml:"responses"`
}

//NewServer creates a new tcp server
func (cfg *Config) NewServer() *Server {

	server := &Server{cfg, nil}

	return server
}

func (cfg Config) String() string {
	return fmt.Sprintf("{port: %d, responses: %v}", cfg.Port, cfg.Responses)
}

func (response Stream) String() string {
	s := strings.Builder{}

	if response.Read != nil {
		s.WriteString(fmt.Sprintf(" Read: %v", *response.Read))
	}
	if response.Write != nil {
		s.WriteString(fmt.Sprintf(" Write: %v", *response.Write))
	}

	return fmt.Sprintf("{%s}", s.String())
}

func (tcpResponse Response) String() string {
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
