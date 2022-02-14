package routing_test

import (
	"Dp218Go/routing"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestRouting(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Routing Suite")
}

var _ = Describe("Handlers", func() {
	var (
		req     *http.Request
		handler http.Handler
		writer  *httptest.ResponseRecorder
	)
	Describe("HandlePOST", func() {
		BeforeEach(func() {
			handler = http.HandlerFunc(routing.ChooseScooter)
			req = httptest.NewRequest(http.MethodPost, "/choose-scooter", nil)
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
})
