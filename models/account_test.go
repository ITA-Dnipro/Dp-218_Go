package models_test

import (
	"Dp218Go/models"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("AccountTransaction", func() {
	var transaction *models.AccountTransaction

	BeforeEach(func() {
		transaction = &models.AccountTransaction{AmountCents: 0}
	})

	Context("has zero amount", func() {
		It("should return money with 0 dollars 0 cents", func() {
			money := transaction.GetAmountInMoney()
			Expect(money.Dollars).To(BeZero())
			Expect(money.Cents).To(BeZero())
		})
	})

	Context("with 100 cents in amount", func() {
		JustBeforeEach(func() {
			transaction.AmountCents = 100
		})
		Specify("money Dollars = 1, cents = 0", func() {
			money := transaction.GetAmountInMoney()
			Expect(money.Dollars).To(BeNumerically("==", 1))
			Expect(money.Cents).To(BeZero())
		})
	})

	Context("1 cent was added to 101 cents", func() {
		var startMoney models.Money

		BeforeEach(func() {
			transaction.AmountCents = 101
			startMoney = transaction.GetAmountInMoney()
		})
		Context("nothing added", func() {
			It("has the same money", func() {
				Expect(transaction.GetAmountInMoney().Cents).To(Equal(startMoney.Cents))
				Expect(transaction.GetAmountInMoney().Dollars).To(Equal(startMoney.Dollars))
			})
		})

		When("1 more cent was added", func() {
			It("should be one more cent", func() {
				transaction.ChangeMoneyInTransaction(1)
				Expect(transaction.GetAmountInMoney().Cents).To(Equal(startMoney.Cents + 1))
			})
		})
	})
})
