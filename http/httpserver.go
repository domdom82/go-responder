package http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/domdom82/go-responder/common"
)

type HttpServer struct {
	config *Config
}

func (srv *HttpServer) Run() {
	fmt.Println("Starting http server on port", srv.config.Port)

	for path, response := range srv.config.Responses {
		http.HandleFunc(path, genHandler(response))
	}

	http.ListenAndServe(fmt.Sprintf(":%s", srv.config.Port), nil)
}

func genHandler(response *Response) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if response.Get != nil && r.Method == http.MethodGet {
			handleHttpResponseOptions(response.Get, w, r)
		}
		//TODO: handle other methods
	}

}

func handleHttpResponseOptions(httpResponseOptions *HttpResponseOptions, w http.ResponseWriter, r *http.Request) {
	if httpResponseOptions.Static != nil {
		handleHttpResponse(httpResponseOptions.Static, w, r)
	} else if httpResponseOptions.Seq != nil {
		//TODO: handle sequence. maintain state
	}

}

func handleHttpResponse(httpResponse *HttpResponse, w http.ResponseWriter, r *http.Request) {
	if httpResponse.Delay != nil {
		time.Sleep(*httpResponse.Delay)
	}
	w.WriteHeader(httpResponse.Status)

	if httpResponse.BigBody != nil {
		data := common.GenResponseData(httpResponse.BigBody.Type, httpResponse.BigBody.Size)
		w.Write(data)
	} else if httpResponse.Body != nil {
		w.Write([]byte(*httpResponse.Body))
	}
}
