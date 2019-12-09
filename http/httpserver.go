package http

import (
	"fmt"
	"net/http"
	"sort"
	"time"

	"github.com/domdom82/go-responder/common"
)

//Server represents a http server instance
type Server struct {
	config *Config
	server *http.Server
}

//Run runs the http server
func (srv *Server) Run() {
	fmt.Println("Starting http server on port", srv.config.Port)
	mux := http.NewServeMux()

	for path, response := range srv.config.Responses {
		mux.HandleFunc(path, genHandler(response))
	}

	srv.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", srv.config.Port),
		Handler: mux,
	}

	go func() {
		err := srv.server.ListenAndServe()
		if err != nil {
			fmt.Println(err)
		}
	}()

}

//Stop stops the http server
func (srv *Server) Stop() {
	_ = srv.server.Close()
}

func genHandler(response *Method) func(w http.ResponseWriter, r *http.Request) {
	seqNum := 0
	return func(w http.ResponseWriter, r *http.Request) {
		if response.Get != nil && r.Method == http.MethodGet {
			handleHTTPResponseOptions(response.Get, &seqNum, w, r)
		}
		if response.Put != nil && r.Method == http.MethodPut {
			handleHTTPResponseOptions(response.Put, &seqNum, w, r)
		}
		if response.Post != nil && r.Method == http.MethodPost {
			handleHTTPResponseOptions(response.Post, &seqNum, w, r)
		}
		if response.Delete != nil && r.Method == http.MethodDelete {
			handleHTTPResponseOptions(response.Delete, &seqNum, w, r)
		}
	}
}

func handleHTTPResponseOptions(httpResponseOptions *ResponseOptions, seqNum *int, w http.ResponseWriter, r *http.Request) {
	if httpResponseOptions.Static != nil {
		handleHTTPResponse(httpResponseOptions.Static, w, r)
	} else if httpResponseOptions.Seq != nil {
		handleHTTPResponse(httpResponseOptions.Seq[*seqNum], w, r)
		if *seqNum < len(httpResponseOptions.Seq)-1 {
			*seqNum++
		}
	} else if httpResponseOptions.Loop != nil {
		handleHTTPResponse(httpResponseOptions.Loop[*seqNum], w, r)
		*seqNum = (*seqNum + 1) % (len(httpResponseOptions.Loop))
	}

}

func handleHTTPResponse(httpResponse *Response, w http.ResponseWriter, r *http.Request) {
	if httpResponse.Delay != nil {
		time.Sleep(*httpResponse.Delay)
	}
	if len(httpResponse.Headers) > 0 {
		for key, values := range httpResponse.Headers {
			for _, value := range values {
				w.Header().Add(key, value)
			}
		}
	}
	w.WriteHeader(httpResponse.Status)

	if httpResponse.BigBody != nil {
		data := common.GenResponseData(httpResponse.BigBody.Type, httpResponse.BigBody.Size)
		w.Write(data)
	} else if httpResponse.Body != nil {
		w.Write([]byte(*httpResponse.Body))
	}
	if httpResponse.ShowHeaders {
		w.Write([]byte("\n\n----- Headers Received -----\n"))
		var keys []string
		for header := range r.Header {
			keys = append(keys, header)
		}
		sort.Strings(keys)
		for _, key := range keys {
			w.Write([]byte(fmt.Sprintf("%s : %v \n", key, r.Header[key])))
		}
	}
}
