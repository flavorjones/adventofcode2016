package adventofcode2016_test

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type TilePredictor struct {
	floor [][]bool // true if trap
}

func NewTilePredictor(input string) *TilePredictor {
	var floor [][]bool
	row := make([]bool, len(input))

	for j, char := range []byte(input) {
		if char == '^' {
			row[j] = true
		} else {
			row[j] = false
		}
	}

	floor = append(floor, row)
	return &TilePredictor{floor}
}

// returns a pointer to self so I can chain.
func (tp *TilePredictor) next() *TilePredictor {
	current_row := tp.floor[len(tp.floor)-1]
	next_row := make([]bool, len(current_row))

	for j := 0; j < len(current_row); j++ {
		var left, center, right bool
		if j == 0 {
			left = false
		} else {
			left = current_row[j-1]
		}
		if j == len(current_row)-1 {
			right = false
		} else {
			right = current_row[j+1]
		}
		center = current_row[j]

		next_row[j] = ((left && center && !right) ||
			(center && right && !left) ||
			(left && !center && !right) ||
			(right && !center && !left))
	}
	tp.floor = append(tp.floor, next_row)

	return tp
}

func (tp *TilePredictor) currentString() string {
	row := tp.floor[len(tp.floor)-1]
	rval := make([]byte, len(row))
	for j, tile := range row {
		if tile {
			rval[j] = '^'
		} else {
			rval[j] = '.'
		}
	}
	return string(rval)
}

func (tp *TilePredictor) safeCount() int {
	count := 0
	for j := 0; j < len(tp.floor); j++ {
		for k := 0; k < len(tp.floor[j]); k++ {
			if !tp.floor[j][k] {
				count++
			}
		}
	}
	return count
}

var _ = Describe("Day18", func() {
	Describe("TilePredictor", func() {
		Describe("#currentString", func() {
			It("outputs a readable interpretation of the most recent row", func() {
				input := `..^^.`
				tp := NewTilePredictor(input)
				Expect(tp.currentString()).To(Equal(`..^^.`))
			})
		})

		It("predicts a small set of tiles", func() {
			input := `..^^.`
			tp := NewTilePredictor(input)
			Expect(tp.next().currentString()).To(Equal(`.^^^^`))
			Expect(tp.next().currentString()).To(Equal(`^^..^`))
		})

		It("predicts a larger set of tiles", func() {
			tp := NewTilePredictor(`.^^.^.^^^^`)
			Expect(tp.next().currentString()).To(Equal(`^^^...^..^`))
			Expect(tp.next().currentString()).To(Equal(`^.^^.^.^^.`))
			Expect(tp.next().currentString()).To(Equal(`..^^...^^^`))
			Expect(tp.next().currentString()).To(Equal(`.^^^^.^^.^`))
			Expect(tp.next().currentString()).To(Equal(`^^..^.^^..`))
			Expect(tp.next().currentString()).To(Equal(`^^^^..^^^.`))
			Expect(tp.next().currentString()).To(Equal(`^..^^^^.^^`))
			Expect(tp.next().currentString()).To(Equal(`.^^^..^.^^`))
			Expect(tp.next().currentString()).To(Equal(`^^.^^^..^^`))
			Expect(tp.safeCount()).To(Equal(38))
		})
	})

	Describe("the puzzle", func() {
		var tp *TilePredictor

		BeforeEach(func() {
			tp = NewTilePredictor(`^.^^^..^^...^.^..^^^^^.....^...^^^..^^^^.^^.^^^^^^^^.^^.^^^^...^^...^^^^.^.^..^^..^..^.^^.^.^.......`)
		})

		It("star 1", func() {
			for j := 1; j < 40; j++ {
				tp.next()
			}
			fmt.Println("star 1: there are", tp.safeCount(), "safe tiles")
		})

		It("star 2", func() {
			for j := 1; j < 400000; j++ {
				tp.next()
			}
			fmt.Println("star 2: there are", tp.safeCount(), "safe tiles")
		})
	})
})
