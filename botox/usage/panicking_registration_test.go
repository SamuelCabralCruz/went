//go:build test

package usage_test

import (
	"errors"
	"github.com/SamuelCabralCruz/going/botox"
	"github.com/SamuelCabralCruz/going/botox/usage/fixture"
	. "github.com/SamuelCabralCruz/going/kinggo"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = DescribeFunction(botox.Resolve[any], func() {
	var panickedError error

	BeforeEach(func() {
		panickedError = errors.New("something went wrong")
	})

	AfterEach(func() {
		botox.Reset()
	})

	Context("with non singleton registration", func() {
		var observedInstance fixture.Stateless
		var observedError error

		act := func() {
			observedInstance, observedError = botox.Resolve[fixture.Stateless]()
		}

		Context("with panicking supplier", func() {
			BeforeEach(func() {
				botox.RegisterSupplier(func() fixture.Stateless {
					panic(panickedError)
				})
			})

			It("should return recovered error", func() {
				act()

				Expect(observedInstance).To(BeZero())
				Expect(observedError).NotTo(BeNil())
				Expect(observedError).To(Equal(panickedError))
			})
		})

		Context("with panicking producer", func() {
			BeforeEach(func() {
				botox.RegisterProducer(func() (fixture.Stateless, error) {
					panic(panickedError)
				})
			})

			It("should return recovered error", func() {
				act()

				Expect(observedInstance).To(BeZero())
				Expect(observedError).NotTo(BeNil())
				Expect(observedError).To(Equal(panickedError))
			})
		})
	})

	Context("with singleton registration", func() {
		var observedInstance *fixture.Stateless
		var observedError error

		act := func() {
			observedInstance, observedError = botox.Resolve[*fixture.Stateless]()
		}

		Context("with panicking singleton supplier", func() {
			BeforeEach(func() {
				botox.RegisterSingletonSupplier(func() fixture.Stateless {
					panic(panickedError)
				})
			})

			It("should return recovered error", func() {
				act()

				Expect(observedInstance).To(BeZero())
				Expect(observedError).NotTo(BeNil())
				Expect(observedError).To(Equal(panickedError))
			})
		})

		Context("with panicking singleton producer", func() {
			BeforeEach(func() {
				botox.RegisterSingletonProducer(func() (fixture.Stateless, error) {
					panic(panickedError)
				})
			})

			It("should return recovered error", func() {
				act()

				Expect(observedInstance).To(BeZero())
				Expect(observedError).NotTo(BeNil())
				Expect(observedError).To(Equal(panickedError))
			})
		})
	})
})
