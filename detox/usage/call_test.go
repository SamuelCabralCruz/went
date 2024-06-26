//go:build test

package usage_test

import (
	"fmt"
	"github.com/SamuelCabralCruz/going/detox"
	"github.com/SamuelCabralCruz/going/detox/usage/fixture"
	. "github.com/SamuelCabralCruz/going/kinggo"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = DescribeType[detox.Detox[any]](func() {
	var cut fixture.Interface1Mock
	var mocked detox.Mocked[func(string) string]

	BeforeEach(func() {
		cut = fixture.NewInterface1Mock()
		mocked = detox.When(cut.Detox, cut.SingleArgSingleReturn)
	})

	AfterEach(func() {
		cut.Reset()
	})

	DescribeFunction(detox.Mocked[any].Call, func() {
		var observed string

		act := func() {
			mocked.Call(func(s string) string {
				return fmt.Sprintf("this return has been mocked - %s", s)
			})
			observed = cut.SingleArgSingleReturn("some input value")
		}

		It("should resolve fake implementation", func() {
			act()

			Expect(observed).To(Equal("this return has been mocked - some input value"))
		})

		It("should be persistent", func() {
			act()

			Expect(func() { cut.SingleArgSingleReturn("1st additional invocation") }).NotTo(Panic())
			Expect(func() { cut.SingleArgSingleReturn("2nd additional invocation") }).NotTo(Panic())
			Expect(func() { cut.SingleArgSingleReturn("3rd additional invocation") }).NotTo(Panic())
			Expect(func() { cut.SingleArgSingleReturn("4th additional invocation") }).NotTo(Panic())
			Expect(func() { cut.SingleArgSingleReturn("5th additional invocation") }).NotTo(Panic())
		})

		Context("with already registered persistent implementation", func() {
			BeforeEach(func() {
				mocked.Call(func(_ string) string {
					return "already registered"
				})
			})

			It("should override the previous implementation", func() {
				act()

				Expect(observed).To(Equal("this return has been mocked - some input value"))
			})
		})
	})
})
