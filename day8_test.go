package adventofcode2016_test

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
)

type TinyDisplayCommand interface {
	Rect(int, int)
	RotateCol(int, int)
	RotateRow(int, int)
}

type TinyDisplay struct {
	xSize, ySize int
	pixels       [][]bool
}

func NewTinyDisplay(xSize, ySize int) TinyDisplay {
	td := TinyDisplay{xSize, ySize, make([][]bool, ySize)}
	for jrow, _ := range td.pixels {
		td.pixels[jrow] = make([]bool, xSize)
	}
	return td
}

func (td *TinyDisplay) Rect(xSize, ySize int) {
	for y := 0; y < ySize; y++ {
		for x := 0; x < xSize; x++ {
			td.pixels[y][x] = true
		}
	}
}

func (td *TinyDisplay) RotateCol(colIndex, len int) {
	for times := 0; times < len; times++ {
		bottomPixel := td.pixels[td.ySize-1][colIndex]
		for y := td.ySize - 1; y > 0; y-- {
			td.pixels[y][colIndex] = td.pixels[y-1][colIndex]
		}
		td.pixels[0][colIndex] = bottomPixel
	}
}

func (td *TinyDisplay) RotateRow(rowIndex, len int) {
	row := td.pixels[rowIndex]
	for times := 0; times < len; times++ {
		rightPixel := row[td.xSize-1]
		for x := td.xSize - 1; x > 0; x-- {
			row[x] = row[x-1]
		}
		row[0] = rightPixel
	}
}

func (td TinyDisplay) String() string {
	output := "\n"
	for _, row := range td.pixels {
		for _, pixel := range row {
			if pixel {
				output += "#"
			} else {
				output += "."
			}
		}
		output += "\n"
	}
	return output
}

//
//  command input
//
var tdRectCommandRe = regexp.MustCompile(`rect (\d+)x(\d+)`)
var tdRotateColCommandRe = regexp.MustCompile(`rotate column x=(\d+) by (\d+)`)
var tdRotateRowCommandRe = regexp.MustCompile(`rotate row y=(\d+) by (\d+)`)

func TDCommandDispatch(command string, subject TinyDisplayCommand) {
	if match := tdRectCommandRe.FindStringSubmatch(command); match != nil {
		arg1, _ := strconv.Atoi(match[1])
		arg2, _ := strconv.Atoi(match[2])
		subject.Rect(arg1, arg2)
	} else if match := tdRotateColCommandRe.FindStringSubmatch(command); match != nil {
		arg1, _ := strconv.Atoi(match[1])
		arg2, _ := strconv.Atoi(match[2])
		subject.RotateCol(arg1, arg2)
	} else if match := tdRotateRowCommandRe.FindStringSubmatch(command); match != nil {
		arg1, _ := strconv.Atoi(match[1])
		arg2, _ := strconv.Atoi(match[2])
		subject.RotateRow(arg1, arg2)
	}
}

//
//  mock tiny display
//
type MockTD struct {
	method string
	arg1   int
	arg2   int
}

func (mtd *MockTD) Rect(arg1 int, arg2 int) {
	mtd.method = "Rect"
	mtd.arg1 = arg1
	mtd.arg2 = arg2
}

func (mtd *MockTD) RotateCol(arg1 int, arg2 int) {
	mtd.method = "RotateCol"
	mtd.arg1 = arg1
	mtd.arg2 = arg2
}

func (mtd *MockTD) RotateRow(arg1 int, arg2 int) {
	mtd.method = "RotateRow"
	mtd.arg1 = arg1
	mtd.arg2 = arg2
}

var _ = Describe("Day8", func() {
	Describe("TinyDisplay", func() {
		Describe("#String", func() {
			It("renders the pixels", func() {
				display := NewTinyDisplay(7, 3)
				Expect(display.String()).To(Equal("\n.......\n.......\n.......\n"))

				display.pixels[1][1] = true
				Expect(display.String()).To(Equal("\n.......\n.#.....\n.......\n"))
			})
		})

		Describe("#Rect", func() {
			It("draws a filled rectangle near the origin", func() {
				display := NewTinyDisplay(7, 3)
				display.Rect(3, 2)
				Expect(display.String()).To(Equal("\n###....\n###....\n.......\n"))
			})
		})

		Describe("#RotateCol", func() {
			It("rotates a column down", func() {
				display := NewTinyDisplay(7, 3)
				display.Rect(3, 2)
				display.RotateCol(1, 2)
				Expect(display.String()).To(Equal("\n###....\n#.#....\n.#.....\n"))
			})
		})

		Describe("#RotateRow", func() {
			It("rotates a row to the right", func() {
				display := NewTinyDisplay(7, 3)
				display.Rect(3, 2)
				display.RotateRow(0, 5)
				Expect(display.String()).To(Equal("\n#....##\n###....\n.......\n"))
			})
		})
	})

	Describe("#TDCommandDispatch", func() {
		Describe("rect", func() {
			It("calls Rect with appropriate args on the subject", func() {
				mtd := MockTD{}
				TDCommandDispatch("rect 3x2", &mtd)
				Expect(mtd.method).To(Equal("Rect"))
				Expect(mtd.arg1).To(Equal(3))
				Expect(mtd.arg2).To(Equal(2))

				TDCommandDispatch("rect 8x9", &mtd)
				Expect(mtd.method).To(Equal("Rect"))
				Expect(mtd.arg1).To(Equal(8))
				Expect(mtd.arg2).To(Equal(9))
			})
		})

		Describe("rotate row", func() {
			It("calls Rect with appropriate args on the subject", func() {
				mtd := MockTD{}
				TDCommandDispatch("rotate row y=1 by 5", &mtd)
				Expect(mtd.method).To(Equal("RotateRow"))
				Expect(mtd.arg1).To(Equal(1))
				Expect(mtd.arg2).To(Equal(5))

				TDCommandDispatch("rotate row y=2 by 12", &mtd)
				Expect(mtd.method).To(Equal("RotateRow"))
				Expect(mtd.arg1).To(Equal(2))
				Expect(mtd.arg2).To(Equal(12))
			})
		})

		Describe("rotate col", func() {
			It("calls Rect with appropriate args on the subject", func() {
				mtd := MockTD{}
				TDCommandDispatch("rotate column x=1 by 5", &mtd)
				Expect(mtd.method).To(Equal("RotateCol"))
				Expect(mtd.arg1).To(Equal(1))
				Expect(mtd.arg2).To(Equal(5))

				TDCommandDispatch("rotate column x=2 by 12", &mtd)
				Expect(mtd.method).To(Equal("RotateCol"))
				Expect(mtd.arg1).To(Equal(2))
				Expect(mtd.arg2).To(Equal(12))
			})
		})
	})

	Describe("the puzzle", func() {
		var parseFile = func(filename string) []string {
			data, _ := ioutil.ReadFile(filename)
			return strings.Split(string(data), "\n")
		}

		It("star 1 and 2", func() {
			commands := parseFile("day8.txt")
			td := NewTinyDisplay(50, 6)
			for _, command := range commands {
				TDCommandDispatch(command, &td)
			}

			litPixels := 0
			for _, row := range td.pixels {
				for _, pixel := range row {
					if pixel {
						litPixels++
					}
				}
			}
			fmt.Println("star 1: there are", litPixels, "lit pixels")
			fmt.Println(td)
		})
	})
})
