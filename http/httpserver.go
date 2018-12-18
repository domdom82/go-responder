package http

import (
	"bytes"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"code.cloudfoundry.org/bytefmt"
)

type HttpServer struct {
	config *Config
}

const ResponseTypeBinary string = "binary"
const ResponseTypeLorem string = "lorem"

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
		data := genResponseData(httpResponse.BigBody.Type, httpResponse.BigBody.Size)
		w.Write(data)
	} else if httpResponse.Body != nil {
		w.Write([]byte(*httpResponse.Body))
	}
}

//TODO: refactor into own file to be used by tcp and websocket servers
func genResponseData(responseType string, size string) []byte {
	payload := new(bytes.Buffer)

	sizeBytes, err := bytefmt.ToBytes(size)

	if err != nil {
		fmt.Println(err)
		return payload.Bytes()
	}

	switch responseType {
	case ResponseTypeBinary:
		for i := uint64(0); i < sizeBytes; i++ {
			payload.WriteByte(byte(rand.Intn(255)))
		}
	case ResponseTypeLorem:
		//TODO: use lorem generator
	}

	return payload.Bytes()
}
