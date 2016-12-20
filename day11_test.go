package adventofcode2016_test

import (
	"bytes"
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"regexp"
	"strings"
)

type RTFArtifact interface {
	Element() string
	Artifact() string
}

type RTFGenerator struct {
	element string
}

func (rtfg RTFGenerator) Element() string {
	return rtfg.element
}

func (rtfg RTFGenerator) Artifact() string {
	return "generator"
}

func (rtfg RTFGenerator) String() string {
	return fmt.Sprintf("%.2s-%.1s", rtfg.Element(), rtfg.Artifact())
}

type RTFMicrochip struct {
	element string
}

func (rtfm RTFMicrochip) Element() string {
	return rtfm.element
}

func (rtfm RTFMicrochip) Artifact() string {
	return "microchip"
}

func (rtfm RTFMicrochip) String() string {
	return fmt.Sprintf("%.2s-%.2s", rtfm.Element(), rtfm.Artifact())
}

type RTFConfig [][]RTFArtifact

func (rtfc RTFConfig) String() string {
	var output bytes.Buffer
	for floor := len(rtfc) - 1; floor >= 0; floor-- {
		output.WriteString(fmt.Sprintf("f(%d): %s\n", floor+1, rtfc[floor]))
	}
	return output.String()
}

func RTFConfigRead(setup []string) RTFConfig {
	chipRe := regexp.MustCompile(`\b(\w+)-compatible microchip`)
	genRe := regexp.MustCompile(`\b(\w+) generator`)
	rtfc := make(RTFConfig, len(setup))

	for floor, description := range setup {
		rtfc[floor] = []RTFArtifact{}

		chipMatches := chipRe.FindAllStringSubmatch(description, -1)
		for _, chipMatch := range chipMatches {
			rtfc[floor] = append(rtfc[floor], RTFMicrochip{chipMatch[1]})
		}

		genMatches := genRe.FindAllStringSubmatch(description, -1)
		for _, genMatch := range genMatches {
			rtfc[floor] = append(rtfc[floor], RTFGenerator{genMatch[1]})
		}
	}
	return rtfc
}

type RadioisotopeTestingFacility struct {
	config RTFConfig
	ePos   int
}

func NewRadioisotopeTestingFacility(config RTFConfig) RadioisotopeTestingFacility {
	return RadioisotopeTestingFacility{
		config,
		0,
	}
}

func (rtf *RadioisotopeTestingFacility) ok() bool {
	for _, floor := range rtf.config {
	artifact_loop:
		for _, artifact := range floor {
			if artifact.Artifact() != "microchip" {
				continue
			}

			for _, other := range floor {
				_, ok := other.(RTFGenerator)
				if ok && other.Element() == artifact.Element() {
					continue artifact_loop
				}
			}

			for _, other := range floor {
				_, ok := other.(RTFGenerator)
				if ok && other.Element() != artifact.Element() {
					return false
				}
			}
		}
	}
	return true
}

func RTFTripPlan(config []RTFConfig) []RTFConfig {
	
}


var _ = Describe("Day11", func() {
	testData := `The first floor contains a hydrogen-compatible microchip and a lithium-compatible microchip.
The second floor contains a hydrogen generator.
The third floor contains a lithium generator.
The fourth floor contains nothing relevant.`
	testSetup := strings.Split(string(testData), "\n")

	Describe("RTFG", func() {
		rtfg := RTFGenerator{"polonium"}

		Describe("#Element()", func() {
			It("returns the element name", func() {
				Expect(rtfg.Element()).To(Equal("polonium"))
			})
		})

		Describe("#Artifact()", func() {
			It("returns the string 'generator'", func() {
				Expect(rtfg.Artifact()).To(Equal("generator"))
			})
		})
	})

	Describe("RTFM", func() {
		rtfg := RTFMicrochip{"polonium"}

		Describe("#Element()", func() {
			It("returns the element name", func() {
				Expect(rtfg.Element()).To(Equal("polonium"))
			})
		})

		Describe("#Artifact()", func() {
			It("returns the string 'microchip'", func() {
				Expect(rtfg.Artifact()).To(Equal("microchip"))
			})
		})
	})

	Describe("RTFConfigRead", func() {
		It("generates a config", func() {
			actual := RTFConfigRead(testSetup)
			expected := RTFConfig{
				[]RTFArtifact{RTFMicrochip{"hydrogen"}, RTFMicrochip{"lithium"}},
				[]RTFArtifact{RTFGenerator{"hydrogen"}},
				[]RTFArtifact{RTFGenerator{"lithium"}},
				[]RTFArtifact{},
			}
			Expect(actual).To(Equal(expected))
		})
	})

	Describe("RTF", func() {
		Describe("#ok", func() {
			Context("all chips are with their generators", func() {
				rtf := NewRadioisotopeTestingFacility(RTFConfig{
					[]RTFArtifact{RTFGenerator{"a"}, RTFMicrochip{"a"}},
					[]RTFArtifact{RTFGenerator{"b"}, RTFMicrochip{"b"}},
					[]RTFArtifact{RTFGenerator{"c"}, RTFMicrochip{"c"}, RTFGenerator{"d"}, RTFMicrochip{"d"}},
					[]RTFArtifact{},
				})

				It("is ok", func() {
					Expect(rtf.ok()).To(BeTrue())
				})
			})

			Context("all isolated chips not near a generator", func() {
				rtf := NewRadioisotopeTestingFacility(RTFConfig{
					[]RTFArtifact{RTFGenerator{"a"}},
					[]RTFArtifact{RTFMicrochip{"b"}},
					[]RTFArtifact{RTFGenerator{"c"}, RTFMicrochip{"c"}},
					[]RTFArtifact{},
				})

				It("is ok", func() {
					Expect(rtf.ok()).To(BeTrue())
				})
			})

			Context("an isolated chip is near a generator", func() {
				rtf := NewRadioisotopeTestingFacility(RTFConfig{
					[]RTFArtifact{RTFGenerator{"a"}},
					[]RTFArtifact{RTFMicrochip{"b"}, RTFGenerator{"a"}},
					[]RTFArtifact{RTFGenerator{"c"}, RTFMicrochip{"c"}},
					[]RTFArtifact{},
				})

				It("is not ok", func() {
					Expect(rtf.ok()).To(BeFalse())
				})
			})
		})
	})

	Describe("the test", func() {
		It("finds a solution", func() {
			config := RTFConfigRead(testSetup)
			solution := RTFTripPlan(config)
			Expect(len(solution)).To(Equal(11))
		})
	})

	// 	Describe("the puzzle", func() {
	// 		data := `The first floor contains a polonium generator, a thulium generator, a thulium-compatible microchip, a promethium generator, a ruthenium generator, a ruthenium-compatible microchip, a cobalt generator, and a cobalt-compatible microchip.
	// The second floor contains a polonium-compatible microchip and a promethium-compatible microchip.
	// The third floor contains nothing relevant.
	// The fourth floor contains nothing relevant.`
	// 		setup := strings.Split(string(data), "\n")

	// 	})
})
