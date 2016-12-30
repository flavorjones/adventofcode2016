package adventofcode2016_test

import (
	"crypto/md5"
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"strconv"
)

var _ = fmt.Sprintf

type KeyGenerator struct {
	salt         string
	stretch      int
	foundKeys    map[int]int // ordinal → index
//	cachedHashes map[int][]byte
}

func NewKeyGenerator(salt string) *KeyGenerator {
	return &KeyGenerator{salt, 0, make(map[int]int)}
}

func NewStretchedKeyGenerator(salt string) *KeyGenerator {
	return &KeyGenerator{salt, 2016, make(map[int]int)}
}

func (kg *KeyGenerator) Hash(index int) []byte {
	value := kg.salt + strconv.Itoa(index)
	hash := []byte(fmt.Sprintf("%x", md5.Sum([]byte(value))))
	for j := 0; j < kg.stretch; j++ {
		hash = []byte(fmt.Sprintf("%x", md5.Sum(hash)))
	}
	return hash
}

func (kg *KeyGenerator) calculateKey(position int) int {
	var start int
	if position == 1 {
		start = 0
	} else {
		start = kg.Key(position-1) + 1
	}

	for j := start; ; j++ {
		repeatEh, repeatChar := any3Repeat(kg.Hash(j))
		if repeatEh {
			for k := j + 1; k < j+1000; k++ {
				if specific5Repeat(kg.Hash(k), repeatChar) {
					return j
				}
			}
		}
	}
	return -1
}

func (kg *KeyGenerator) Key(position int) int {
	// this is a read-through cache
	value, ok := kg.foundKeys[position]
	if !ok {
		value = kg.calculateKey(position) // will recurse if necessary
		kg.foundKeys[position] = value
	}
	return value
}

// ----------------------------------------
// utility methods

func any3Repeat(hash []byte) (bool, byte) {
	fmt.Printf("MIKE: looking for 3 repeats in %s\n", hash)
	for j := 0; j < len(hash)-2; j++ {
		if hash[j] == hash[j+1] && hash[j] == hash[j+2] {
			return true, hash[j]
		}
	}
	return false, '-'
}

func specific5Repeat(hash []byte, char byte) bool {
	//	fmt.Printf("MIKE: looking for 5 repeats of %c in %s\n", char, hash)
	for j := 0; j < len(hash)-4; j++ {
		if hash[j] == char &&
			hash[j] == hash[j+1] &&
			hash[j] == hash[j+2] &&
			hash[j] == hash[j+3] &&
			hash[j] == hash[j+4] {
			return true
		}
	}
	return false
}

// ----------------------------------------
// tests

var _ = Describe("Day14", func() {
	Describe("KeyGenerator", func() {
		Context("single hash", func() {
			It("generates the first based on a salt", func() {
				Expect(NewKeyGenerator("abc").Key(1)).To(Equal(39))
			})

			It("generates the second key based on a salt, implicitly calculating previous", func() {
				Expect(NewKeyGenerator("abc").Key(2)).To(Equal(92))
			})

			It("scales", func() {
				Expect(NewKeyGenerator("abc").Key(64)).To(Equal(22728))
			})
		})

		Context("stretched", func() {
			It("generates the first based on a salt", func() {
				Expect(NewStretchedKeyGenerator("abc").Key(1)).To(Equal(10))
			})

			It("scales", func() {
				Expect(NewStretchedKeyGenerator("abc").Key(64)).To(Equal(22551))
			})
		})
	})

	Describe("the puzzle", func() {
		Describe("star 1", func() {
			It("finds the 64th key", func() {
				index := NewKeyGenerator("qzyelonm").Key(64)
				fmt.Println("star 1: 64th keys index is", index)
			})
		})
	})
})
