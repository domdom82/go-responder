package tcp

import (
	"net"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestTcp(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Tcp Suite")
}

var _ = Describe("TcpServer", func() {

	var (
		tcpServer *TcpServer
		config    *Config
	)

	AfterEach(func() {
		tcpServer.Stop()
	})

	Describe("Basic TcpServer Tests", func() {
		Context("An empty tcp server", func() {
			config = &Config{Port: 8081}
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
