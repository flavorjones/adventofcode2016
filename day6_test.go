package adventofcode2016_test

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"sort"
	"strings"
)

type RepetitionDecoder struct {
	messages []string
}

func (d RepetitionDecoder) calculateFrequencyDistribution() []SortableNameComponents {
	messageLen := len(d.messages[0])
	counts := make([]map[byte]uint, messageLen)

	for j := 0; j < messageLen; j++ {
		counts[j] = make(map[byte]uint)
	}

	for _, message := range d.messages {
		for j, char := range []byte(message) {
			count := counts[j]
			if n, ok := count[char]; ok {
				count[char] = n + 1
			} else {
				count[char] = 1
			}
		}
	}

	distribution := make([]SortableNameComponents, messageLen)
	for j := 0; j < messageLen; j++ {
		// convert to an array of SortableNameComponent (from day 4)
		count := counts[j]
		components := make(SortableNameComponents, 0, len(count)/2)
		for element, occurrences := range count {
			components = append(components, SortableNameComponent{element, occurrences})
		}
		sort.Sort(components)
		distribution[j] = components
	}
	return distribution
}

func (d RepetitionDecoder) decode() string {
	distribution := d.calculateFrequencyDistribution()
	decodedMessage := ""

	for _, components := range distribution {
		decodedMessage += string(components[0].element)
	}

	return string(decodedMessage)
}

func (d RepetitionDecoder) decode2() string {
	distribution := d.calculateFrequencyDistribution()
	decodedMessage := ""

	for _, components := range distribution {
		decodedMessage += string(components[len(components)-1].element)
	}

	return string(decodedMessage)
}

var _ = Describe("Day6", func() {
	var parseFile = func(filename string) []string {
		data, _ := ioutil.ReadFile(filename)
		return strings.Split(string(data), "\n")
	}

	Describe("RepetitionDecoder", func() {
		messages := parseFile("day6_test.txt")

		Describe("#decode", func() {
			It("decodes properly", func() {
				Expect(RepetitionDecoder{messages}.decode()).To(Equal("easter"))
			})
		})

		Describe("#decode2", func() {
			It("decodes properly", func() {
				Expect(RepetitionDecoder{messages}.decode2()).To(Equal("advent"))
			})
		})
	})

	Describe("RepetitionDecoder", func() {
		messages := parseFile("day6_data.txt")

		Describe("star 1", func() {
			It("finds the answer", func() {
				fmt.Println("star 1:", RepetitionDecoder{messages}.decode())
			})
		})

		Describe("star 2", func() {
			It("finds the answer", func() {
				fmt.Println("star 2:", RepetitionDecoder{messages}.decode2())
			})
		})
	})
})
