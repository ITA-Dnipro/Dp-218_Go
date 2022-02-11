package services_test

import (
	"Dp218Go/models"
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

	Describe(".scooterService testing", func() {
		type expectedResult struct {
			status models.ScooterStatus
			err    error
		}

		var (
			scooterService *services.ScooterService
			repoScooter    *mock.MockScooterRepo
		)

		BeforeEach(func() {
			repoScooter = mock.NewMockScooterRepo(mockCtrl)
			scooterService = &services.ScooterService{RepoScooter: repoScooter}
			expectedError = errors.New("expectedError")

		})

		AfterEach(func() {
			mockCtrl.Finish()
		})

		DescribeTable(".scooterService GetScooterStatus table testing",
			func(id int, result expectedResult) {
				repoScooter.EXPECT().GetScooterStatus(id).Return(result.status, result.err).AnyTimes()
				status, err := scooterService.GetScooterStatus(id)
				Expect(status).Should(BeEquivalentTo(result.status))
				Expect(err).Should(MatchError(result.err))
			},
			Entry("should work with noError", 10, expectedResult{status: models.ScooterStatus{Scooter: models.ScooterDTO{}, Location: models.Coordinate{}, BatteryRemain: 22.2, StationID: 4}, err: errors.New("noError")}),
			Entry("should work with someError", -1, expectedResult{status: models.ScooterStatus{Scooter: models.ScooterDTO{}, Location: models.Coordinate{}, BatteryRemain: 22.2, StationID: 4}, err: errors.New("someError")}),
		)
	})
})
