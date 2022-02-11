package services_test

import (
	"Dp218Go/models"
	"Dp218Go/services"
	"Dp218Go/services/mock"
	"errors"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe(".orderService testing", func() {
	var (
		order         *services.OrderService
		mockCtrl      *gomock.Controller
		repoOrder     *mock.MockOrderRepo
		fakeOrder     models.Order
		fakeMileage   float64
		expectedError error
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		repoOrder = mock.NewMockOrderRepo(mockCtrl)
		order = &services.OrderService{RepoOrder: repoOrder}
		expectedError = errors.New("expectedError")
		fakeOrder = models.Order{ID: 1, UserID: 2, ScooterID: 3, StatusEndID: 2, StatusStartID: 1, Distance: 100.5, Amount: 777}
		fakeMileage = 233.33
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

	Context("when DeleteOrder returns an error", func() {
		BeforeEach(func() {
			repoOrder.EXPECT().DeleteOrder(1).Return(expectedError)
		})
		It("should return error and match it by value", func() {
			Expect(order.DeleteOrder(1)).To(MatchError("expectedError"))
		})
	})

	Context("when CreateOrder returns a correct result", func() {
		BeforeEach(func() {
			repoOrder.EXPECT().CreateOrder(models.User{}, 1, 1, 1, 1.0).Return(models.Order{}, nil)
		})
		It("should return correct", func() {
			Expect(order.CreateOrder(models.User{}, 1, 1, 1, 1.0)).Error().ShouldNot(HaveOccurred())
		})
	})

	Context("when CreateOrder returns an error", func() {
		BeforeEach(func() {
			repoOrder.EXPECT().CreateOrder(models.User{}, 1, 1, 1, 1.0).Return(models.Order{}, expectedError)
		})
		It("should return error", func() {
			Expect(order.CreateOrder(models.User{}, 1, 1, 1, 1.0)).Error().Should(HaveOccurred())
		})
	})

	Context("when CreateOrder returns an order model, and nil error", func() {
		BeforeEach(func() {
			repoOrder.EXPECT().CreateOrder(models.User{}, 1, 1, 1, 1.0).Return(fakeOrder, nil).AnyTimes()
		})
		It("should return an equal order to fakeOrder and nil error", func() {
			Expect(order.CreateOrder(models.User{}, 1, 1, 1, 1.0)).Should(BeEquivalentTo(fakeOrder))
			Expect(order.CreateOrder(models.User{}, 1, 1, 1, 1.0)).Error().ShouldNot(HaveOccurred())
		})
	})

	Context("when GetScooterMileageByID returns a correct mileage ", func() {
		JustBeforeEach(func() {
			repoOrder.EXPECT().GetScooterMileageByID(1).Return(fakeMileage, nil).AnyTimes()
		})
		It("should return a value is equal to fakeMileage", func() {
			Expect(order.GetScooterMileageByID(1)).Should(BeNumerically("==", fakeMileage))
			Expect(order.GetScooterMileageByID(1)).Error().ShouldNot(HaveOccurred())
		})
	})

	Context("when GetScooterMileageByID returns different value to expected mileage with nil error", func() {
		BeforeEach(func() {
			repoOrder.EXPECT().GetScooterMileageByID(1).Return(fakeMileage+1, nil).AnyTimes()
		})
		It("should pass with higher value than fakeMileage", func() {
			Expect(order.GetScooterMileageByID(1)).Should(BeNumerically(">", fakeMileage))
			_, err := order.GetScooterMileageByID(1)
			Expect(err).ShouldNot(HaveOccurred())
		})
	})
})
