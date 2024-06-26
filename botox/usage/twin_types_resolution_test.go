//go:build test

package usage_test

import (
	"github.com/SamuelCabralCruz/going/botox"
	"github.com/SamuelCabralCruz/going/botox/usage/fixture/pkg1"
	"github.com/SamuelCabralCruz/going/botox/usage/fixture/pkg2"
	. "github.com/SamuelCabralCruz/going/kinggo"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = DescribeFunction(botox.MustResolve[any], func() {
	var observedTwin1 pkg1.SomeStruct
	var observedTwin2 pkg2.SomeStruct

	act := func() {
		observedTwin1 = botox.MustResolve[pkg1.SomeStruct]()
		observedTwin2 = botox.MustResolve[pkg2.SomeStruct]()
	}

	AfterEach(func() {
		botox.Reset()
	})

	Context("with twin types registrations", func() {
		var twin1 pkg1.SomeStruct
		var twin2 pkg2.SomeStruct

		BeforeEach(func() {
			twin1 = pkg1.SomeStruct{}
			twin2 = pkg2.SomeStruct{}
			botox.RegisterInstance(twin1)
			botox.RegisterInstance(twin2)
		})

		It("should resolve types specifically", func() {
			act()

			Expect(observedTwin1).To(Equal(twin1))
			Expect(observedTwin2).To(Equal(twin2))
		})
	})
})
