package tcp_test

import (
	"net"
	"time"

	"github.com/domdom82/go-responder/tcp"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("TcpServer", func() {

	var (
		tcpServer *tcp.TcpServer
		config    *tcp.Config
	)

	AfterEach(func() {
		tcpServer.Stop()
	})

	Describe("Basic TcpServer Tests", func() {
		Context("An empty tcp server", func() {
			config = &tcp.Config{Port: 8081}
			tcpServer = config.NewServer()
			It("should at least open a port", func() {
				tcpServer.Run()
				time.Sleep(1 * time.Second)
				_, err := net.Dial("tcp", ":8081")
				Expect(err).To(BeNil())
			})

			It("should be stoppable via stop function", func() {
				tcpServer.Run()
				time.Sleep(1 * time.Second)
				_, err := net.Dial("tcp", ":8081")
				Expect(err).To(BeNil())
				tcpServer.Stop()
				_, err = net.Dial("tcp", ":8081")
				Expect(err).To(Not(BeNil()))

			})

		})

	})
})
