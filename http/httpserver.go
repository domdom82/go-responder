package http

import (
	"bytes"
	"fmt"
	"github.com/domdom82/datarate"
	"io"
	"math"
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
	var data []byte
	if httpResponse.BigBody != nil {
		data = common.GenResponseData(httpResponse.BigBody.Type, httpResponse.BigBody.Size)
	} else if httpResponse.Body != nil {
		data = []byte(*httpResponse.Body)
	}
	if httpResponse.Rate != nil {
		writeAtRate(w, data, httpResponse.Rate)
	} else {
		w.Write(data)
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

func writeAtRate(w http.ResponseWriter, data []byte, rate *datarate.Datarate) {
	bytesPerSecond := rate.BytesPerSecond()
	bytesPerSecondInt := int(math.Ceil(bytesPerSecond))
	buf := make([]byte, bytesPerSecondInt)
	byteReader := bytes.NewReader(data)

	for {
		tStart := time.Now()
		n, err := byteReader.Read(buf)
		if n == 0 && err == io.EOF {
			break
		}
		w.Write(buf[:n])
		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		}
		tEnd := time.Now()
		tDuration := tEnd.Sub(tStart)
		tWait := time.Second - tDuration

		if tWait > 0 && n == bytesPerSecondInt {
			time.Sleep(tWait)
		}
	}
}
