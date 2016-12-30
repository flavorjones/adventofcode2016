package adventofcode2016_test

import (
	"fmt"
	"github.com/MakeNowJust/heredoc"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"regexp"
	"strconv"
	"strings"
)

type AssembunnyProcessor struct {
	a, b, c, d   int // registers
	ip           int // instruction pointer
	instructions []string
}

func NewAssembunnyProcessor() *AssembunnyProcessor {
	return &AssembunnyProcessor{0, 0, 0, 0, 0, nil}
}

// instruction set
var apInteger = `-?[0-9]+`
var apIntegerRe = regexp.MustCompile(apInteger)
var apRegister = `[abcd]`
var apRegisterRe = regexp.MustCompile(apRegister)
var apCpyRe = regexp.MustCompile(fmt.Sprintf("cpy (%s|%s) (%s)", apInteger, apRegister, apRegister))
var apIncRe = regexp.MustCompile(fmt.Sprintf("inc (%s)", apRegister))
var apDecRe = regexp.MustCompile(fmt.Sprintf("dec (%s)", apRegister))
var apJnzRe = regexp.MustCompile(fmt.Sprintf("jnz (%s|%s) (%s)", apInteger, apRegister, apInteger))

func (ap *AssembunnyProcessor) Register(registerName string) *int {
	switch registerName {
	case "a":
		return &(ap.a)
	case "b":
		return &(ap.b)
	case "c":
		return &(ap.c)
	case "d":
		return &(ap.d)
	default:
		panic(fmt.Sprintf("unknown register: `%s`", registerName))
	}
}

func (ap *AssembunnyProcessor) Next() {
	ap.ip++
}

func (ap *AssembunnyProcessor) Jump(offset int) {
	ap.ip += offset
}

func (ap *AssembunnyProcessor) Run(program string) {
	ap.ip = 0
	ap.instructions = strings.Split(program, "\n")

	for ap.ip < len(ap.instructions) {
		instruction := ap.instructions[ap.ip]

		if blankStringRe.MatchString(instruction) {
			ap.Next()
			continue
		}

		switch {
		case apCpyRe.MatchString(instruction):
			matches := apCpyRe.FindStringSubmatch(instruction)
			src, dst := matches[1], matches[2]

			switch {
			case apIntegerRe.MatchString(src):
				value, _ := strconv.Atoi(src)
				*(ap.Register(dst)) = value
			case apRegisterRe.MatchString(src):
				*(ap.Register(dst)) = *ap.Register(src)
			default:
				panic(fmt.Sprintf("unknown src arg for cpy: `%s`", src))
			}
			ap.Next()

		case apIncRe.MatchString(instruction):
			matches := apIncRe.FindStringSubmatch(instruction)
			register := matches[1]

			*(ap.Register(register))++
			ap.Next()

		case apDecRe.MatchString(instruction):
			matches := apDecRe.FindStringSubmatch(instruction)
			register := matches[1]

			*(ap.Register(register))--
			ap.Next()

		case apJnzRe.MatchString(instruction):
			matches := apJnzRe.FindStringSubmatch(instruction)
			subject := matches[1]
			offset, _ := strconv.Atoi(matches[2])

			var subjectValue int
			switch {
			case apRegisterRe.MatchString(subject):
				subjectValue = *(ap.Register(subject))
			case apIntegerRe.MatchString(subject):
				subjectValue, _ = strconv.Atoi(subject)
			default:
				panic(fmt.Sprintf("unknown subject arg for jnz: `%s`", subject))
			}

			if subjectValue == 0 {
				ap.Next()
				continue
			}
			ap.Jump(offset)

		default:
			panic(fmt.Sprintf("unknown instruction: `%s`", instruction))
		}
	}
}

// ----------------------------------------
// tests

var _ = Describe("Day12", func() {
	Describe("AssembunnyProcessor", func() {
		var ap *AssembunnyProcessor

		BeforeEach(func() {
			ap = NewAssembunnyProcessor()
		})

		Describe("`cpy`", func() {
			It("copies an integer to a register", func() {
				ap.Run("cpy 3 a")
				Expect(ap.a).To(Equal(3))

				ap.Run("cpy -3 b")
				Expect(ap.b).To(Equal(-3))

				ap.Run("cpy 99 c")
				Expect(ap.c).To(Equal(99))

				ap.Run("cpy 999 d")
				Expect(ap.d).To(Equal(999))
			})

			It("copies an register to a register", func() {
				ap.Run("cpy 3 a")
				ap.Run("cpy a b")
				Expect(ap.b).To(Equal(3))
			})
		})

		Describe("`inc`", func() {
			It("increments a register", func() {
				ap.Run("cpy 10 a")
				ap.Run("inc a")
				Expect(ap.a).To(Equal(11))
				ap.Run("inc a")
				Expect(ap.a).To(Equal(12))
			})
		})

		Describe("`dec`", func() {
			It("decrements a register", func() {
				ap.Run("cpy 10 a")
				ap.Run("dec a")
				Expect(ap.a).To(Equal(9))
				ap.Run("dec a")
				Expect(ap.a).To(Equal(8))
			})
		})

		Describe("instructions", func() {
			It("runs until the program reaches its end", func() {
				instructions := heredoc.Doc(`
					cpy 10 a
					dec a
					inc b
				`)
				ap.Run(instructions)
				Expect(ap.a).To(Equal(9))
				Expect(ap.b).To(Equal(1))
			})
		})

		Describe("`jnz`", func() {
			Context("subject is a register", func() {
				It("continues if value is zero", func() {
					instructions := heredoc.Doc(`
            jnz a 2
					  inc b
					  inc c
  				`)
					ap.Run(instructions)
					Expect(ap.b).To(Equal(1)) // not skipped
					Expect(ap.c).To(Equal(1))
				})

				It("jumps ahead if value is positive", func() {
					instructions := heredoc.Doc(`
					  inc a
            jnz a 2
					  inc b
					  inc c
  				`)
					ap.Run(instructions)
					Expect(ap.b).To(Equal(0)) // skipped that line
					Expect(ap.c).To(Equal(1))
				})

				It("jumps back if value is negative", func() {
					instructions := heredoc.Doc(`
						cpy -5 a
					  inc a
            jnz a -1
					  inc b
					  inc c
  				`)
					ap.Run(instructions)
					Expect(ap.a).To(Equal(0)) // after being incremented a few times
					Expect(ap.b).To(Equal(1))
					Expect(ap.c).To(Equal(1))
				})
			})

			Context("subject is an integer", func() {
				It("continues if value is zero", func() {
					instructions := heredoc.Doc(`
            jnz 0 2
					  inc b
					  inc c
  				`)
					ap.Run(instructions)
					Expect(ap.b).To(Equal(1)) // not skipped
					Expect(ap.c).To(Equal(1))
				})

				It("jumps ahead if value is positive", func() {
					instructions := heredoc.Doc(`
            jnz 1 2
					  inc b
					  inc c
  				`)
					ap.Run(instructions)
					Expect(ap.b).To(Equal(0)) // skipped that line
					Expect(ap.c).To(Equal(1))
				})

				It("jumps back if value is negative", func() {
					instructions := heredoc.Doc(`
					  dec a
						inc a
						jnz a 3
            jnz 1 -2
					  inc b
					  inc c
  				`)
					ap.Run(instructions)
					Expect(ap.a).To(Equal(1)) // after being incremented twice
					Expect(ap.b).To(Equal(0)) // skipped that line
					Expect(ap.c).To(Equal(1))
				})
			})
		})

		It("is Turing-complete", func() {
			instructions := heredoc.Doc(`
        cpy 41 a
        inc a
        inc a
        dec a
        jnz a 2
        dec a
			`)
			ap := NewAssembunnyProcessor()
			ap.Run(instructions)
			Expect(ap.a).To(Equal(42))
		})
	})

	Describe("the puzzle", func() {
		instructions := heredoc.Doc(`
			cpy 1 a
      cpy 1 b
      cpy 26 d
      jnz c 2
      jnz 1 5
      cpy 7 c
      inc d
      dec c
      jnz c -2
      cpy a c
      inc a
      dec b
      jnz b -2
      cpy c b
      dec d
      jnz d -6
      cpy 19 c
      cpy 11 d
      inc a
      dec d
      jnz d -2
      dec c
      jnz c -5
		`)

		Describe("star 1", func() {
			It("doesn't halt", func() {
				ap := NewAssembunnyProcessor()
				ap.Run(instructions)
				fmt.Println("star 1: register 'a' is", ap.a)
			})
		})

		Describe("star 2", func() {
			It("doesn't halt", func() {
				ap := NewAssembunnyProcessor()
				ap.Run("cpy 1 c")
				ap.Run(instructions)
				fmt.Println("star 2: register 'a' is", ap.a)
			})
		})
	})
})
