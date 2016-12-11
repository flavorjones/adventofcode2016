package adventofcode2016_test

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"regexp"
	"strings"
)

func stringHasAbbaNature(word string) bool {
	for j := 0; j < len(word)-3; j++ {
		if word[j] == word[j+3] &&
			word[j+1] == word[j+2] &&
			word[j] != word[j+1] {
			return true
		}
	}
	return false
}

type IPv7Part struct {
	word       string
	isHypernet bool
}

type IPv7 struct {
	address string
}

var ipv7PartsRe = regexp.MustCompile(`(\b\w+\b)+`)

func (ip IPv7) parts() []IPv7Part {
	matches := ipv7PartsRe.FindAllStringSubmatch(ip.address, -1)
	parts := make([]IPv7Part, 0, 3)

	isHypernet := false
	for _, match := range matches {
		parts = append(parts, IPv7Part{match[0], isHypernet})
		isHypernet = !isHypernet
	}

	return parts
}

func (ip IPv7) supportsTLS() bool {
	abbaSomewhere := false
	for _, part := range ip.parts() {
		if part.isHypernet {
			if stringHasAbbaNature(part.word) {
				return false
			}
		} else {
			if stringHasAbbaNature(part.word) {
				abbaSomewhere = true
			}
		}
	}
	return abbaSomewhere
}

var _ = Describe("Day7", func() {
	Describe("#stringHasAbbaNature", func() {
		It("looks for the abba pattern", func() {
			Expect(stringHasAbbaNature("abba")).To(BeTrue(), "abba")
			Expect(stringHasAbbaNature("abcd")).To(BeFalse(), "abcd")
			Expect(stringHasAbbaNature("aaaa")).To(BeFalse(), "aaaa")
			Expect(stringHasAbbaNature("ioxxoj")).To(BeTrue(), "ioxxoj")
			Expect(stringHasAbbaNature("ababababatuut")).To(BeTrue(), "ababababatuut")
			Expect(stringHasAbbaNature("tuutababababa")).To(BeTrue(), "tuutababababa")
		})
	})

	Describe("IPv7", func() {
		Describe("#parts", func() {
			Context("address has three parts", func() {
				It("returns the address parts", func() {
					parts := IPv7{"foo[bar]quux"}.parts()
					Expect(parts[0]).To(Equal(IPv7Part{"foo", false}))
					Expect(parts[1]).To(Equal(IPv7Part{"bar", true}))
					Expect(parts[2]).To(Equal(IPv7Part{"quux", false}))
				})
			})

			Context("address has five parts", func() {
				It("returns the address parts", func() {
					parts := IPv7{"foo[bar]bazz[quux]quuux"}.parts()
					Expect(parts[0]).To(Equal(IPv7Part{"foo", false}))
					Expect(parts[1]).To(Equal(IPv7Part{"bar", true}))
					Expect(parts[2]).To(Equal(IPv7Part{"bazz", false}))
					Expect(parts[3]).To(Equal(IPv7Part{"quux", true}))
					Expect(parts[4]).To(Equal(IPv7Part{"quuux", false}))
				})
			})
		})

		Describe("#supportsTLS", func() {
			It("parses the address to determine support", func() {
				Expect(IPv7{"abba[mnop]qrst"}.supportsTLS()).To(BeTrue())
				Expect(IPv7{"qrst[mnop]abba"}.supportsTLS()).To(BeTrue())
				Expect(IPv7{"abcd[bddb]xyyx"}.supportsTLS()).To(BeFalse())
				Expect(IPv7{"xyyx[bddb]abcd"}.supportsTLS()).To(BeFalse())
				Expect(IPv7{"aaaa[qwer]tyui"}.supportsTLS()).To(BeFalse())
				Expect(IPv7{"bccbaaaa[qwer]tyui"}.supportsTLS()).To(BeTrue())
				Expect(IPv7{"aaaabccb[qwer]tyui"}.supportsTLS()).To(BeTrue())
				Expect(IPv7{"ioxxoj[asdfgh]zxcvbn"}.supportsTLS()).To(BeTrue())
				Expect(IPv7{"a[b]c[d]effe"}.supportsTLS()).To(BeTrue())
				Expect(IPv7{"a[b]c[deed]effe"}.supportsTLS()).To(BeFalse())
			})
		})
	})

	Describe("the puzzle", func() {
		var parseFile = func(filename string) []string {
			data, _ := ioutil.ReadFile(filename)
			return strings.Split(string(data), "\n")
		}
		addresses := parseFile("day7.txt")

		Describe("star 1", func() {
			Specify("count the addresses that support TLS", func() {
				nMatches := 0
				for _, address := range addresses {
					if (IPv7{address}).supportsTLS() {
						nMatches++
					}
				}
				fmt.Println("star 1:", nMatches, "matches")
			})
		})
	})
})
