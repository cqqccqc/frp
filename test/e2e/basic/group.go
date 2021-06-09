package basic

import (
	"github.com/fatedier/frp/test/e2e/framework"

	. "github.com/onsi/ginkgo"
)

var _ = Describe("[Feature: Group]", func() {
	f := framework.NewDefaultFramework()

	It("Load Balancing by group", func() {
		// TODO
		_ = f
	})

	Describe("Health Check", func() {
		It("TCP", func() {
			// TODO
		})
		It("HTTP", func() {
			// TODO
		})
	})
})
