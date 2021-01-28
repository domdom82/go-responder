package http_test

import (
	"github.com/domdom82/datarate"
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
	RunSpecs(t, "Server Suite")
}

var _ = Describe("Server", func() {

	Context("An empty http server", func() {
		config := &httpResponder.Config{Port: 8080}
		httpServer := config.NewServer()

		BeforeEach(func() {
			httpServer.Run()
			time.Sleep(1 * time.Second)
		})

		AfterEach(func() {
			httpServer.Stop()
		})

		It("should at least open a port", func() {
			_, err := net.Dial("tcp", ":8080")
			Expect(err).To(BeNil())
		})

		It("should be stoppable via stop function", func() {
			_, err := net.Dial("tcp", ":8080")
			Expect(err).To(BeNil())
			httpServer.Stop()
			time.Sleep(1 * time.Second)
			_, err = net.Dial("tcp", ":8080")
			Expect(err).To(Not(BeNil()))

		})
	})

	Context("A simple endpoint", func() {
		config := &httpResponder.Config{Port: 8081, Responses: map[string]*httpResponder.Method{
			"/test": {Get: &httpResponder.ResponseOptions{Static: &httpResponder.Response{Status: 200}}},
		}}
		httpServer := config.NewServer()

		BeforeEach(func() {
			httpServer.Run()
			time.Sleep(1 * time.Second)
		})

		AfterEach(func() {
			httpServer.Stop()
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
			Expect(err).To(BeNil())
		})

	})

	Context("An endpoint with a body", func() {
		body := "OK"
		config := &httpResponder.Config{Port: 8081, Responses: map[string]*httpResponder.Method{
			"/test2": {Get: &httpResponder.ResponseOptions{Static: &httpResponder.Response{Status: 200, Body: &body}}},
		}}
		httpServer := config.NewServer()

		BeforeEach(func() {
			httpServer.Run()
			time.Sleep(1 * time.Second)
		})

		AfterEach(func() {
			httpServer.Stop()
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

	Context("An endpoint with headers configured", func() {
		body := "OK"
		config := &httpResponder.Config{Port: 8081, Responses: map[string]*httpResponder.Method{
			"/headers": {Get: &httpResponder.ResponseOptions{Static: &httpResponder.Response{
				Headers: http.Header{"X-TestHeader": []string{"value"}},
				Status:  200,
				Body:    &body,
			}}},
		}}
		httpServer := config.NewServer()

		BeforeEach(func() {
			httpServer.Run()
			time.Sleep(1 * time.Second)
		})

		AfterEach(func() {
			httpServer.Stop()
		})

		It("should return headers", func() {
			response, err := http.Get("http://localhost:8081/headers")
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			bytes, err := ioutil.ReadAll(response.Body)
			Expect(string(bytes)).To(Equal(body))
			Expect(err).To(BeNil())
			Expect(response.Header.Get("X-TestHeader")).To(Equal("value"))
		})
	})

	Context("An endpoint with a slow data rate", func() {
		body := "123"
		rate, _ := datarate.Parse("1B/s")
		config := &httpResponder.Config{Port: 8081, Responses: map[string]*httpResponder.Method{
			"/slow": {Get: &httpResponder.ResponseOptions{Static: &httpResponder.Response{
				Rate:   rate,
				Status: 200,
				Body:   &body,
			}}},
		}}
		httpServer := config.NewServer()

		BeforeEach(func() {
			httpServer.Run()
			time.Sleep(1 * time.Second)
		})

		AfterEach(func() {
			httpServer.Stop()
		})

		It("should take 3 seconds to transmit 3 bytes at a rate of 1B/s", func() {
			response, err := http.Get("http://localhost:8081/slow")
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			tStart := time.Now()
			bytes, err := ioutil.ReadAll(response.Body)
			tEnd := time.Now()
			tDuration := tEnd.Sub(tStart)
			Expect(string(bytes)).To(Equal(body))
			Expect(err).To(BeNil())
			Expect(tDuration.Round(time.Second)).To(Equal(3 * time.Second))

		})
	})

	Context("An endpoint with a data rate larger than the message", func() {
		body := "123"
		rate, _ := datarate.Parse("10B/s")
		config := &httpResponder.Config{Port: 8081, Responses: map[string]*httpResponder.Method{
			"/slow": {Get: &httpResponder.ResponseOptions{Static: &httpResponder.Response{
				Rate:   rate,
				Status: 200,
				Body:   &body,
			}}},
		}}
		httpServer := config.NewServer()

		BeforeEach(func() {
			httpServer.Run()
			time.Sleep(1 * time.Second)
		})

		AfterEach(func() {
			httpServer.Stop()
		})

		It("should deliver the message completely and immediately", func() {
			response, err := http.Get("http://localhost:8081/slow")
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			tStart := time.Now()
			bytes, err := ioutil.ReadAll(response.Body)
			tEnd := time.Now()
			tDuration := tEnd.Sub(tStart)
			Expect(string(bytes)).To(Equal(body))
			Expect(err).To(BeNil())
			Expect(tDuration.Round(time.Second)).To(Equal(0 * time.Second))
		})
	})

})
