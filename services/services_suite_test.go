package services_test

import (
	"Dp218Go/services"
	"Dp218Go/services/mock"
	"errors"
	"github.com/golang/mock/gomock"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestServices(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Services Suite")
}

var _ = Describe(".order Create", func() {
	var (
		order         *services.OrderService
		mockCtrl      *gomock.Controller
		repoOrder     *mock.MockOrderRepo
		expectedError error
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		repoOrder = mock.NewMockOrderRepo(mockCtrl)
		order = &services.OrderService{RepoOrder: repoOrder}
		expectedError = errors.New("expectedError")
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})
	Context("when DeleteOrder is correct", func() {
		BeforeEach(func() {
			repoOrder.EXPECT().DeleteOrder(1).Return(nil)
		})
		It("should return correct", func() {
			Expect(order.DeleteOrder(1)).To(Succeed())
		})
	})

	Context("when DeleteOrder return err", func() {
		BeforeEach(func() {
			repoOrder.EXPECT().DeleteOrder(1).Return(expectedError)
		})
		It("should return correct", func() {
			Expect(order.DeleteOrder(1)).To(MatchError("expectedError"))
		})
	})
})
