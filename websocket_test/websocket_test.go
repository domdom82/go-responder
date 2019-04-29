package websocket_test

import (
	"net"
	"time"

	"github.com/domdom82/go-responder/websocket"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Websocket Server", func() {

	var (
		wsServer *websocket.WsServer
		config   *websocket.Config
	)

	AfterEach(func() {
		wsServer.Stop()
	})

	Describe("Basic Websocket Tests", func() {
		Context("An empty websocket server", func() {
			config = &websocket.Config{Port: 8082}
			wsServer = config.NewServer()
			It("should at least open a port", func() {
				wsServer.Run()
				time.Sleep(1 * time.Second)
				_, err := net.Dial("tcp", ":8082")
				Expect(err).To(BeNil())
			})

			It("should be stoppable via stop function", func() {
				wsServer.Run()
				time.Sleep(1 * time.Second)
				_, err := net.Dial("tcp", ":8082")
				Expect(err).To(BeNil())
				wsServer.Stop()
				_, err = net.Dial("tcp", ":8082")
				Expect(err).To(Not(BeNil()))

			})
		})
	})
})
