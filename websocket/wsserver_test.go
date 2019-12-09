package websocket

import (
	"net"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestWebsocket(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Websocket Suite")
}

var _ = Describe("Websocket Server", func() {

	var (
		wsServer *WsServer
		config   *Config
	)

	AfterEach(func() {
		wsServer.Stop()
	})

	Describe("Basic Websocket Tests", func() {
		Context("An empty websocket server", func() {
			config = &Config{Port: 8082}
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
