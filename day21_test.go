package adventofcode2016_test

import (
	"bytes"
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
)

type PasswordScrambler struct {
	pw []byte
}

var pwsSwap1Re = regexp.MustCompile(`swap position (\d+) with position (\d+)`)
var pwsSwap2Re = regexp.MustCompile(`swap letter (\w) with letter (\w)`)
var pwsRevRe = regexp.MustCompile(`reverse positions (\d+) through (\d+)`)
var pwsRot1Re = regexp.MustCompile(`rotate (\w+) (\d+) steps?`)
var pwsRot2Re = regexp.MustCompile(`rotate based on position of letter (\w)`)
var pwsMovRe = regexp.MustCompile(`move position (\d+) to position (\d+)`)

func NewPasswordScrambler(password string) *PasswordScrambler {
	return &PasswordScrambler{[]byte(password)}
}

func (ps *PasswordScrambler) password() string {
	return string(ps.pw)
}

var rotLeft = func(password []byte) {
	first := password[0]
	for j := 0; j < len(password)-1; j++ {
		password[j] = password[j+1]
	}
	password[len(password)-1] = first
}

var rotRight = func(password []byte) {
	last := password[len(password)-1]
	for j := len(password) - 1; j > 0; j-- {
		password[j] = password[j-1]
	}
	password[0] = last
}

var reverse = func(password []byte, pos1, pos2 int) {
	for j := 0; j <= (pos2-pos1)/2; j++ {
		j1, j2 := pos1+j, pos2-j
		password[j1], password[j2] = password[j2], password[j1]
	}
}

func (ps *PasswordScrambler) do(command string) {
	switch {
	case pwsSwap1Re.MatchString(command):
		matches := pwsSwap1Re.FindStringSubmatch(command)
		pos1s, pos2s := matches[1], matches[2]
		pos1, _ := strconv.Atoi(pos1s)
		pos2, _ := strconv.Atoi(pos2s)
		ps.pw[pos1], ps.pw[pos2] = ps.pw[pos2], ps.pw[pos1]

	case pwsSwap2Re.MatchString(command):
		matches := pwsSwap2Re.FindStringSubmatch(command)
		char1, char2 := matches[1][0], matches[2][0]
		pos1 := bytes.IndexByte(ps.pw, char1)
		pos2 := bytes.IndexByte(ps.pw, char2)
		ps.pw[pos1], ps.pw[pos2] = ps.pw[pos2], ps.pw[pos1]

	case pwsRevRe.MatchString(command):
		matches := pwsRevRe.FindStringSubmatch(command)
		pos1s, pos2s := matches[1], matches[2]
		pos1, _ := strconv.Atoi(pos1s)
		pos2, _ := strconv.Atoi(pos2s)
		reverse(ps.pw, pos1, pos2)

	case pwsRot1Re.MatchString(command):
		matches := pwsRot1Re.FindStringSubmatch(command)
		direction, stepsS := matches[1], matches[2]
		steps, _ := strconv.Atoi(stepsS)
		for jstep := 0; jstep < steps; jstep++ {
			if direction == "left" {
				rotLeft(ps.pw)
			} else {
				rotRight(ps.pw)
			}
		}

	case pwsRot2Re.MatchString(command):
		matches := pwsRot2Re.FindStringSubmatch(command)
		char := matches[1][0]
		index := bytes.IndexByte(ps.pw, char)
		rotRight(ps.pw)
		for j := 0; j < index; j++ {
			rotRight(ps.pw)
		}
		if index >= 4 {
			rotRight(ps.pw)
		}

	case pwsMovRe.MatchString(command):
		matches := pwsMovRe.FindStringSubmatch(command)
		pos1s, pos2s := matches[1], matches[2]
		pos1, _ := strconv.Atoi(pos1s)
		pos2, _ := strconv.Atoi(pos2s)
		char := ps.pw[pos1]
		ps.pw = append(ps.pw[:pos1], ps.pw[pos1+1:]...)                        // remove char
		ps.pw = append(ps.pw[:pos2], append([]byte{char}, ps.pw[pos2:]...)...) // insert char at pos2

	default:
		panic(fmt.Sprintf("unknown command '%s'", command))
	}
}

func (ps *PasswordScrambler) undo(command string) {
	switch {
	case pwsSwap1Re.MatchString(command):
		ps.do(command)

	case pwsSwap2Re.MatchString(command):
		ps.do(command)

	case pwsRevRe.MatchString(command):
		ps.do(command)

	case pwsRot1Re.MatchString(command):
		matches := pwsRot1Re.FindStringSubmatch(command)
		direction, stepsS := matches[1], matches[2]
		steps, _ := strconv.Atoi(stepsS)
		for jstep := 0; jstep < steps; jstep++ {
			if direction == "left" {
				rotRight(ps.pw)
			} else {
				rotLeft(ps.pw)
			}
		}

	case pwsRot2Re.MatchString(command):
		// brute force because I don't care
		var save []byte
		save = append(save, ps.pw...)
		fmt.Println("MIKE: target is", string(save))
		for j := 0; j < len(ps.pw); j++ {
			copy(ps.pw, save)
			for k := 0; k < j; k++ {
				rotLeft(ps.pw)
			}
			fmt.Println("MIKE: trying", string(ps.pw))
			ps.do(command)
			fmt.Println("MIKE: →", string(ps.pw), "(compared to", string(save), ")")
			if bytes.Equal(ps.pw, save) {
				for k := 0; k < j; k++ {
					rotLeft(ps.pw)
				}
				break
			}
		}

	case pwsMovRe.MatchString(command):
		matches := pwsMovRe.FindStringSubmatch(command)
		pos1s, pos2s := matches[1], matches[2]
		pos1, _ := strconv.Atoi(pos1s)
		pos2, _ := strconv.Atoi(pos2s)
		char := ps.pw[pos2]
		ps.pw = append(ps.pw[:pos2], ps.pw[pos2+1:]...)                        // remove char
		ps.pw = append(ps.pw[:pos1], append([]byte{char}, ps.pw[pos1:]...)...) // insert char at pos1

	default:
		panic(fmt.Sprintf("unknown command '%s'", command))
	}
}

var _ = Describe("Day21", func() {
	Describe("PasswordScrambler", func() {
		It("does a bunch of shit", func() {
			ps := NewPasswordScrambler(`abcde`)

			ps.do(`swap position 4 with position 0`)
			Expect(ps.password()).To(Equal(`ebcda`))

			ps.do(`swap letter d with letter b`)
			Expect(ps.password()).To(Equal(`edcba`))

			ps.do(`reverse positions 0 through 4`)
			Expect(ps.password()).To(Equal(`abcde`))

			ps.do(`rotate left 1 step`)
			Expect(ps.password()).To(Equal(`bcdea`))
			ps.do(`rotate right 1 step`)
			Expect(ps.password()).To(Equal(`abcde`))
			ps.do(`rotate left 1 step`)
			Expect(ps.password()).To(Equal(`bcdea`))

			ps.do(`move position 1 to position 4`)
			Expect(ps.password()).To(Equal(`bdeac`))

			ps.do(`move position 3 to position 0`)
			Expect(ps.password()).To(Equal(`abdec`))

			ps.do(`rotate based on position of letter b`)
			Expect(ps.password()).To(Equal(`ecabd`))

			ps.do(`rotate based on position of letter d`)
			Expect(ps.password()).To(Equal(`decab`))
		})

		It("undoes a bunch of shit", func() {
			ps := NewPasswordScrambler(`decab`)

			ps.undo(`rotate based on position of letter d`)
			Expect(ps.password()).To(Equal(`ecabd`))

			ps.undo(`rotate based on position of letter b`)
			Expect(ps.password()).To(Equal(`abdec`))

			ps.undo(`move position 3 to position 0`)
			Expect(ps.password()).To(Equal(`bdeac`))

			ps.undo(`move position 1 to position 4`)
			Expect(ps.password()).To(Equal(`bcdea`))

			ps.undo(`rotate left 1 step`)
			Expect(ps.password()).To(Equal(`abcde`))
			ps.undo(`rotate right 1 step`)
			Expect(ps.password()).To(Equal(`bcdea`))
			ps.undo(`rotate left 1 step`)
			Expect(ps.password()).To(Equal(`abcde`))

			ps.undo(`reverse positions 0 through 4`)
			Expect(ps.password()).To(Equal(`edcba`))

			ps.undo(`swap letter d with letter b`)
			Expect(ps.password()).To(Equal(`ebcda`))

			ps.undo(`swap position 4 with position 0`)
			Expect(ps.password()).To(Equal(`abcde`))
		})

		It("reverses properly", func() {
			ps := NewPasswordScrambler(`gdhcbaef`)
			ps.do(`reverse positions 3 through 6`)
			Expect(ps.password()).To(Equal(`gdheabcf`))
		})

		It("unreverses properly", func() {
			ps := NewPasswordScrambler(`gdheabcf`)
			ps.undo(`reverse positions 3 through 6`)
			Expect(ps.password()).To(Equal(`gdhcbaef`))
		})

		// It("undoes rotation based on position", func() {
		// 	ps := NewPasswordScrambler(`abcdef`)
		// 	ps.do(`rotate based on position of letter a`)
		// 	ps.undo(`rotate based on position of letter a`)
		// 	Expect(ps.password()).To(Equal(`abcdef`))

		// 	ps.do(`rotate based on position of letter b`)
		// 	ps.undo(`rotate based on position of letter b`)
		// 	Expect(ps.password()).To(Equal(`abcdef`))

		// 	ps.do(`rotate based on position of letter c`)
		// 	ps.undo(`rotate based on position of letter c`)
		// 	Expect(ps.password()).To(Equal(`abcdef`))

		// 	ps.do(`rotate based on position of letter d`)
		// 	ps.undo(`rotate based on position of letter d`)
		// 	Expect(ps.password()).To(Equal(`abcdef`))

		// 	ps.do(`rotate based on position of letter e`)
		// 	ps.undo(`rotate based on position of letter e`)
		// 	Expect(ps.password()).To(Equal(`abcdef`))

		// 	ps.do(`rotate based on position of letter f`)
		// 	ps.undo(`rotate based on position of letter f`)
		// 	Expect(ps.password()).To(Equal(`abcdef`))
		// })
	})

	Describe("the puzzle", func() {
		data, _ := ioutil.ReadFile("day21.txt")
		commands := strings.Split(string(data), "\n")

		It("star 1", func() {
			ps := NewPasswordScrambler("abcdefgh")
			for _, command := range commands {
				if blankStringRe.MatchString(command) {
					continue
				}
				ps.do(command)
			}
			fmt.Println("day 21 star 1: password is", ps.password())
		})

		It("star 2", func() {
			ps := NewPasswordScrambler("fbgdceah")
			for j := len(commands)-1; j >= 0; j-- {
				if blankStringRe.MatchString(commands[j]) {
					continue
				}
				ps.undo(commands[j])
			}
			fmt.Println("day 21 star 2: unscrambled password is", ps.password())
		})
	})
})
