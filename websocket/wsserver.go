package websocket

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"code.cloudfoundry.org/bytefmt"

	"github.com/domdom82/go-responder/common"

	"github.com/gorilla/websocket"
)

type WsServer struct {
	config *Config
	server *http.Server
}

func (srv *WsServer) Run() {
	fmt.Println("Starting web socket server on port", srv.config.Port)
	srv.server = &http.Server{Addr: fmt.Sprintf(":%d", srv.config.Port)}

	for path, response := range srv.config.Responses {
		http.HandleFunc(path, genHandler(response))
	}

	go func() {
		err := srv.server.ListenAndServe()
		if err != nil {
			fmt.Println(err)
		}
	}()

}

func (srv *WsServer) Stop() {
	_ = srv.server.Close()
}

func genHandler(response *Response) func(w http.ResponseWriter, r *http.Request) {

	readBufSize, _ := bytefmt.ToBytes(response.Read.Bufsize)
	writeBufSize, _ := bytefmt.ToBytes(response.Write.Bufsize)

	upgrader := &websocket.Upgrader{
		ReadBufferSize:  int(readBufSize),
		WriteBufferSize: int(writeBufSize),
	}

	upgrader.CheckOrigin = func(r *http.Request) bool {
		//TODO: maybe add some fancy html clients under /<ws-path>/client later to allow origin check again
		return true
	}

	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)

		if err != nil {
			panic(err)
		}

		go handleConn(response, conn)
	}

}

func handleConn(response *Response, conn *websocket.Conn) {
	fmt.Printf("\n(%v) New connection from %v\n", conn.LocalAddr(), conn.RemoteAddr())
	defer func() {
		_ = conn.Close()
	}()
	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func() {
		for {
			if response.Read.Delay != nil {
				time.Sleep(*response.Read.Delay)
			}

			messageType, buf, err := conn.ReadMessage()
			fmt.Printf("\nread %d bytes of type %d from %v", len(buf), messageType, conn.RemoteAddr())

			if err != nil {
				fmt.Printf("%s : %s\n", conn.RemoteAddr(), err)
				break
			}
		}
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		for {
			if response.Write.Delay != nil {
				time.Sleep(*response.Write.Delay)
			}

			dataType := common.ResponseTypeBinary
			if response.Write.Type != nil {
				dataType = *response.Write.Type
			}
			data := common.GenResponseData(dataType, response.Write.Bufsize)
			msgType := msgTypeFromResponseType(dataType)

			err := conn.WriteMessage(msgType, data)
			fmt.Printf("\nwrote %d bytes of type %d to %v", len(data), msgType, conn.RemoteAddr())

			if err != nil {
				fmt.Printf("%s : %s\n", conn.RemoteAddr(), err)
				break
			}
		}
		wg.Done()
	}()

	wg.Wait()
}

func msgTypeFromResponseType(responseType string) int {
	msgType := 0
	switch responseType {
	case common.ResponseTypeBinary:
		msgType = websocket.BinaryMessage
	case common.ResponseTypeLorem:
		msgType = websocket.TextMessage
	}
	return msgType
}
