package routing_test

import (
	"Dp218Go/routing"
	"fmt"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"net/http"
	"net/http/httptest"
	"net/url"
)

var _ = Describe("Handlers", func() {
	var (
		req     *http.Request
		handler http.Handler
		writer  *httptest.ResponseRecorder
	)

	XDescribe("HandlePOST of ChooseScooter, but now pending", func() {
		BeforeEach(func() {
			fmt.Printf("namespace-%d\n", GinkgoParallelProcess())
			By("Registering a handler and a request")
			handler = http.HandlerFunc(routing.ChooseScooter)
			req = httptest.NewRequest(http.MethodPost, "/choose-scooter", nil)

			By("Setting a form value and creating a test writer")
			req.Form = url.Values{"id": {"2"}}
			writer = httptest.NewRecorder()
		})

		It("Processes a POST request successfully", func() {
			handler.ServeHTTP(writer, req)
			Expect(writer.Code).To(Equal(http.StatusOK))
		})

		Context("When we post the wrong data", func() {
			JustBeforeEach(func() {
				req.Form = url.Values{"id": {"noID"}}
			})
			It("Returns a bad request server error", func() {
				handler.ServeHTTP(writer, req)
				Expect(writer.Code).To(
					Equal(http.StatusBadRequest))
			})
		})
	})

	FDescribe("HandlePOST of ChooseStation in one thread with FOCUS mode", Serial, func() {
		BeforeEach(func() {
			handler = http.HandlerFunc(routing.ChooseStation)
			req = httptest.NewRequest(http.MethodPost, "/choose-station", nil)
			req.Form = url.Values{"id": {"12"}}
			writer = httptest.NewRecorder()
		})

		Context("when post correct data", func() {
			It("returns code 200", func() {
				handler.ServeHTTP(writer, req)
				Expect(writer.Code).To(Equal(http.StatusOK))
			})
		})

		Context("when post incorrect data", func() {
			JustBeforeEach(func() {
				req.Form = url.Values{"id": {"notID"}}
			})
			It("returns code 400", func() {
				handler.ServeHTTP(writer, req)
				Expect(writer.Code).To(Equal(http.StatusBadRequest))
			})
		})
	})
})
