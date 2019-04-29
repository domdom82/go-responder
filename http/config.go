package http

import (
	"fmt"
	"strings"
	"time"
)

type BigBody struct {
	Size string `yaml:"size"`
	Type string `yaml:"type"`
}

type HttpResponse struct {
	Status      int            `yaml:"status"`
	ShowHeaders bool           `yaml:"showheaders,omitempty"`
	Body        *string        `yaml:"body,omitempty"`
	BigBody     *BigBody       `yaml:"bigbody,omitempty"`
	Delay       *time.Duration `yaml:"delay,omitempty"`
}

type HttpResponseOptions struct {
	Static *HttpResponse   `yaml:"static,omitempty"`
	Seq    []*HttpResponse `yaml:"seq,omitempty"`
	Loop   []*HttpResponse `yaml:"loop,omitempty"`
}

type Response struct {
	Get    *HttpResponseOptions `yaml:"get,omitempty"`
	Put    *HttpResponseOptions `yaml:"put,omitempty"`
	Post   *HttpResponseOptions `yaml:"post,omitempty"`
	Delete *HttpResponseOptions `yaml:"delete,omitempty"`
}

type Config struct {
	Port      int                  `yaml:"port"`
	Responses map[string]*Response `yaml:"responses"`
}

func (cfg *Config) NewServer() *HttpServer {

	server := &HttpServer{cfg, nil}

	return server
}

func (cfg Config) String() string {
	return fmt.Sprintf("{port: %d, responses: %v}", cfg.Port, cfg.Responses)
}

func (resp Response) String() string {
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

func (respOptions HttpResponseOptions) String() string {
	s := strings.Builder{}

	if respOptions.Static != nil {
		s.WriteString(fmt.Sprintf(" Static: %v", respOptions.Static))
	}
	if respOptions.Seq != nil {
		s.WriteString(fmt.Sprintf(" Seq: %v", respOptions.Seq))
	}

	return fmt.Sprintf("{%s}", s.String())
}

func (httpResponse HttpResponse) String() string {
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
