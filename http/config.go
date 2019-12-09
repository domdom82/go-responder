package http

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

//BigBody is a special kind of response bode with large size. Can be lorem ipsum text or binary random.
type BigBody struct {
	Size string `yaml:"size"`
	Type string `yaml:"type"`
}

//Response represents a single http response
type Response struct {
	Headers     http.Header    `yaml:"headers"`
	Status      int            `yaml:"status"`
	ShowHeaders bool           `yaml:"showheaders,omitempty"`
	Body        *string        `yaml:"body,omitempty"`
	BigBody     *BigBody       `yaml:"bigbody,omitempty"`
	Delay       *time.Duration `yaml:"delay,omitempty"`
}

//ResponseOption enumerates types of possible http responses: static, sequence or loop
type ResponseOptions struct {
	Static *Response   `yaml:"static,omitempty"`
	Seq    []*Response `yaml:"seq,omitempty"`
	Loop   []*Response `yaml:"loop,omitempty"`
}

//Method maps a HTTP method to a response option
type Method struct {
	Get    *ResponseOptions `yaml:"get,omitempty"`
	Put    *ResponseOptions `yaml:"put,omitempty"`
	Post   *ResponseOptions `yaml:"post,omitempty"`
	Delete *ResponseOptions `yaml:"delete,omitempty"`
}

//Config represents a http server configuration
type Config struct {
	Port      int                `yaml:"port"`
	Responses map[string]*Method `yaml:"responses"`
}

//NewServer creates a new HTTP server from a http.config struct
func (cfg *Config) NewServer() *HttpServer {

	server := &HttpServer{cfg, nil}

	return server
}

func (cfg Config) String() string {
	return fmt.Sprintf("{port: %d, responses: %v}", cfg.Port, cfg.Responses)
}

func (resp Method) String() string {
	s := strings.Builder{}

	if resp.Get != nil {
		s.WriteString(fmt.Sprintf(" GET: %v", resp.Get))
	}
	if resp.Put != nil {
		s.WriteString(fmt.Sprintf(" PUT: %v", resp.Put))
	}
	if resp.Post != nil {
		s.WriteString(fmt.Sprintf(" POST: %v", resp.Post))
	}
	if resp.Delete != nil {
		s.WriteString(fmt.Sprintf(" DELETE: %v", resp.Delete))
	}

	return fmt.Sprintf("{%s}", s.String())
}

func (respOptions ResponseOptions) String() string {
	s := strings.Builder{}

	if respOptions.Static != nil {
		s.WriteString(fmt.Sprintf(" Static: %v", respOptions.Static))
	}
	if respOptions.Seq != nil {
		s.WriteString(fmt.Sprintf(" Seq: %v", respOptions.Seq))
	}

	return fmt.Sprintf("{%s}", s.String())
}

func (httpResponse Response) String() string {
	s := strings.Builder{}

	s.WriteString(fmt.Sprintf("Status: %d", httpResponse.Status))

	if httpResponse.Body != nil {
		s.WriteString(fmt.Sprintf(" Body: %s", *httpResponse.Body))
	}
	if httpResponse.BigBody != nil {
		s.WriteString(fmt.Sprintf(" BigBody: %v", httpResponse.BigBody))
	}
	if httpResponse.Delay != nil {
		s.WriteString(fmt.Sprintf(" Delay: %s", *httpResponse.Delay))
	}
	s.WriteString(fmt.Sprintf(" ShowHeaders: %v", httpResponse.ShowHeaders))

	return fmt.Sprintf("{%s}", s.String())
}

func (bigBody BigBody) String() string {
	return fmt.Sprintf("{type: %s, size: %s}", bigBody.Type, bigBody.Size)
}
