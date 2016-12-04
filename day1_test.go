package adventofcode2016_test

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"math"
	"regexp"
	"strconv"
	"strings"
)

type Coordinates struct {
	x, y int
}

var NORTH = Coordinates{0, 1}
var EAST = Coordinates{1, 0}
var SOUTH = Coordinates{0, -1}
var WEST = Coordinates{-1, 0}
var pathSegmentRe = regexp.MustCompile("([LR])([0-9]+)")

type Position struct {
	location Coordinates
	heading  Coordinates
}

func NewPosition() *Position {
	return &Position{heading: NORTH}
}

func (self *Position) move(segment string) {
	match := pathSegmentRe.FindStringSubmatch(segment)
	direction := match[1]
	distance, err := strconv.ParseUint(match[2], 10, 16)
	if err != nil {
		panic(fmt.Sprintf("cannot parse '%s' as a uint", match[2]))
	}

	self.turn(direction)
	self.walk(uint(distance))
}

func (self *Position) turn(direction string) {
	switch direction {
	case "L":
		switch self.heading {
		case NORTH:
			self.heading = WEST
		case WEST:
			self.heading = SOUTH
		case SOUTH:
			self.heading = EAST
		default:
			self.heading = NORTH
		}
	default:
		switch self.heading {
		case NORTH:
			self.heading = EAST
		case EAST:
			self.heading = SOUTH
		case SOUTH:
			self.heading = WEST
		default:
			self.heading = NORTH
		}
	}
}

func (self *Position) walk(distance uint) {
	self.location.x += self.heading.x * int(distance)
	self.location.y += self.heading.y * int(distance)
}

func (self Position) taxicabGeometry() uint {
	return uint(math.Abs(float64(self.location.x)) + math.Abs(float64(self.location.y)))
}

type GridPath struct {
	path string
}

func (self GridPath) segments() []string {
	return strings.Split(self.path, ", ")
}

func (self GridPath) distance() uint {
	position := NewPosition()
	for _, segment := range self.segments() {
		position.move(segment)
	}
	return position.taxicabGeometry()
}

var _ = Describe("Day1", func() {
	Describe("Position", func() {
		Describe("move", func() {
			It("updates heading correctly rightwise", func() {
				position := Position{heading: NORTH} //NewPosition()
				Expect(position.heading).To(Equal(NORTH))

				position.move("R1")
				Expect(position.heading).To(Equal(EAST))

				position.move("R1")
				Expect(position.heading).To(Equal(SOUTH))

				position.move("R1")
				Expect(position.heading).To(Equal(WEST))
			})

			It("updates heading correctly leftwise", func() {
				position := NewPosition()
				Expect(position.heading).To(Equal(NORTH))

				position.move("L0")
				Expect(position.heading).To(Equal(WEST))

				position.move("L0")
				Expect(position.heading).To(Equal(SOUTH))

				position.move("L0")
				Expect(position.heading).To(Equal(EAST))
			})

			It("updates location", func() {
				position := NewPosition()
				Expect(position.location).To(Equal(Coordinates{0, 0}))

				position.move("R2")
				Expect(position.location).To(Equal(Coordinates{2, 0}))

				position.move("R2")
				Expect(position.location).To(Equal(Coordinates{2, -2}))

				position.move("R2")
				Expect(position.location).To(Equal(Coordinates{0, -2}))

				position.move("R2")
				Expect(position.location).To(Equal(Coordinates{0, 0}))
			})
		})
	})

	Describe("GridPath", func() {
		Describe("#distance", func() {
			It("adds two segments", func() {
				Expect(GridPath{"R2, L3"}.distance()).To(Equal(uint(5)))
			})

			It("adds three segments", func() {
				Expect(GridPath{"R2, R2, R2"}.distance()).To(Equal(uint(2)))
			})

			It("adds four segments", func() {
				Expect(GridPath{"R5, L5, R5, R3"}.distance()).To(Equal(uint(12)))
			})
		})
	})

	Describe("the puzzle", func() {
		It("star 1", func() {
			path := "L1, L5, R1, R3, L4, L5, R5, R1, L2, L2, L3, R4, L2, R3, R1, L2, R5, R3, L4, R4, L3, R3, R3, L2, R1, L3, R2, L1, R4, L2, R4, L4, R5, L3, R1, R1, L1, L3, L2, R1, R3, R2, L1, R4, L4, R2, L189, L4, R5, R3, L1, R47, R4, R1, R3, L3, L3, L2, R70, L1, R4, R185, R5, L4, L5, R4, L1, L4, R5, L3, R2, R3, L5, L3, R5, L1, R5, L4, R1, R2, L2, L5, L2, R4, L3, R5, R1, L5, L4, L3, R4, L3, L4, L1, L5, L5, R5, L5, L2, L1, L2, L4, L1, L2, R3, R1, R1, L2, L5, R2, L3, L5, L4, L2, L1, L2, R3, L1, L4, R3, R3, L2, R5, L1, L3, L3, L3, L5, R5, R1, R2, L3, L2, R4, R1, R1, R3, R4, R3, L3, R3, L5, R2, L2, R4, R5, L4, L3, L1, L5, L1, R1, R2, L1, R3, R4, R5, R2, R3, L2, L1, L5"

			fmt.Println("star 1 distance is", GridPath{path}.distance())
		})
	})
})
