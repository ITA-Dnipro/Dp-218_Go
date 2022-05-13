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

var _ = Describe("ScooterGinkgo", func() {
	type expectedResult struct {
		status models.ScooterStatus
		err    error
	}

	var (
		scooterService *services.ScooterService
		repoScooter    *mock.MockScooterRepo
		mockCtrl       *gomock.Controller
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		repoScooter = mock.NewMockScooterRepo(mockCtrl)
		scooterService = &services.ScooterService{RepoScooter: repoScooter}

	})

	DescribeTable(".scooterService GetScooterStatus table testing",
		func(id int, result expectedResult) {
			By("Mocking GetScooterStatus")
			repoScooter.EXPECT().GetScooterStatus(id).Return(result.status, result.err).AnyTimes()

			By("Calling GetScooterStatus and comparing results")
			status, err := scooterService.GetScooterStatus(id)
			Expect(status).Should(BeEquivalentTo(result.status))
			Expect(err).Should(MatchError(result.err))
		},
		EntryDescription("parameters: ID=%v, ExpectedResult:%v"),
		Entry(nil, 10, expectedResult{status: models.ScooterStatus{Scooter: models.ScooterDTO{}, Location: models.Coordinate{}, BatteryRemain: 22.2, StationID: 4}, err: errors.New("noError")}),
		Entry(nil, -1, expectedResult{status: models.ScooterStatus{Scooter: models.ScooterDTO{}, Location: models.Coordinate{}, BatteryRemain: 22.2, StationID: 4}, err: errors.New("someError")}),
	)
})
