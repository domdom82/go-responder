package http_test

import (
	"net"
	"time"

	"github.com/domdom82/go-responder/http"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("HttpServer", func() {

	var (
		httpServer *http.HttpServer
		config     *http.Config
	)

	AfterEach(func() {
		httpServer.Stop()
	})

	Describe("Basic HttpServer Tests", func() {
		Context("An empty http server", func() {
			config = &http.Config{Port: 8080}
			httpServer = config.NewServer()
			It("should at least open a port", func() {
				httpServer.Run()
				time.Sleep(1 * time.Second)
				_, err := net.Dial("tcp", ":8080")
				Expect(err).To(BeNil())
			})

			It("should be stoppable via stop function", func() {
				httpServer.Run()
				time.Sleep(1 * time.Second)
				_, err := net.Dial("tcp", ":8080")
				Expect(err).To(BeNil())
				httpServer.Stop()
				time.Sleep(1 * time.Second)
				_, err = net.Dial("tcp", ":8080")
				Expect(err).To(Not(BeNil()))

			})
		})
	})
})
