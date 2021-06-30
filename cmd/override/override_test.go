package override

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestOverride(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Override Suite")
}

var (
	res1 = Resource{name: "rone"}
	res2 = Resource{name: "rtwo"}
)

var (
	comp1 = Component{name: "cone", resources: Resources{res1}}
	comp2 = Component{name: "ctwo", resources: Resources{res1, res2}}
)

var comps = Components{comp1, comp2}

var _ = Describe("Component", func() {
	It("should have name of cone", func() {
		Expect(comp1.Name()).To(Equal("cone"))
	})

	It("should have a resource named rone", func() {
		Expect(comp1.Resource("rone")).To(Equal(res1))
	})

	It("should return strings rone and rtwo", func() {
		Expect(comp1.ListResources()).To(Equal([]string{res1.name}))
	})

	It("should not have a resource named rthree", func() {
		_, err := comp1.Resource("rthree")
		Expect(err).ToNot(BeNil())
	})
})

var _ = Describe("Components", func() {

	It("should have 2 components", func() {
		Expect(len(comps)).To(Equal(2))
	})

	It("should return strings cone and ctwo", func() {
		Expect(comps.Strings()).To(Equal([]string{comp1.name, comp2.name}))
	})

	It("should have a component named 'cone'", func() {
		Expect(comps.Component("cone")).To(Equal(comp1))
	})

	It("should not have a component named 'cthree'", func() {
		_, err := comps.Component("cthree")
		Expect(err).ToNot(BeNil())
	})
})
