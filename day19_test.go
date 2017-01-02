package adventofcode2016_test

import (
	"fmt"
	"github.com/Workiva/go-datastructures/bitarray"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = fmt.Printf

type WhiteElephantParty struct {
	n_elves uint64
}

func NewWhiteElephantParty(n_elves uint64) *WhiteElephantParty {
	return &WhiteElephantParty{n_elves}
}

func (wep *WhiteElephantParty) winner() uint64 {
	n_elves := wep.n_elves
	elves := bitarray.NewBitArray(n_elves)

one_left:
	for {
		for jelf := uint64(0); jelf < wep.n_elves; jelf++ {
			if bit, _ := elves.GetBit(jelf); bit {
				continue
			}

			for jnext := uint64(1); jnext < wep.n_elves; jnext++ {
				jnext_elf := (jelf + jnext) % wep.n_elves
				if bit, _ := elves.GetBit(jnext_elf); bit {
					continue
				}
//				fmt.Printf("MIKE: %d takes from %d, leaving %d\n", jelf+1, jnext_elf+1, n_elves-1)
				elves.SetBit(jnext_elf)
				n_elves--
				break
			}

			if n_elves == 1 {
				break one_left
			}
		}
	}

	for jelf := uint64(0); jelf < wep.n_elves; jelf++ {
		if bit, _ := elves.GetBit(jelf); !bit {
			return jelf + 1 // map back into ordinal space
		}
	}
	panic("no elf found")
}

func compressIntSlice(slice []int) []int {
	// find first hole
	hole := -1
	for j := 0; j < len(slice); j++ {
		if slice[j] == -1 {
			hole = j
			break
		}
	}
	if hole == -1 {
		panic("could not compress, no hole")
	}

	k := hole
	for j := hole+1; j < len(slice); j++ {
		if slice[j] != -1 {
			slice[k] = slice[j]
			k++
		}
	}
	slice = slice[:k]
	return slice
}

func (wep *WhiteElephantParty) winner2() int {
	elves := make([]int, wep.n_elves)

	// populate the array with elf numbers
	for j := 0; j < int(wep.n_elves); j++ {
		elves[j] = j+1
	}

	previous_elf := 0
	n_elves := int(wep.n_elves)
	jelf := 0
	buffer := 0
	for n_elves > 1 {
		if elves[jelf] == -1 {
			elves = compressIntSlice(elves)
			buffer = 0
			
			// reset jelf pointer to the right place
			for j := 0; j < len(elves); j++ {
				if elves[j] == previous_elf {
					jelf = j + 1
				}
			}
			if jelf >= len(elves) {
				jelf = 0
			}
		}

		var jnext_elf int
		jnext_elf = (jelf + buffer + (n_elves / 2)) % len(elves)

		// fmt.Printf("MIKE: %d (idx %d) takes from %d (idx %d) (there are %d left)\n",
		// 	elves[jelf], jelf, elves[jnext_elf], jnext_elf, n_elves-1)
		elves[jnext_elf] = -1 // marker

		n_elves--
		buffer++

		previous_elf = elves[jelf]

		jelf++
		if jelf >= len(elves) {
			jelf = 0
		}
	}
	return previous_elf
}

var _ = Describe("Day19", func() {
	Describe("WhiteElephantParty", func() {
		It("picks the winner (by method 1)", func() {
			wep := NewWhiteElephantParty(5)
			Expect(wep.winner()).To(Equal(uint64(3)))
		})

		It("picks the winner (by method 2)", func() {
			wep := NewWhiteElephantParty(5)
			Expect(wep.winner2()).To(Equal(2))
		})
	})

	Describe("the puzzle", func() {
		It("experiment", func() {
			wep := NewWhiteElephantParty(10)
			winner := wep.winner2()
			fmt.Println("MIKE: experiment winner", winner)
		})

		It("star 1", func() {
			wep := NewWhiteElephantParty(3017957)
			winner := wep.winner()
			fmt.Println("star 1: winning elf is", winner)
		})

		It("star 2", func() {
			wep := NewWhiteElephantParty(3017957)
			winner := wep.winner2()
			fmt.Println("star 2: winning elf is", winner)
		})
	})
})
