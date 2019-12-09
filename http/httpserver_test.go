package http_test

import (
	"io/ioutil"
	"net"
	"net/http"
	"testing"
	"time"

	httpResponder "github.com/domdom82/go-responder/http"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestHttp(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "HttpServer Suite")
}

var _ = Describe("HttpServer", func() {

	Context("An empty http server", func() {
		config := &httpResponder.Config{Port: 8080}
		httpServer := config.NewServer()

		BeforeEach(func() {
			err := httpServer.Run()
			Expect(err).To(BeNil())
			time.Sleep(1 * time.Second)
		})

		AfterEach(func() {
			_ = httpServer.Stop()
		})

		It("should at least open a port", func() {
			_, err := net.Dial("tcp", ":8080")
			Expect(err).To(BeNil())
		})

		It("should be stoppable via stop function", func() {
			_, err := net.Dial("tcp", ":8080")
			Expect(err).To(BeNil())
			err = httpServer.Stop()
			Expect(err).To(BeNil())
			time.Sleep(1 * time.Second)
			_, err = net.Dial("tcp", ":8080")
			Expect(err).To(Not(BeNil()))

		})
	})

	Context("A simple endpoint", func() {
		config := &httpResponder.Config{Port: 8081, Responses: map[string]*httpResponder.Response{
			"/test": {Get: &httpResponder.HttpResponseOptions{Static: &httpResponder.HttpResponse{Status: 200}}},
		}}
		httpServer := config.NewServer()

		BeforeEach(func() {
			err := httpServer.Run()
			Expect(err).To(BeNil())
			time.Sleep(1 * time.Second)
		})

		AfterEach(func() {
			err := httpServer.Stop()
			Expect(err).To(BeNil())
			time.Sleep(1 * time.Second)
		})

		It("should reply to a request", func() {
			response, err := http.Get("http://localhost:8081/test")
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
		})

		It("should have no body", func() {
			response, err := http.Get("http://localhost:8081/test")
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			bytes, err := ioutil.ReadAll(response.Body)
			Expect(string(bytes)).To(Equal(""))
		})

	})

	Context("An endpoint with a body", func() {
		body := "OK"
		config := &httpResponder.Config{Port: 8081, Responses: map[string]*httpResponder.Response{
			"/test2": {Get: &httpResponder.HttpResponseOptions{Static: &httpResponder.HttpResponse{Status: 200, Body: &body}}},
		}}
		httpServer := config.NewServer()

		BeforeEach(func() {
			err := httpServer.Run()
			Expect(err).To(BeNil())
			time.Sleep(1 * time.Second)
		})

		AfterEach(func() {
			_ = httpServer.Stop()
		})

		It("should have a body", func() {
			response, err := http.Get("http://localhost:8081/test2")
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			bytes, err := ioutil.ReadAll(response.Body)
			Expect(string(bytes)).To(Equal(body))
			Expect(err).To(BeNil())
		})
	})

})
