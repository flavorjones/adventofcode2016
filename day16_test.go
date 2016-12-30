package adventofcode2016_test

// http://adventofcode.com/2016/day/16

import (
	"fmt"
	"github.com/Workiva/go-datastructures/bitarray"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = fmt.Println

type DragonData struct {
	length uint64
	bits   bitarray.BitArray
}

func NewDragonData(input string) *DragonData {
	length := uint64(len(input))
	bits := bitarray.NewBitArray(length)
	for jchar, char := range []byte(input) {
		if char == '1' {
			bits.SetBit(uint64(jchar))
		}
	}
	return &DragonData{length, bits}
}

func (dd *DragonData) Cycle() {
	newLength := dd.length*2 + 1
	newBits := bitarray.NewBitArray(newLength)

	for j := uint64(0); j < dd.length; j++ {
		if flipped, _ := dd.bits.GetBit(j); !flipped {
			k := newLength - j - 1
			if error := newBits.SetBit(k); error != nil {
				panic(fmt.Sprintf("could not set bit %d on bitarray length %d cap %d\n", k, newLength, newBits.Capacity()))
			}
		}
	}
	newBits = newBits.Or(dd.bits)

	dd.length = newLength
	dd.bits = newBits
}

func (dd *DragonData) CycleToFill(diskSize uint64) {
	for dd.length < diskSize {
		dd.Cycle()
	}
	dd.length = diskSize
}

func (dd *DragonData) Checksum() string {
	length := dd.length
	bits := bitarray.NewBitArray(length).Or(dd.bits) // make a copy

	for (length % 2) == 0 {
		nextLength := length / 2
		nextBits := bitarray.NewBitArray(nextLength)

		for j := uint64(0); j < length; j += 2 {
			flipped1, _ := bits.GetBit(j)
			flipped2, _ := bits.GetBit(j + 1)
			if flipped1 == flipped2 {
				nextBits.SetBit(j / 2)
			}
		}
		length = nextLength
		bits = nextBits
	}

	return sprintbits(bits, length)
}

func (dd *DragonData) dataString() string {
	return sprintbits(dd.bits, dd.length)
}

// ----------------------------------------
// utilities functions
func sprintbits(ba bitarray.BitArray, length uint64) string {
	output := make([]byte, length)

	for j := uint64(0); j < length; j++ {
		if flipped, _ := ba.GetBit(j); flipped {
			output[j] = '1'
		} else {
			output[j] = '0'
		}
	}

	return string(output)
}

// ----------------------------------------
// tests
var _ = Describe("Day16", func() {
	Describe("DragonData", func() {
		Describe(".NewDragonData", func() {
			It("sets data properly", func() {
				dd := NewDragonData("11111111100001010")
				Expect(dd.dataString()).To(Equal("11111111100001010"))
			})
		})

		Describe("#Cycle", func() {
			It("generates one iteration of dragon-curve data", func() {
				dd := NewDragonData("1")
				dd.Cycle()
				Expect(dd.dataString()).To(Equal("100"))
			})

			It("generates one iteration of dragon-curve data", func() {
				dd := NewDragonData("0")
				dd.Cycle()
				Expect(dd.dataString()).To(Equal("001"))
			})

			It("generates one iteration of dragon-curve data", func() {
				dd := NewDragonData("11111")
				dd.Cycle()
				Expect(dd.dataString()).To(Equal("11111000000"))
			})

			It("generates one iteration of dragon-curve data", func() {
				dd := NewDragonData("111100001010")
				dd.Cycle()
				Expect(dd.dataString()).To(Equal("1111000010100101011110000"))
			})
		})

		Describe("#CycleToFill", func() {
			It("cycles until it fills the disk", func() {
				dd := NewDragonData("111100001010")
				dd.CycleToFill(23)
				Expect(dd.dataString()).To(Equal("11110000101001010111100"))
			})
		})

		Describe("#Checksum", func() {
			It("calculates the checksum as specified", func() {
				dd := NewDragonData("110010110100")
				Expect(dd.Checksum()).To(Equal("100"))
			})
		})

		Describe("smoketest", func() {
			It("combines these functions correctly", func() {
				dd := NewDragonData("10000")
				dd.CycleToFill(20)
				Expect(dd.Checksum()).To(Equal("01100"))
			})
		})
	})

	Describe("the puzzle", func() {
		It("star 1", func() {
			dd := NewDragonData("01000100010010111")
			dd.CycleToFill(272)
			fmt.Println("day 16 star 1: checksum is", dd.Checksum())
		})

		It("star 2", func() {
			dd := NewDragonData("01000100010010111")
			dd.CycleToFill(35651584)
			fmt.Println("day 16 star 2: checksum is", dd.Checksum())
		})
	})

})
