package tcp

import (
	"fmt"
	"net"
	"sync"
	"time"

	"code.cloudfoundry.org/bytefmt"

	"github.com/domdom82/go-responder/common"
)

//Server represents a tcp server instance
type Server struct {
	config   *Config
	listener net.Listener
}

func (srv *Server) handleConn(conn net.Conn) {
	fmt.Printf("\n(%v) New connection from %v\n", conn.LocalAddr(), conn.RemoteAddr())
	defer func() {
		_ = conn.Close()
	}()
	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func() {
		if srv.config.Responses != nil {
			readBufSize, _ := bytefmt.ToBytes(srv.config.Responses.Read.Bufsize)
			readBuf := make([]byte, int(readBufSize))
			for {
				if srv.config.Responses.Read.Delay != nil {
					time.Sleep(*srv.config.Responses.Read.Delay)
				}
				nbytes, err := conn.Read(readBuf)
				fmt.Printf("\nread %d bytes from %v", nbytes, conn.RemoteAddr())
				if err != nil {
					fmt.Printf(" (%s)\n", err)
					break
				}
			}
		}
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		for {
			if srv.config.Responses != nil {
				if srv.config.Responses.Write.Delay != nil {
					time.Sleep(*srv.config.Responses.Write.Delay)
				}
				dataType := common.ResponseTypeBinary
				if srv.config.Responses.Write.Type != nil {
					dataType = *srv.config.Responses.Write.Type
				}
				data := common.GenResponseData(dataType, srv.config.Responses.Write.Bufsize)

				nbytes, err := conn.Write(data)
				fmt.Printf("\nwrote %d bytes to %v", nbytes, conn.RemoteAddr())
				if err != nil {
					fmt.Printf(" (%s)\n", err)
					break
				}
			}
		}
		wg.Done()
	}()

	wg.Wait()
}

//Run starts a tcp server
func (srv *Server) Run() {
	fmt.Println("Starting tcp server on port", srv.config.Port)
	listener, _ := net.Listen("tcp", fmt.Sprintf(":%d", srv.config.Port))
	srv.listener = listener

	go func() {
		defer func() {
			_ = listener.Close()
		}()
		for {
			conn, err := listener.Accept()
			if err != nil {
				return
			}

			go srv.handleConn(conn)
		}
	}()
}

//Stop stops a tcp server
func (srv *Server) Stop() {
	fmt.Println("Closing tcp server on", srv.listener.Addr().String())
	_ = srv.listener.Close()
}
