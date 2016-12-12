package adventofcode2016_test

import (
	"bytes"
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io"
	"io/ioutil"
	"regexp"
	"strconv"
)

var _ = fmt.Println

type ExpFormat struct {
	content string
}

var expFormatMarkerRe = regexp.MustCompile(`\((\d+)x(\d+)\)`)

func getMarkerData(content []byte) (nchars int, times int) {
	match := expFormatMarkerRe.FindSubmatch(content)
	nchars, _ = strconv.Atoi(string(match[1]))
	times, _ = strconv.Atoi(string(match[2]))
	return
}

func (ef ExpFormat) Decompress() string {
	content := bytes.NewBufferString(ef.content)
	decompressed := bytes.Buffer{}

	for {
		byte, err := content.ReadByte()
		if err == io.EOF {
			break
		}

		if byte == '(' {
			content.UnreadByte()

			marker, _ := content.ReadBytes(')')
			nchars, times := getMarkerData(marker)

			repeatingSegment := content.Next(nchars)
			for j := 0; j < times; j++ {
				decompressed.Write(repeatingSegment)
			}
		} else {
			decompressed.WriteByte(byte)
		}
	}

	rval := decompressed.String()
	return rval
}

func (ef ExpFormat) Decompress2Len() int {
	content := bytes.NewBufferString(ef.content)
	byteCount := 0

	for {
		byte, err := content.ReadByte()
		if err == io.EOF {
			break
		}

		if byte == '(' {
			content.UnreadByte()

			marker, _ := content.ReadBytes(')')
			nchars, times := getMarkerData(marker)

			nchars = ExpFormat{string(content.Next(nchars))}.Decompress2Len()
			byteCount += nchars * times
		} else {
			byteCount++
		}
	}

	return byteCount
}

var _ = Describe("Day9", func() {
	Describe("ExpFormat", func() {
		Describe("#Decompress", func() {
			It("decompresses markerless text", func() {
				Expect(ExpFormat{"ADVENT"}.Decompress()).To(Equal("ADVENT"))
			})

			It("decompresses simple markers", func() {
				Expect(ExpFormat{"A(1x5)BC"}.Decompress()).To(Equal("ABBBBBC"))
				Expect(ExpFormat{"(3x3)XYZ"}.Decompress()).To(Equal("XYZXYZXYZ"))
			})

			It("decompresses multiple markers", func() {
				Expect(ExpFormat{"A(2x2)BCD(2x2)EFG"}.Decompress()).To(Equal("ABCBCDEFEFG"))
			})

			It("ignores markers that are part of repeating segments", func() {
				Expect(ExpFormat{"(6x1)(1x3)A"}.Decompress()).To(Equal("(1x3)A"))
				Expect(ExpFormat{"X(8x2)(3x3)ABCY"}.Decompress()).To(Equal("X(3x3)ABC(3x3)ABCY"))
			})
		})

		Describe("#Decompress2Len", func() {
			It("returns the length of the document decompressed by alternative algo", func() {
				Expect(ExpFormat{"(3x3)XYZ"}.Decompress2Len()).To(Equal(9))
				Expect(ExpFormat{"X(8x2)(3x3)ABCY"}.Decompress2Len()).To(Equal(len("XABCABCABCABCABCABCY")))
				Expect(ExpFormat{"(27x12)(20x12)(13x14)(7x10)(1x12)A"}.Decompress2Len()).To(Equal(241920))
				Expect(ExpFormat{"(25x3)(3x3)ABC(2x3)XY(5x2)PQRSTX(18x9)(3x2)TWO(5x7)SEVEN"}.Decompress2Len()).To(Equal(445))
			})
		})
	})

	Describe("the puzzle", func() {
		data, _ := ioutil.ReadFile("day9.txt")

		Describe("star 1", func() {
			It("prints the decompressed size of the puzzle data", func() {
				ef := ExpFormat{string(data)}
				fmt.Println("star 1: decompressed size is", len(ef.Decompress()))
			})
		})

		Describe("star 2", func() {
			It("prints the alt-decompressed size of the puzzle data", func() {
				ef := ExpFormat{string(data)}
				fmt.Println("star 2: decompressed size is", ef.Decompress2Len())
			})
		})
	})
})
