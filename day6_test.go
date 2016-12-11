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

func (d RepetitionDecoder) decode() string {
	messageLen := len(d.messages[0])
	counts := make([]map[byte]uint, messageLen)
	decodedMessage := make([]byte, messageLen)

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

	for j := 0; j < messageLen; j++ {
		// convert to an array of SortableNameComponent (from day 4)
		count := counts[j]
		components := make(SortableNameComponents, 0, len(count)/2)
		for element, occurrences := range count {
			components = append(components, SortableNameComponent{element, occurrences})
		}
		sort.Sort(components)
		decodedMessage[j] = components[0].element
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
	})

	Describe("RepetitionDecoder", func() {
		messages := parseFile("day6_data.txt")

		Describe("star 1", func() {
			It("finds the answer", func() {
				fmt.Println("star 1:", RepetitionDecoder{messages}.decode())
			})
		})
	})
})
