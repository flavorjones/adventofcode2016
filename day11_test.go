package adventofcode2016_test

import (
	"bytes"
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"regexp"
	"sort"
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

type RTFArtifacts []RTFArtifact
type RTFConfig []RTFArtifacts

func (c RTFArtifacts) Len() int      { return len(c) }
func (c RTFArtifacts) Swap(j, k int) { c[j], c[k] = c[k], c[j] }
func (c RTFArtifacts) Less(j, k int) bool {
	if c[j].Element() == c[k].Element() {
		return c[j].Artifact() < c[k].Artifact()
	}
	return c[j].Element() < c[k].Element()
}

func NewRTFConfig(config ...RTFArtifacts) RTFConfig {
	rtfc := make(RTFConfig, len(config))
	for j, _ := range rtfc {
		rtfc[j] = config[j]
		sort.Sort(rtfc[j])
	}
	return rtfc
}

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

type RTFHistory []RadioisotopeTestingFacility
type RTFPermutations []RadioisotopeTestingFacility

func elevatorPermutations(nElements int) [][]int {
	rval := [][]int{}
	for j := 0; j < nElements; j++ {
		rval = append(rval, []int{j})
		for k := j + 1; k < nElements; k++ {
			rval = append(rval, []int{j, k})
		}
	}
	return rval
}

func (rtf RadioisotopeTestingFacility) Permutations() RTFPermutations {
	permutations := RTFPermutations{}

	for _, indexPermutation := range elevatorPermutations(len(rtf.config[rtf.ePos])) {
		artifacts := RTFArtifacts{}
		for _, index := range indexPermutation {
			artifacts = append(artifacts, rtf.config[rtf.ePos][index])
		}

		modifiedFloor := make(RTFArtifacts, len(rtf.config[rtf.ePos]))
		copy(modifiedFloor, rtf.config[rtf.ePos])

		for j := len(indexPermutation) - 1; j >= 0; j-- {
			index := indexPermutation[j]
			if index < len(modifiedFloor)-1 {
				modifiedFloor = append(
					modifiedFloor[:index],
					modifiedFloor[index+1:]...,
				)
			} else {
				modifiedFloor = modifiedFloor[:index]
			}
		}

		if rtf.ePos > 0 {
			newPos := rtf.ePos - 1
			permutation := NewRTFConfig(rtf.config...)
			permutation[newPos] = append(permutation[newPos], artifacts...)
			permutation[rtf.ePos] = modifiedFloor
			permutations = append(permutations,
				RadioisotopeTestingFacility{permutation, newPos})
		}

		if rtf.ePos < len(rtf.config)-1 {
			newPos := rtf.ePos + 1
			permutation := NewRTFConfig(rtf.config...)
			permutation[newPos] = append(permutation[newPos], artifacts...)
			permutation[rtf.ePos] = modifiedFloor
			permutations = append(permutations,
				RadioisotopeTestingFacility{permutation, newPos})
		}
	}

	return permutations
}

func (rtf RadioisotopeTestingFacility) ValidPermutations() RTFPermutations {
	permutations := rtf.Permutations()
	validPermutations := make(RTFPermutations, 0, len(permutations))
	for j, _ := range permutations {
		if permutations[j].ok() {
			validPermutations = append(validPermutations, permutations[j])
		}
	}
	return validPermutations
}

func RTFTripPlanImpl(stateHistory RTFHistory) (RTFHistory, bool) {
	return append(stateHistory, NewRadioisotopeTestingFacility(RTFConfig{})), true
}

func RTFTripPlan(config RTFConfig) []RadioisotopeTestingFacility {
	history := RTFHistory{NewRadioisotopeTestingFacility(config)}
	history, _ = RTFTripPlanImpl(history)
	return history
}

var _ = Describe("Day11", func() {
	testData := `The first floor contains a hydrogen-compatible microchip and a lithium-compatible microchip.
The second floor contains a hydrogen generator.
The third floor contains a lithium generator.
The fourth floor contains nothing relevant.`
	testSetup := strings.Split(string(testData), "\n")

	Describe("RTFGenerator", func() {
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

	Describe("RTFMicrochip", func() {
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

	Describe("NewRTFConfig()", func() {
		It("creates a new config in canonical (sorted) order", func() {
			actual := NewRTFConfig(
				RTFArtifacts{RTFGenerator{"a"}, RTFMicrochip{"a"}},
				RTFArtifacts{RTFMicrochip{"b"}, RTFGenerator{"z"}, RTFGenerator{"b"}},
				RTFArtifacts{RTFMicrochip{"z"}, RTFMicrochip{"a"}},
				RTFArtifacts{RTFGenerator{"z"}, RTFGenerator{"a"}},
			)
			canonical := RTFConfig{
				RTFArtifacts{RTFGenerator{"a"}, RTFMicrochip{"a"}},
				RTFArtifacts{RTFGenerator{"b"}, RTFMicrochip{"b"}, RTFGenerator{"z"}},
				RTFArtifacts{RTFMicrochip{"a"}, RTFMicrochip{"z"}},
				RTFArtifacts{RTFGenerator{"a"}, RTFGenerator{"z"}},
			}
			Expect(actual).To(Equal(canonical))
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

	Describe("RadioisotopeTestingFacility", func() {
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

		Describe("#Permutations", func() {
			It("returns all possible next-steps", func() {
				initial := RadioisotopeTestingFacility{
					NewRTFConfig(
						RTFArtifacts{},
						RTFArtifacts{RTFMicrochip{"a"}, RTFMicrochip{"b"}, RTFMicrochip{"c"}},
						RTFArtifacts{},
						RTFArtifacts{},
					), 1}

				permutations := RTFPermutations{
					RadioisotopeTestingFacility{
						NewRTFConfig(
							RTFArtifacts{RTFMicrochip{"a"}},
							RTFArtifacts{RTFMicrochip{"b"}, RTFMicrochip{"c"}},
							RTFArtifacts{},
							RTFArtifacts{},
						), 0},
					RadioisotopeTestingFacility{
						NewRTFConfig(
							RTFArtifacts{},
							RTFArtifacts{RTFMicrochip{"b"}, RTFMicrochip{"c"}},
							RTFArtifacts{RTFMicrochip{"a"}},
							RTFArtifacts{},
						), 2},

					RadioisotopeTestingFacility{
						NewRTFConfig(
							RTFArtifacts{RTFMicrochip{"a"}, RTFMicrochip{"b"}},
							RTFArtifacts{RTFMicrochip{"c"}},
							RTFArtifacts{},
							RTFArtifacts{},
						), 0},
					RadioisotopeTestingFacility{
						NewRTFConfig(
							RTFArtifacts{},
							RTFArtifacts{RTFMicrochip{"c"}},
							RTFArtifacts{RTFMicrochip{"a"}, RTFMicrochip{"b"}},
							RTFArtifacts{},
						), 2},

					RadioisotopeTestingFacility{
						NewRTFConfig(
							RTFArtifacts{RTFMicrochip{"a"}, RTFMicrochip{"c"}},
							RTFArtifacts{RTFMicrochip{"b"}},
							RTFArtifacts{},
							RTFArtifacts{},
						), 0},
					RadioisotopeTestingFacility{
						NewRTFConfig(
							RTFArtifacts{},
							RTFArtifacts{RTFMicrochip{"b"}},
							RTFArtifacts{RTFMicrochip{"a"}, RTFMicrochip{"c"}},
							RTFArtifacts{},
						), 2},

					RadioisotopeTestingFacility{
						NewRTFConfig(
							RTFArtifacts{RTFMicrochip{"b"}},
							RTFArtifacts{RTFMicrochip{"a"}, RTFMicrochip{"c"}},
							RTFArtifacts{},
							RTFArtifacts{},
						), 0},
					RadioisotopeTestingFacility{
						NewRTFConfig(
							RTFArtifacts{},
							RTFArtifacts{RTFMicrochip{"a"}, RTFMicrochip{"c"}},
							RTFArtifacts{RTFMicrochip{"b"}},
							RTFArtifacts{},
						), 2},

					RadioisotopeTestingFacility{
						NewRTFConfig(
							RTFArtifacts{RTFMicrochip{"b"}, RTFMicrochip{"c"}},
							RTFArtifacts{RTFMicrochip{"a"}},
							RTFArtifacts{},
							RTFArtifacts{},
						), 0},
					RadioisotopeTestingFacility{
						NewRTFConfig(
							RTFArtifacts{},
							RTFArtifacts{RTFMicrochip{"a"}},
							RTFArtifacts{RTFMicrochip{"b"}, RTFMicrochip{"c"}},
							RTFArtifacts{},
						), 2},
					RadioisotopeTestingFacility{
						NewRTFConfig(
							RTFArtifacts{RTFMicrochip{"c"}},
							RTFArtifacts{RTFMicrochip{"a"}, RTFMicrochip{"b"}},
							RTFArtifacts{},
							RTFArtifacts{},
						), 0},
					RadioisotopeTestingFacility{
						NewRTFConfig(
							RTFArtifacts{},
							RTFArtifacts{RTFMicrochip{"a"}, RTFMicrochip{"b"}},
							RTFArtifacts{RTFMicrochip{"c"}},
							RTFArtifacts{},
						), 2},
				}

				actual := initial.Permutations()
				Expect(len(actual)).To(Equal(len(permutations)))
				Expect(actual).To(ConsistOf(permutations))
			})
		})

		Describe("#ValidPermutations", func() {
			It("returns all valid next-steps", func() {
				initial := RadioisotopeTestingFacility{
					NewRTFConfig(
						RTFArtifacts{},
						RTFArtifacts{RTFMicrochip{"a"}, RTFGenerator{"a"}, RTFMicrochip{"c"}},
						RTFArtifacts{},
						RTFArtifacts{},
					), 1}

				permutations := RTFPermutations{
					RadioisotopeTestingFacility{
						NewRTFConfig(
							RTFArtifacts{RTFMicrochip{"a"}, RTFGenerator{"a"}},
							RTFArtifacts{RTFMicrochip{"c"}},
							RTFArtifacts{},
							RTFArtifacts{},
						), 0},
					RadioisotopeTestingFacility{
						NewRTFConfig(
							RTFArtifacts{RTFGenerator{"a"}},
							RTFArtifacts{RTFMicrochip{"a"}, RTFMicrochip{"c"}},
							RTFArtifacts{},
							RTFArtifacts{},
						), 0},
					RadioisotopeTestingFacility{
						NewRTFConfig(
							RTFArtifacts{},
							RTFArtifacts{RTFMicrochip{"c"}},
							RTFArtifacts{RTFMicrochip{"a"}, RTFGenerator{"a"}},
							RTFArtifacts{},
						), 2},
					RadioisotopeTestingFacility{
						NewRTFConfig(
							RTFArtifacts{},
							RTFArtifacts{RTFMicrochip{"a"}, RTFMicrochip{"c"}},
							RTFArtifacts{RTFGenerator{"a"}},
							RTFArtifacts{},
						), 2},

					RadioisotopeTestingFacility{
						NewRTFConfig(
							RTFArtifacts{RTFMicrochip{"c"}},
							RTFArtifacts{RTFGenerator{"a"}, RTFMicrochip{"a"}},
							RTFArtifacts{},
							RTFArtifacts{},
						), 0},

					RadioisotopeTestingFacility{
						NewRTFConfig(
							RTFArtifacts{RTFMicrochip{"a"}, RTFMicrochip{"c"}},
							RTFArtifacts{RTFGenerator{"a"}},
							RTFArtifacts{},
							RTFArtifacts{},
						), 0},
					RadioisotopeTestingFacility{
						NewRTFConfig(
							RTFArtifacts{},
							RTFArtifacts{RTFGenerator{"a"}},
							RTFArtifacts{RTFMicrochip{"a"}, RTFMicrochip{"c"}},
							RTFArtifacts{},
						), 2},
					RadioisotopeTestingFacility{
						NewRTFConfig(
							RTFArtifacts{},
							RTFArtifacts{RTFMicrochip{"a"}, RTFGenerator{"a"}},
							RTFArtifacts{RTFMicrochip{"c"}},
							RTFArtifacts{},
						), 2},
				}

				actual := initial.ValidPermutations()
				Expect(actual).To(ConsistOf(permutations))
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
