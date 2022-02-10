package routing_test

import (
	"fmt"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestRouting(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Routing Suite")
}

var _ = BeforeSuite(func() {
	fmt.Println("Test suite is starting... Here we can declare our parameters")
})

var _ = AfterSuite(func() {
	fmt.Println("Test suite finished. We can add here some logic like reset parameters.")
})
