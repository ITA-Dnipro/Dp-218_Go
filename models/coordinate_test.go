package models_test

import (
	"Dp218Go/models"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Coordinate", func() {
	var coordinate1, coordinate2 *models.Coordinate
	var actualDistance float64

	BeforeEach(func() {
		coordinate1 = &models.Coordinate{
			Latitude:  40.5,
			Longitude: 45.5,
		}

		coordinate2 = &models.Coordinate{
			Latitude:  55.5,
			Longitude: 50.0,
		}

		actualDistance = coordinate1.Distance(*coordinate2)
	})

	Describe("Calculate distance between two coordinates", func() {
		Context("has a distance as counted before", func() {
			It("should be an equal", func() {
				Expect(coordinate1.Distance(*coordinate2)).To(Equal(actualDistance))
			})
		})
	})
})
