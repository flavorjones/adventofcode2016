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

	jbyte := 0
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

		jbyte++
	}

	rval := decompressed.String()
	return rval
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
	})

	Describe("the puzzle", func() {
		Describe("star 1", func() {
			It("prints the decompressed size of the puzzle data", func() {
				data, _ := ioutil.ReadFile("day9.txt")
				ef := ExpFormat{string(data)}
				fmt.Println("star 1: decompressed size is", len(ef.Decompress()))
			})
		})
	})
})
