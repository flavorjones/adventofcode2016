package adventofcode2016_test

import (
	"bytes"
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"reflect"
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
	return fmt.Sprintf("%.2s-%.1s", rtfm.Element(), rtfm.Artifact())
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

func (c RTFArtifacts) Remove(artifacts RTFArtifacts) RTFArtifacts {
	var rval RTFArtifacts
original:
	for _, original := range c {
		for _, artifact := range artifacts {
			if reflect.DeepEqual(original, artifact) {
				continue original
			}
		}
		rval = append(rval, original)
	}
	return rval
}

func (rtfc RTFConfig) Canonicalize() {
	for j, _ := range rtfc {
		sort.Sort(rtfc[j])
	}
}

func NewRTFConfig(config ...RTFArtifacts) RTFConfig {
	rtfc := make(RTFConfig, len(config))
	for j, _ := range config {
		rtfc[j] = make(RTFArtifacts, len(config[j]))
		copy(rtfc[j], config[j])
	}
	rtfc.Canonicalize()
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

func (rtf RadioisotopeTestingFacility) String() string {
	var output bytes.Buffer
	output.WriteString("\n")
	for floor := len(rtf.config) - 1; floor >= 0; floor-- {
		var indicator string
		if rtf.ePos == floor {
			indicator = "*"
		} else {
			indicator = " "
		}
		output.WriteString(fmt.Sprintf("     f(%d): %s %s\n", floor+1, indicator, rtf.config[floor]))
	}
	return output.String()
}

func (rtf RadioisotopeTestingFacility) Equals(rhs RadioisotopeTestingFacility) bool {
	if rtf.ePos != rhs.ePos || len(rtf.config) != len(rhs.config) {
		return false
	}
	for j, _ := range rtf.config {
		if len(rtf.config[j]) != len(rhs.config[j]) {
			return false
		}
		for k, _ := range rtf.config[j] {
			if rtf.config[j][k] != rhs.config[j][k] {
				return false
			}
		}
	}
	return true
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

func (rtfh RTFHistory) Contains(state RadioisotopeTestingFacility) bool {
	for _, rtf := range rtfh {
		if state.Equals(rtf) {
			return true
		}
	}
	return false
}

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

		if rtf.ePos > 0 {
			newPos := rtf.ePos - 1
			permutedConfig := NewRTFConfig(rtf.config...)
			permutedConfig[newPos] = append(permutedConfig[newPos], artifacts...)
			permutedConfig[rtf.ePos] = permutedConfig[rtf.ePos].Remove(artifacts)
			permutedConfig.Canonicalize()
			permutations = append(permutations,
				RadioisotopeTestingFacility{permutedConfig, newPos})
		}

		if rtf.ePos < len(rtf.config)-1 {
			newPos := rtf.ePos + 1
			permutedConfig := NewRTFConfig(rtf.config...)
			permutedConfig[newPos] = append(permutedConfig[newPos], artifacts...)
			permutedConfig[rtf.ePos] = permutedConfig[rtf.ePos].Remove(artifacts)
			permutedConfig.Canonicalize()
			permutations = append(permutations,
				RadioisotopeTestingFacility{permutedConfig, newPos})
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

func (rtf RadioisotopeTestingFacility) done() bool {
	rtfc := rtf.config
	return len(rtfc[0]) == 0 &&
		len(rtfc[1]) == 0 &&
		len(rtfc[2]) == 0 &&
		len(rtfc[3]) > 0
}

func RTFTripPlanImpl(stateHistory RTFHistory, maxDepth int) (RTFHistory, bool) {
	current := stateHistory[len(stateHistory)-1]
	permutations := current.ValidPermutations()
	for _, permutation := range permutations {
		if stateHistory.Contains(permutation) {
			continue
		}
		if permutation.done() {
			return append(stateHistory, permutation), true
		}

		if len(stateHistory) < maxDepth {
			fullStateHistory, ok := RTFTripPlanImpl(append(stateHistory, permutation), maxDepth)
			if ok {
				return fullStateHistory, true
			}
		}
	}

	return stateHistory, false
}

func RTFTripPlan(config RTFConfig) []RadioisotopeTestingFacility {
	// omg so inefficient, I'm embarassed but I'm ready to move onto the next puzzle.
	for depth := 1; ; depth++ {
		fmt.Println("MIKE: ---------- depth", depth, "----------")
		start := RTFHistory{NewRadioisotopeTestingFacility(config)}
		win, done := RTFTripPlanImpl(start, depth)
		if done {
			return win
		}
	}
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
		Describe("#done", func() {
			Context("everything's not on the fourth floor", func() {
				It("returns false", func() {
					rtf := NewRadioisotopeTestingFacility(RTFConfig{
						[]RTFArtifact{RTFGenerator{"a"}, RTFMicrochip{"a"}},
						[]RTFArtifact{RTFGenerator{"b"}, RTFMicrochip{"b"}},
						[]RTFArtifact{RTFGenerator{"c"}, RTFMicrochip{"c"}, RTFGenerator{"d"}, RTFMicrochip{"d"}},
						[]RTFArtifact{},
					})
					Expect(rtf.done()).To(BeFalse())
				})

				It("returns false", func() {
					rtf := NewRadioisotopeTestingFacility(RTFConfig{
						[]RTFArtifact{RTFGenerator{"a"}},
						[]RTFArtifact{},
						[]RTFArtifact{},
						[]RTFArtifact{RTFGenerator{"b"}},
					})
					Expect(rtf.done()).To(BeFalse())
				})
			})

			Context("everything IS on the fourth floor", func() {
				It("returns true", func() {
					rtf := NewRadioisotopeTestingFacility(RTFConfig{
						[]RTFArtifact{},
						[]RTFArtifact{},
						[]RTFArtifact{},
						[]RTFArtifact{RTFGenerator{"a"}, RTFGenerator{"b"}},
					})
					Expect(rtf.done()).To(BeTrue())
				})
			})
		})

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

	Describe("RTFHistory / RTFPermutations", func() {
		Describe("#Contains", func() {
			history := RTFHistory{
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
			}

			Context("history contains the thing", func() {
				It("returns true", func() {
					Expect(history.Contains(
						RadioisotopeTestingFacility{
							NewRTFConfig(
								RTFArtifacts{},
								RTFArtifacts{RTFGenerator{"a"}},
								RTFArtifacts{RTFMicrochip{"a"}, RTFMicrochip{"c"}},
								RTFArtifacts{},
							), 2})).To(BeTrue())
					Expect(history.Contains(
						RadioisotopeTestingFacility{
							NewRTFConfig(
								RTFArtifacts{RTFMicrochip{"a"}, RTFMicrochip{"c"}},
								RTFArtifacts{RTFGenerator{"a"}},
								RTFArtifacts{},
								RTFArtifacts{},
							), 0})).To(BeTrue())
				})
			})

			Context("history does not contain the thing", func() {
				Context("because it has a different ePos", func() {
					state := RadioisotopeTestingFacility{
						NewRTFConfig(
							RTFArtifacts{},
							RTFArtifacts{RTFGenerator{"a"}},
							RTFArtifacts{RTFMicrochip{"a"}, RTFMicrochip{"c"}},
							RTFArtifacts{},
						), 1}

					It("returns false", func() {
						Expect(history.Contains(state)).To(BeFalse())
					})
				})

				Context("because it has a different config", func() {
					state := RadioisotopeTestingFacility{
						NewRTFConfig(
							RTFArtifacts{},
							RTFArtifacts{RTFGenerator{"c"}},
							RTFArtifacts{RTFMicrochip{"a"}, RTFMicrochip{"c"}},
							RTFArtifacts{},
						), 2}

					It("returns false", func() {
						Expect(history.Contains(state)).To(BeFalse())
					})
				})
			})
		})
	})

	Describe("the test", func() {
		It("finds a solution", func() {
			config := RTFConfigRead(testSetup)
			solution := RTFTripPlan(config)
			Expect(len(solution)-1).To(Equal(11))
		})
	})

	Describe("the puzzle", func() {
		data := `The first floor contains a polonium generator, a thulium generator, a thulium-compatible microchip, a promethium generator, a ruthenium generator, a ruthenium-compatible microchip, a cobalt generator, and a cobalt-compatible microchip.
	The second floor contains a polonium-compatible microchip and a promethium-compatible microchip.
	The third floor contains nothing relevant.
	The fourth floor contains nothing relevant.`
		setup := strings.Split(string(data), "\n")
		
		It("finds a solution", func() {
			config := RTFConfigRead(setup)
			solution := RTFTripPlan(config)
			fmt.Println("MIKE: solution is", solution)
			fmt.Println("MIKE: took", len(solution)-1, "steps")
		})
	})
})
